# Plano de Implementação: Minimapa (2c)

> **Para workers agênticos:** SUB-SKILL OBRIGATÓRIA: use superpowers:subagent-driven-development (recomendado) ou superpowers:executing-plans para implementar este plano task por task. Os passos usam a sintaxe de checkbox (`- [ ]`) para acompanhamento.

**Objetivo:** Mostrar um overlay simples no canto da tela com as salas já
visitadas da construção atual, destacando qual é a sala onde o jogador
está agora — a sala não-visitada simplesmente não aparece (sem o estado
intermediário "visível mas não visitada" do Isaac — corte de escopo
deliberado, ver Arquitetura).

**Arquitetura:** `building.Building` ganha rastreio de quais salas já
foram visitadas (hoje só sabe qual é a sala atual). A posição/cor de
cada célula do minimapa é calculada por uma função pura e testável
(`minimap.LayoutCells`), separada do desenho de verdade — o mesmo padrão
já usado no resto do projeto pra manter lógica testável sem precisar de
tela real (`Player.Move` testável sem `InputHandler`, `Game.CheckRoomTransition`
testável sem rodar `Update()`). O desenho em si usa
`github.com/hajimehoshi/ebiten/v2/vector` (subpacote do ebiten já
presente no `go.mod`, nenhuma dependência nova) pra pintar retângulos
sólidos — não existe nenhum asset de minimapa, só cor sólida.

**Simplificação deliberada**: o spec menciona "minimapa mostrando salas
visitadas/não visitadas", no estilo Isaac, que também mostra salas
*conhecidas-mas-não-visitadas* (adjacentes a uma visitada) com um
contorno apagado. Esse plano só tem 2 estados — visitada (aparece) e
não-visitada (não aparece) — porque é o suficiente pra confirmar que o
sistema de salas existe e orientar o jogador; o terceiro estado fica pra
quando (se) isso incomodar na prática.

**Stack técnica:** Go, `github.com/hajimehoshi/ebiten/v2` (incluindo
`.../v2/vector`), `github.com/quasilyte/ebitengine-input`,
`github.com/stretchr/testify`.

**Spec:** [docs/superpowers/specs/2026-10-17-dengue-boy-isaac-redesign-design.md](../specs/2026-10-17-dengue-boy-isaac-redesign-design.md)

## Restrições globais

- Stack inalterada. `.../v2/vector` é parte do módulo
  `github.com/hajimehoshi/ebiten/v2` já presente no `go.mod` — não conta
  como dependência nova.
- `go build ./...` e `go test ./...` precisam passar no final do plano.
- Este plano modifica `building.Building` (campo novo + um construtor
  novo) e o corpo de `building.TryMove` e `generator.Generate` dos
  planos 2a/2b — sem mudar a assinatura pública de nenhum dos dois.

## Foco de revisão

- A sala de spawn precisa nascer marcada como visitada — se não nascer,
  o minimapa abre vazio até o jogador se mexer uma vez, o que parece bug
  — responsabilidade da Task 1.
- Uma sala não-visitada nunca pode aparecer no minimapa, mesmo que seja
  vizinha direta da sala atual — responsabilidade da Task 2
  (`LayoutCells`).
- A posição de cada célula no minimapa precisa ser relativa à caixa
  delimitadora de **todas** as salas da construção (visitadas ou não),
  não só das visitadas — senão o layout do minimapa "pula"/reposiciona
  toda vez que uma sala nova é visitada, o que é desorientador —
  responsabilidade da Task 2.
- Exatamente uma célula (a sala atual) precisa estar marcada como
  `Current: true` no resultado de `LayoutCells`, nunca zero nem mais de
  uma — responsabilidade da Task 2.
- Uma construção de uma sala só (sem porta, caso `roomCount == 1` do
  plano 2b) não pode fazer `LayoutCells` falhar nem dividir por zero no
  cálculo da caixa delimitadora (min/max com um único ponto) —
  responsabilidade da Task 2; como `Draw` (Task 3) só itera o que
  `LayoutCells` devolve, cobrir isso na Task 2 já cobre `Draw` também.

---

### Task 1: Rastreio de salas visitadas em `Building`

**Arquivos:**
- Modificar: `internal/modules/building/building.go`
- Modificar: `internal/modules/building/building_test.go`
- Modificar: `internal/modules/generator/generator.go`

