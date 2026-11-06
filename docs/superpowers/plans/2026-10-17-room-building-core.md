# Plano de Implementação: Núcleo de Salas & Construção (2a)

> **Para workers agênticos:** SUB-SKILL OBRIGATÓRIA: use superpowers:subagent-driven-development (recomendado) ou superpowers:executing-plans para implementar este plano task por task. Os passos usam a sintaxe de checkbox (`- [ ]`) para acompanhamento.

**Objetivo:** Substituir o mapa único e fixo de hoje por um sistema de
salas conectadas por portas: o jogador anda por uma sala, encosta numa
porta, e troca de sala — com a colisão e a renderização trocando junto.
Prova o fluxo ponta a ponta com 2-3 salas desenhadas à mão; geração
procedural (2b) e minimapa (2c) constroem em cima disso em planos
separados.

**Arquitetura:** Três pacotes novos/alterados. `maptile` ganha um tile de
porta (`Door`), que já fica andável de graça (a colisão hoje só bloqueia
`Wall`). Um pacote novo `room` define `Room` (um `*maptile.Map` + quais
lados têm porta + se a sala está "limpa") e as funções puras de geometria
de porta (posição fixa por lado, já que toda sala usa a grade 25×20 de
hoje). Um pacote novo `building` define `Building` (grade de salas +
sala atual) e a navegação entre elas (`TryMove`), incluindo a trava por
`Cleared` — que fica com a lógica real mas sempre `true` nas salas deste
plano, já que ainda não existe nada (foco de dengue, inimigo) pra setar
`false`; isso é o gancho que os planos de Limpeza e Inimigos vão usar
depois. Câmera não é uma tarefa separada: como cada sala continua do
tamanho exato da tela (25×20 tiles, decisão já tomada na Fundação),
"câmera" é só desenhar a sala atual na origem — nenhum código novo é
necessário além de trocar `Game.Map` por `Game.Building.CurrentRoom().Map`.

**Stack técnica:** Go, `github.com/hajimehoshi/ebiten/v2`,
`github.com/quasilyte/ebitengine-input`, `github.com/stretchr/testify`.

**Spec:** [docs/superpowers/specs/2026-10-17-dengue-boy-isaac-redesign-design.md](../specs/2026-10-17-dengue-boy-isaac-redesign-design.md)

## Restrições globais

- Stack inalterada: Go + `github.com/hajimehoshi/ebiten/v2` +
  `github.com/quasilyte/ebitengine-input`. Nenhuma dependência nova.
- Toda sala usa a grade 25×20 tiles de 32px (`maptile.ImageSize`) — igual
  ao mapa único de hoje. Tamanho de sala variável fica fora de escopo
  deste plano.
- Porta fica sempre na mesma posição fixa por lado (não configurável por
  sala) — simplificação deste plano; a geração procedural (2b) também
  vai seguir essa mesma convenção.
- Sem câmera com scroll — cada sala renderiza na origem (0,0), decisão já
  tomada no plano da Fundação.
- `go build ./...` precisa compilar no final do plano. Entre tasks que
  trocam `Player.CollisionSystem`'s conteúdo e `Game.Map`→`Game.Building`,
  o build pode quebrar temporariamente — mesmo padrão já usado no plano
  da Fundação; cada task escopa sua verificação ao próprio pacote até a
  task final restaurar o build completo.

## Foco de revisão

- Uma sala declara porta num lado mas a sala vizinha não declara a porta
  correspondente do lado oposto (grafo assimétrico) — `TryMove` precisa
  falhar no lado que não reciproca, nunca teleportar o jogador pra uma
  sala sem volta — responsabilidade da Task 2 (`building.TryMove`).
- Uma sala sem porta num lado (`HasDoor` falso) nunca pode disparar
  transição ali, mesmo que o jogador esteja exatamente na posição fixa
  onde a porta ficaria se existisse — a checagem de zona de porta tem
  que ser condicionada a `HasDoor`, não só à coordenada — responsabilidade
  da Task 4 (wiring no `Game`).
- Uma sala com `Cleared = false` precisa bloquear todas as portas dela,
  não só "mostrar trancada" visualmente — responsabilidade da Task 2.