**Interfaces:**
- Produz:
  - `Building.Visited map[Coord]bool` (campo novo).
  - `func building.New(rooms map[Coord]*room.Room, start Coord) *Building`
    — monta `&Building{Rooms: rooms, Current: start, Visited:
    map[Coord]bool{start: true}}`.
  - `TryMove` (já existente) passa a marcar `b.Visited[target] = true`
    no momento em que a transição é bem-sucedida (mesmo ponto onde já
    atualiza `b.Current`).
  - Consumida pela Task 2 (`minimap.LayoutCells`).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
func TestNewMarksStartRoomAsVisited(t *testing.T) {
	spawn := &room.Room{Cleared: true}
	b := building.New(map[building.Coord]*room.Room{{0, 0}: spawn}, building.Coord{0, 0})
	require.True(t, b.Visited[building.Coord{0, 0}])
	require.Len(t, b.Visited, 1)
}

func TestTryMoveMarksDestinationAsVisited(t *testing.T) {
	b := twoConnectedRooms() // já existe em building_test.go
	b.Visited = map[building.Coord]bool{{0, 0}: true}

	b.TryMove(room.East)

	require.True(t, b.Visited[building.Coord{1, 0}])
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/building/... -v`
Esperado: FAIL — `building.New` e `Building.Visited` ainda não existem.

- [ ] **Passo 3: Implementar**

Em `internal/modules/building/building.go`: adicionar o campo
`Visited map[Coord]bool` à struct `Building`; implementar `New` como
descrito; em `TryMove`, logo após `b.Current = target`, adicionar `if
b.Visited == nil { b.Visited = map[Coord]bool{} }; b.Visited[target] =
true` (o `nil`-check cobre `Building`s montados por struct literal direto,
como os testes existentes do plano 2a, que não passam por `New`).

Em `internal/modules/generator/generator.go`: trocar a construção final
do `*building.Building` (struct literal) por `building.New(rooms,
building.Coord{X: 0, Y: 0})`.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/building/... ./internal/modules/generator/... -v`
Esperado: PASS — incluindo os testes já existentes dos planos 2a e 2b
(nenhum deles deve quebrar).

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/building/building.go internal/modules/building/building_test.go internal/modules/generator/generator.go
git commit -m "Track visited rooms in Building"
```

---

### Task 2: Layout puro do minimapa

**Arquivos:**
- Criar: `internal/modules/minimap/minimap.go`
- Criar: `internal/modules/minimap/minimap_test.go`

**Interfaces:**
- Consome: `building.Building`, `building.Coord` (Task 1 + plano 2a).
- Produz:
  - `type minimap.Cell struct { Coord building.Coord; ScreenX, ScreenY float64; Current bool }`
  - `func minimap.LayoutCells(b *building.Building) []Cell` — uma
    célula por sala em `b.Visited` (salas fora de `b.Visited` não geram
    célula). `ScreenX = marginX + float64(coord.X-minX)*(cellSize+cellGap)`,
    `ScreenY = marginY + float64(coord.Y-minY)*(cellSize+cellGap)`, onde
    `minX`/`minY` são o menor `X`/`Y` entre **todas** as chaves de
    `b.Rooms` (visitadas ou não). `Current` é `true` só para a célula
    cujo `Coord == b.Current`. Constantes exportadas:
    `minimap.CellSize = 10.0`, `minimap.CellGap = 2.0`,
    `minimap.MarginX = 12.0`, `minimap.MarginY = 12.0`.
  - Consumida pela Task 3 (`minimap.Draw`).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package minimap_test

import (
	"app/internal/modules/building"
	"app/internal/modules/minimap"
	"app/internal/modules/room"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLayoutCellsOnlyIncludesVisitedRooms(t *testing.T) {
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: {}, {1, 0}: {},
		},
		Current: building.Coord{0, 0},
		Visited: map[building.Coord]bool{{0, 0}: true},
	}

	cells := minimap.LayoutCells(b)

	require.Len(t, cells, 1)
	require.Equal(t, building.Coord{0, 0}, cells[0].Coord)
}

func TestLayoutCellsFlagsExactlyTheCurrentRoom(t *testing.T) {
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: {}, {1, 0}: {},
		},
		Current: building.Coord{1, 0},
		Visited: map[building.Coord]bool{{0, 0}: true, {1, 0}: true},
	}

	cells := minimap.LayoutCells(b)

	currentCount := 0
	for _, c := range cells {
		if c.Current {
			currentCount++
			require.Equal(t, building.Coord{1, 0}, c.Coord)
		}
	}
	require.Equal(t, 1, currentCount)
}

func TestLayoutCellsPositionsRelativeToFullBoundingBox(t *testing.T) {
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: {}, {1, 0}: {}, {2, 0}: {},
		},
		Current: building.Coord{2, 0},
		// {1,0} exists in Rooms but is NOT visited — the bounding box
		// must still account for it.
		Visited: map[building.Coord]bool{{0, 0}: true, {2, 0}: true},
	}

	cells := minimap.LayoutCells(b)
	require.Len(t, cells, 2)

	byCoord := map[building.Coord]minimap.Cell{}
	for _, c := range cells {
		byCoord[c.Coord] = c
	}

	require.Equal(t, minimap.MarginX, byCoord[building.Coord{0, 0}].ScreenX)
	require.Equal(t, minimap.MarginY, byCoord[building.Coord{0, 0}].ScreenY)
	require.Equal(t, minimap.MarginX+2*(minimap.CellSize+minimap.CellGap), byCoord[building.Coord{2, 0}].ScreenX)
	require.Equal(t, minimap.MarginY, byCoord[building.Coord{2, 0}].ScreenY)
}

func TestLayoutCellsWithSingleRoomBuilding(t *testing.T) {
	b := &building.Building{
		Rooms:   map[building.Coord]*room.Room{{0, 0}: {}},
		Current: building.Coord{0, 0},
		Visited: map[building.Coord]bool{{0, 0}: true},
	}

	cells := minimap.LayoutCells(b)

	require.Len(t, cells, 1)
	require.Equal(t, minimap.MarginX, cells[0].ScreenX)
	require.Equal(t, minimap.MarginY, cells[0].ScreenY)
	require.True(t, cells[0].Current)
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/minimap/... -v`
Esperado: FAIL — pacote `minimap` ainda não existe.

- [ ] **Passo 3: Implementar `Cell`, as 4 constantes e `LayoutCells` em `internal/modules/minimap/minimap.go`**

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/minimap/... -v`
Esperado: PASS (4 testes)

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/minimap/minimap.go internal/modules/minimap/minimap_test.go
git commit -m "Add pure minimap cell layout"
```

---

### Task 3: Desenhar o minimapa e ligar no `Game`

**Arquivos:**
- Modificar: `internal/modules/minimap/minimap.go`
- Modificar: `internal/modules/game/game.go`

**Interfaces:**
- Consome: `minimap.LayoutCells`, `minimap.Cell` + constantes (Task 2).
- Produz:
  - `func minimap.Draw(screen *ebiten.Image, b *building.Building)` —
    para cada `Cell` de `LayoutCells(b)`, desenha um quadrado
    `CellSize × CellSize` em `(ScreenX, ScreenY)` via
    `vector.DrawFilledRect`: cor `color.RGBA{255, 255, 255, 255}`
    (branco) se `Cell.Current`, senão `color.RGBA{150, 150, 150, 255}`
    (cinza).
  - `Game.Draw` passa a chamar `minimap.Draw(screen, g.Building)` depois
    de `g.Entities.Draw(screen)`.

- [ ] **Passo 1: Implementar `Draw` em `internal/modules/minimap/minimap.go`**

Não há teste automatizado pra esse passo — desenhar numa `*ebiten.Image`
de verdade não é verificável sem renderizar a tela; a lógica que importa
(quais células, onde, qual cor) já está coberta pelos testes de
`LayoutCells` na Task 2. `Draw` é só o laço fino que aplica isso.

- [ ] **Passo 2: Ligar no `Game`**

Em `internal/modules/game/game.go`, em `Draw`, adicionar a chamada a
`minimap.Draw(screen, g.Building)` depois de `g.Entities.Draw(screen)`.

- [ ] **Passo 3: Verificar que o projeto inteiro compila e os testes passam**

Rodar: `go build ./...`
Esperado: sem erros.

Rodar: `go test ./...`
Esperado: todos os pacotes passam.

- [ ] **Passo 4: Verificação manual**

Rodar: `go run cmd/main.go` (ou `make run`). Confirmar: um quadrado
branco aparece no canto superior esquerdo assim que o jogo abre (a sala
de spawn); andar por uma porta pra uma sala nova faz aparecer um
quadrado cinza na posição relativa certa, e o quadrado branco (sala
atual) se move pra lá.

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/minimap/minimap.go internal/modules/game/game.go
git commit -m "Draw the minimap overlay"
```