- Trocar de sala precisa reconstruir o `CollisionSystem` com as paredes
  da sala NOVA, não manter os objetos da sala antiga — senão o jogador
  atravessa parede (ou trava em parede fantasma) logo depois da
  transição — responsabilidade da Task 4.
- Andar na diagonal em direção a uma porta não pode deixar o jogador
  "cortar caminho" através do canto da parede vizinha à porta (que
  continua sendo `Wall` de verdade) — o deslize por eixo já implementado
  na Fundação deveria cobrir isso, mas merece um teste de regressão
  específico com uma porta real no meio — responsabilidade da Task 4.

---

### Task 1: Tile de porta + tipo `Room`

**Arquivos:**
- Modificar: `internal/modules/maptile/map_tile.go` (novo tile `Door`)
- Modificar: `internal/modules/maptile/renderer.go` (imagem do tile `Door`)
- Criar: `internal/modules/room/room.go`
- Criar: `internal/modules/room/room_test.go`

**Interfaces:**
- Produz:
  - `maptile.Door MapTile = 5` (novo valor no enum existente).
  - `type room.Side int` com `room.North`, `room.South`, `room.East`,
    `room.West` (nessa ordem, `iota`).
  - `func room.Opposite(side Side) Side` — `North↔South`, `East↔West`.
  - `type room.Room struct { Map *maptile.Map; Doors map[Side]bool; Cleared bool }`
  - `func (r *Room) HasDoor(side Side) bool`
  - Constantes exportadas de posição fixa de porta (coluna/linha na
    grade 25×20, `0`-indexado): `room.DoorColumnNS = 12` (coluna da
    porta norte/sul), `room.DoorRowNorth = 0`, `room.DoorRowSouth = 19`,
    `room.DoorRowEW = 9` (linha da porta leste/oeste),
    `room.DoorColumnWest = 0`, `room.DoorColumnEast = 24`.
  - `func room.NewBlankRoomTiles(doors map[Side]bool) [][]int` — gera uma
    grade 25×20 com borda de `Wall` em todo o perímetro, exceto nas
    posições fixas de porta acima, para os lados presentes em `doors`
    (com valor `true`), que viram `Door`. Função pura, sem
    `Images`/`ebiten`.
  - Consumida pela Task 2 (`building.Building`), pela Task 3 (`DoorZone`/
    `EntryPosition` reusam estas mesmas constantes) e pela Task 5 (salas
    de demonstração).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package room_test

import (
	"app/internal/modules/maptile"
	"app/internal/modules/room"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOppositeSide(t *testing.T) {
	require.Equal(t, room.South, room.Opposite(room.North))
	require.Equal(t, room.North, room.Opposite(room.South))
	require.Equal(t, room.West, room.Opposite(room.East))
	require.Equal(t, room.East, room.Opposite(room.West))
}

func TestRoomHasDoor(t *testing.T) {
	r := &room.Room{Doors: map[room.Side]bool{room.North: true}}
	require.True(t, r.HasDoor(room.North))
	require.False(t, r.HasDoor(room.South))
}

func TestNewBlankRoomTilesHasBorderWalls(t *testing.T) {
	tiles := room.NewBlankRoomTiles(map[room.Side]bool{})
	require.Len(t, tiles, 20)
	require.Len(t, tiles[0], 25)
	require.Equal(t, int(maptile.Wall), tiles[0][0], "top-left corner must stay a wall")
	require.Equal(t, int(maptile.Empty), tiles[10][12], "room interior must be empty")
}

func TestNewBlankRoomTilesOpensDoorOnRequestedSide(t *testing.T) {
	tiles := room.NewBlankRoomTiles(map[room.Side]bool{room.North: true})
	require.Equal(t, int(maptile.Door), tiles[0][12], "north door must open at the fixed column")
	require.Equal(t, int(maptile.Wall), tiles[19][12], "south side must stay closed when not requested")
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/room/... -v`
Esperado: FAIL — pacote `room` ainda não existe.

- [ ] **Passo 3: Implementar**

Em `internal/modules/maptile/map_tile.go`: adicionar `Door MapTile = 5`
ao bloco de constantes existente.

Em `internal/modules/maptile/renderer.go`: em `loadImageTiles`, adicionar
`Door: emptyImage` ao mapa retornado — reaproveita a mesma imagem do tile
`Empty` como placeholder visual (sem asset novo; o plano de Sprites, item
7 do roadmap, cuida da arte de verdade depois).

Em `internal/modules/room/room.go`: implementar `Side`, as 4 constantes
de lado, `Opposite`, `Room`, `HasDoor`, as 6 constantes de posição de
porta listadas em Interfaces, e `NewBlankRoomTiles` usando essas
constantes pra decidir em qual coluna/linha trocar `Wall` por `Door`.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/room/... -v`
Esperado: PASS (4 testes)

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/maptile/map_tile.go internal/modules/maptile/renderer.go internal/modules/room/room.go internal/modules/room/room_test.go
git commit -m "Add door tile and Room type"
```

---

### Task 2: `building.Building` — grade de salas e navegação

**Arquivos:**
- Criar: `internal/modules/building/building.go`
- Criar: `internal/modules/building/building_test.go`

**Interfaces:**
- Consome: `room.Room`, `room.Side`, `room.Opposite` (Task 1).
- Produz:
  - `type building.Coord struct { X, Y int }`
  - `func building.Neighbor(c Coord, side room.Side) Coord` — `North:
    {X,Y-1}`, `South: {X,Y+1}`, `East: {X+1,Y}`, `West: {X-1,Y}`.
  - `type building.Building struct { Rooms map[Coord]*room.Room; Current Coord }`
  - `func (b *Building) CurrentRoom() *room.Room`
  - `func (b *Building) TryMove(side room.Side) bool` — retorna `false`
    sem mudar nada se: a sala atual não tem porta nesse lado, OU a sala
    atual não está `Cleared`, OU não existe sala na coordenada vizinha,
    OU a sala vizinha não tem porta no lado oposto. Caso contrário,
    atualiza `Current` para a vizinha e retorna `true`.
  - Consumida pela Task 4 (`Game`) e pelo plano 2b (geração procedural).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package building_test

import (
	"app/internal/modules/building"
	"app/internal/modules/room"
	"testing"

	"github.com/stretchr/testify/require"
)

func twoConnectedRooms() *building.Building {
	spawn := &room.Room{Doors: map[room.Side]bool{room.East: true}, Cleared: true}
	east := &room.Room{Doors: map[room.Side]bool{room.West: true}, Cleared: true}
	return &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: spawn,
			{1, 0}: east,
		},
		Current: building.Coord{0, 0},
	}
}

func TestTryMoveSucceedsWhenBothSidesHaveMatchingDoors(t *testing.T) {
	b := twoConnectedRooms()
	ok := b.TryMove(room.East)
	require.True(t, ok)
	require.Equal(t, building.Coord{1, 0}, b.Current)
}

func TestTryMoveFailsWhenCurrentRoomHasNoDoorThatSide(t *testing.T) {
	b := twoConnectedRooms()
	ok := b.TryMove(room.North)
	require.False(t, ok)
	require.Equal(t, building.Coord{0, 0}, b.Current)
}

func TestTryMoveFailsWhenNeighborRoomDoesNotExist(t *testing.T) {
	b := twoConnectedRooms()
	ok := b.TryMove(room.South)
	require.False(t, ok)
}

func TestTryMoveFailsWhenNeighborHasNoMatchingDoor(t *testing.T) {
	spawn := &room.Room{Doors: map[room.Side]bool{room.East: true}, Cleared: true}
	eastWithoutWestDoor := &room.Room{Doors: map[room.Side]bool{}, Cleared: true}
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: spawn,
			{1, 0}: eastWithoutWestDoor,
		},
		Current: building.Coord{0, 0},
	}
	ok := b.TryMove(room.East)
	require.False(t, ok, "neighbor must reciprocate the door on the opposite side")
}

func TestTryMoveFailsWhenCurrentRoomIsNotCleared(t *testing.T) {
	spawn := &room.Room{Doors: map[room.Side]bool{room.East: true}, Cleared: false}
	east := &room.Room{Doors: map[room.Side]bool{room.West: true}, Cleared: true}
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: spawn,
			{1, 0}: east,
		},
		Current: building.Coord{0, 0},
	}
	ok := b.TryMove(room.East)
	require.False(t, ok)
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/building/... -v`
Esperado: FAIL — pacote `building` ainda não existe.

- [ ] **Passo 3: Implementar `Coord`, `Neighbor`, `Building`, `CurrentRoom`, `TryMove` em `internal/modules/building/building.go`**

`TryMove` é a única função com lógica real; as outras são acesso direto a
campos/mapas.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/building/... -v`
Esperado: PASS (5 testes)

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/building/building.go internal/modules/building/building_test.go
git commit -m "Add Building grid with door-gated room navigation"
```

---

### Task 3: Geometria de porta (posição fixa, zona de colisão, ponto de entrada)

**Arquivos:**
- Modificar: `internal/modules/room/room.go`
- Modificar: `internal/modules/room/room_test.go`

**Interfaces:**
- Consome: `room.Side` e as 6 constantes de posição de porta (Task 1),
  `collision.Object`, `maptile.ImageSize` (já existentes).
- Produz:
  - `func room.DoorZone(side Side) collision.Object` — retorna o
    retângulo 32×32 (em pixels, `coluna*32, linha*32`) da porta daquele
    lado, usando as constantes acima.
  - `func room.EntryPosition(side Side) geom.Position` — a posição em
    pixels, um tile pra dentro da borda, pro lado `side` (usada quando o
    jogador entra numa sala vindo daquele lado).
  - Consumida pela Task 4.

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
func TestDoorZonePositions(t *testing.T) {
	require.Equal(t, collision.Object{Position: geom.Position{X: 384, Y: 0}, Width: 32, Height: 32}, room.DoorZone(room.North))
	require.Equal(t, collision.Object{Position: geom.Position{X: 384, Y: 608}, Width: 32, Height: 32}, room.DoorZone(room.South))
	require.Equal(t, collision.Object{Position: geom.Position{X: 768, Y: 288}, Width: 32, Height: 32}, room.DoorZone(room.East))
	require.Equal(t, collision.Object{Position: geom.Position{X: 0, Y: 288}, Width: 32, Height: 32}, room.DoorZone(room.West))
}

func TestEntryPositions(t *testing.T) {
	require.Equal(t, geom.Position{X: 384, Y: 32}, room.EntryPosition(room.North))
	require.Equal(t, geom.Position{X: 384, Y: 576}, room.EntryPosition(room.South))
	require.Equal(t, geom.Position{X: 736, Y: 288}, room.EntryPosition(room.East))
	require.Equal(t, geom.Position{X: 32, Y: 288}, room.EntryPosition(room.West))
}
```

(Adicionar os imports `app/internal/modules/collision` e
`app/internal/modules/geom` ao `room_test.go` existente.)

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/room/... -run "TestDoorZone|TestEntryPositions" -v`
Esperado: FAIL — `DoorZone`/`EntryPosition` ainda não existem.

- [ ] **Passo 3: Implementar `DoorZone` e `EntryPosition` em `internal/modules/room/room.go`**

Reusar as 6 constantes de posição já definidas na Task 1 — não
redeclarar valores novos, pra garantir que a porta visual/tile e a zona
de colisão da porta concordem no mesmo lugar.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/room/... -v`
Esperado: PASS (6 testes no total do pacote)

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/room/room.go internal/modules/room/room_test.go
git commit -m "Add fixed door geometry: zone and entry position per side"
```

---

### Task 4: `CollisionSystem.Rebuild` + transição de sala no `Game`

**Arquivos:**
- Modificar: `internal/modules/collision/collision.go`
- Modificar: `internal/modules/game/game.go`
- Criar: `internal/modules/game/game_test.go` (arquivo já existe da
  Fundação — adicionar os novos testes a ele)

**Interfaces:**
- Consome: `room.DoorZone`, `room.EntryPosition`, `room.Opposite`,
  `room.Room.HasDoor` (Task 1/3); `building.Building.TryMove`,
  `building.Building.CurrentRoom` (Task 2).
- Produz:
  - `func (c *CollisionSystem) Rebuild(mapTile *maptile.Map)` —
    recalcula `c.Objects` a partir do mapa dado (mesma lógica que hoje
    só roda dentro de `collision.New`; `New` passa a chamar `Rebuild`
    internamente, sem mudar sua assinatura).
  - `Game.Building *building.Building` substitui `Game.Map
    *maptile.Map`.
  - `func New(building *building.Building, player *player.Player, inputSystem input.System) *Game`
    — mesma assinatura de hoje, só com `building *building.Building` no
    lugar do antigo `maptile *maptile.Map`.
  - `func (g *Game) CheckRoomTransition()` — **exportado** (o teste deste
    pacote chama direto, sem passar por `Update()`, pra não precisar
    simular `InputSystem`). Chamado a partir de `Update()`. Para cada
    lado em `[North, South, East, West]`: se a sala atual tem porta
    nesse lado (`HasDoor`) E a posição do jogador colide com
    `room.DoorZone(side)` E `g.Building.TryMove(side)` retorna `true`,
    então: `g.Player.Position = room.EntryPosition(room.Opposite(side))`
    e `g.Player.CollisionSystem.Rebuild(g.Building.CurrentRoom().Map)`.
  - Consumida pela Task 5 (`cmd/main.go`).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package collision_test
// (ou adicionar ao arquivo collision_test.go existente, se já houver — caso
// contrário criar internal/modules/collision/collision_test.go)

func TestRebuildReplacesObjectsFromNewMap(t *testing.T) {
	c := &collision.CollisionSystem{
		Objects: []collision.Object{{Position: geom.Position{X: 999, Y: 999}, Width: 32, Height: 32}},
	}
	newMap := &maptile.Map{Tiles: [][]int{{1, 0}, {0, 1}}}

	c.Rebuild(newMap)

	require.Len(t, c.Objects, 2, "only the two Wall tiles in newMap become objects")
	require.NotContains(t, c.Objects, collision.Object{Position: geom.Position{X: 999, Y: 999}, Width: 32, Height: 32})
}
```

```go
// internal/modules/game/game_test.go — adicionar:

func TestCheckRoomTransitionMovesPlayerAndRebuildsCollision(t *testing.T) {
	spawnTiles := room.NewBlankRoomTiles(map[room.Side]bool{room.East: true})
	eastTiles := room.NewBlankRoomTiles(map[room.Side]bool{room.West: true})
	spawn := &room.Room{Map: &maptile.Map{Tiles: spawnTiles}, Doors: map[room.Side]bool{room.East: true}, Cleared: true}
	east := &room.Room{Map: &maptile.Map{Tiles: eastTiles}, Doors: map[room.Side]bool{room.West: true}, Cleared: true}
	b := &building.Building{
		Rooms:   map[building.Coord]*room.Room{{0, 0}: spawn, {1, 0}: east},
		Current: building.Coord{0, 0},
	}

	cs := &collision.CollisionSystem{}
	cs.Rebuild(spawn.Map)
	p := &player.Player{Position: room.DoorZone(room.East).Position, CollisionSystem: cs}

	g := &game.Game{Building: b, Player: p, Entities: entity.NewManager()}
	g.Entities.Add(p)

	g.CheckRoomTransition()

	require.Equal(t, building.Coord{1, 0}, b.Current)
	require.Equal(t, room.EntryPosition(room.West), p.Position)
}

func TestCheckRoomTransitionDoesNothingWithoutADoor(t *testing.T) {
	windowless := &room.Room{Map: &maptile.Map{Tiles: room.NewBlankRoomTiles(nil)}, Doors: map[room.Side]bool{}, Cleared: true}
	b := &building.Building{Rooms: map[building.Coord]*room.Room{{0, 0}: windowless}, Current: building.Coord{0, 0}}
	cs := &collision.CollisionSystem{}
	cs.Rebuild(windowless.Map)
	// Player standing exactly where a door zone would be if one existed:
	p := &player.Player{Position: room.DoorZone(room.North).Position, CollisionSystem: cs}

	g := &game.Game{Building: b, Player: p, Entities: entity.NewManager()}
	g.Entities.Add(p)

	g.CheckRoomTransition()

	require.Equal(t, building.Coord{0, 0}, b.Current)
	require.Equal(t, room.DoorZone(room.North).Position, p.Position)
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/collision/... ./internal/modules/game/... -v`
Esperado: FAIL — `Rebuild`, `Game.Building` e `CheckRoomTransition`
ainda não existem.

- [ ] **Passo 3: Implementar**

Em `internal/modules/collision/collision.go`: extrair o corpo de
`createCollisionObjects` pra dentro de um novo método
`(c *CollisionSystem) Rebuild(mapTile *maptile.Map)` que faz `c.Objects =
createCollisionObjects(mapTile)`; `New` passa a chamar `Rebuild` no
`CollisionSystem` recém-criado em vez de montar `Objects` direto.

Em `internal/modules/game/game.go`: trocar o campo `Map *maptile.Map`
por `Building *building.Building`; em `Draw`, trocar
`g.Map.Draw(screen)` por `g.Building.CurrentRoom().Map.Draw(screen)`; em
`Update`, chamar `g.CheckRoomTransition()` depois de
`g.Entities.Update()`; implementar `CheckRoomTransition` como descrito
em Interfaces, usando `collision.CollisionSystem{Objects:
[]collision.Object{room.DoorZone(side)}}.HasCollision(...)` pra testar
sobreposição do jogador com a zona da porta (reaproveita o
`HasCollision` já existente, sem nova lógica de AABB). `New(...)`
precisa ser ajustado pra receber `*building.Building` no lugar de
`*maptile.Map` — ajustar a assinatura e quem chama (Task 5 cuida de
`cmd/main.go`).

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/collision/... ./internal/modules/game/... -v`
Esperado: PASS. Nota: `go build ./...` ainda pode falhar aqui —
`cmd/main.go` só é ajustado na Task 5.

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/collision/collision.go internal/modules/game/game.go internal/modules/game/game_test.go
git commit -m "Add collision rebuild and room-transition handling in Game"
```

---

### Task 5: Três salas conectadas + verificação manual

**Arquivos:**
- Modificar: `cmd/main.go`

**Interfaces:**
- Consome: `room.Room`, `room.NewBlankRoomTiles` (Task 1);
  `building.Building`, `building.Coord` (Task 2); `game.New` com a nova
  assinatura (Task 4).

- [ ] **Passo 1: Implementar o layout de 3 salas em `cmd/main.go`**

Substituir a chamada única `maptile.New()` por três salas conectadas,
usando `room.NewBlankRoomTiles` + as imagens já carregadas por
`maptile.loadImageTiles` (reaproveitar o carregamento de imagem
existente, só trocando `Tiles`):

- Sala de spawn em `building.Coord{0, 0}`: porta só a Leste.
- Sala em `building.Coord{1, 0}`: porta a Oeste (volta pro spawn) e ao
  Norte.
- Sala em `building.Coord{1, -1}`: porta só ao Sul (volta pra sala do
  meio).

Todas as 3 salas com `Cleared: true`. Montar o `building.Building` com
essas 3 salas e `Current: building.Coord{0, 0}`. Passar o `Building` pro
`collision.New` (usando a sala atual: `collision.New(initialBuilding.CurrentRoom().Map)`),
pro `player.New` (collision system resultante), e pro `game.New`
(building, player, input system).

- [ ] **Passo 2: Verificar que o projeto inteiro compila e os testes passam**

Rodar: `go build ./...`
Esperado: sem erros.

Rodar: `go test ./...`
Esperado: todos os pacotes passam.

- [ ] **Passo 3: Verificação manual**

Rodar: `go run cmd/main.go` (ou `make run`). Confirmar: anda até a porta
Leste da sala de spawn, atravessa pra sala do meio; anda até a porta
Norte da sala do meio, atravessa pra terceira sala; volta pelas mesmas
portas. Confirmar que não dá pra atravessar onde não tem porta (paredes
continuam bloqueando nos outros 3 lados de cada sala).

- [ ] **Passo 4: Commit**

```bash
git add cmd/main.go
git commit -m "Wire a 3-room connected building into the game"
```
