# Plano de Implementação: Geração Procedural de Salas (2b)

> **Para workers agênticos:** SUB-SKILL OBRIGATÓRIA: use superpowers:subagent-driven-development (recomendado) ou superpowers:executing-plans para implementar este plano task por task. Os passos usam a sintaxe de checkbox (`- [ ]`) para acompanhamento.

**Objetivo:** Trocar as 3 salas desenhadas à mão (plano 2a) por um
gerador que monta um `building.Building` com N salas conectadas por
portas, sorteadas por um passeio aleatório a partir da sala de spawn —
sem nunca produzir uma porta sem sala do outro lado, nem uma sala
inalcançável a partir do spawn.

**Arquitetura:** Um pacote novo, `generator`, com uma única função que
consome as peças já existentes do plano 2a (`room.Room`,
`room.NewBlankRoomTiles`, `building.Building`, `building.Neighbor`) e não
muda nenhuma delas. O algoritmo é um passeio aleatório que cresce a
construção uma sala por vez a partir do spawn, sempre anexando a sala
nova a uma sala já colocada (o que garante conectividade por construção)
e sempre abrindo a porta recíproca nos dois lados ao mesmo tempo (o que
garante o invariante "toda porta tem sala do outro lado" sem precisar de
um passe de validação depois).

**Simplificação deliberada**: o roadmap fala em "sortear a partir de um
banco de salas-modelo" — isso sugere variedade de conteúdo de sala
(tamanhos, obstáculos diferentes). Esse plano só varia a topologia das
portas; o conteúdo de toda sala continua sendo o retângulo vazio do
plano 2a (`room.NewBlankRoomTiles`), porque ainda não existe nada pra
variar — sem foco de dengue (Limpeza) nem inimigo (Inimigos & bosses),
qualquer "modelo de sala" diferente seria decoração sem função. Variedade
de conteúdo real fica pra quando esses planos definirem o que uma sala
pode conter.

**Stack técnica:** Go, `github.com/hajimehoshi/ebiten/v2`,
`github.com/quasilyte/ebitengine-input`, `github.com/stretchr/testify`.

**Spec:** [docs/superpowers/specs/2026-10-17-dengue-boy-isaac-redesign-design.md](../specs/2026-10-17-dengue-boy-isaac-redesign-design.md)

## Restrições globais

- Stack inalterada. Nenhuma dependência nova (usa só `math/rand` da
  biblioteca padrão).
- Reusa a convenção de grade fixa 25×20 + porta em posição fixa por lado
  do plano 2a — nenhuma sala tem formato ou porta diferente disso.
- Todas as salas geradas nascem `Cleared: true`, pelo mesmo motivo do
  plano 2a: ainda não existe foco de dengue nem inimigo pra travar uma
  sala de verdade; o gancho de travar fica pros planos de Limpeza e
  Inimigos.
- `go build ./...` e `go test ./...` precisam passar no final do plano.

## Foco de revisão

- Gerar mais salas do que cabem sem colisão de coordenada (ex: pedir 50
  salas numa grade pequena) não pode travar num loop infinito — o
  gerador precisa parar e devolver o que já construiu quando não houver
  mais expansão possível — responsabilidade da Task 1.
- Toda porta gerada precisa ter sala do lado oposto (nunca uma porta
  "pra lugar nenhum") — responsabilidade da Task 1, já coberta
  explicitamente por teste.
- Toda sala gerada precisa ser alcançável a partir do spawn por portas
  — responsabilidade da Task 1, já coberta explicitamente por teste.
- Rodar o gerador duas vezes com o mesmo `*rand.Rand` (mesma seed)
  precisa produzir o mesmo resultado — determinismo é o que torna esse
  algoritmo testável; se a implementação usar alguma fonte de
  aleatoriedade fora do `rng` recebido (ex: `rand.Intn` do pacote global
  em vez do `rng.Intn`), os testes da Task 1 ficam não-reprodutíveis e
  quebram de forma intermitente — responsabilidade da Task 1.
- Pedir `roomCount = 1` (só o spawn, sem nenhuma porta) não pode
  quebrar — é o caso-base do algoritmo, não um erro — responsabilidade
  da Task 1.

---

### Task 1: Algoritmo de geração

**Arquivos:**
- Criar: `internal/modules/generator/generator.go`
- Criar: `internal/modules/generator/generator_test.go`

**Interfaces:**
- Consome: `room.Room`, `room.NewBlankRoomTiles`, `room.Side` + as 4
  constantes de lado, `room.Opposite` (plano 2a, Task 1); `building.Building`,
  `building.Coord`, `building.Neighbor` (plano 2a, Task 2);
  `maptile.MapTile`, `maptile.Map` (já existentes).
- Produz:
  - `func generator.Generate(roomCount int, rng *rand.Rand, images map[maptile.MapTile]*ebiten.Image) *building.Building`
  - Consumida pela Task 2 (`cmd/main.go`).

**Algoritmo** (a assinatura e os testes abaixo não determinam os
detalhes de busca/sorteio — a lógica de crescimento é descrita aqui):

Mantenha um mapa `doors map[building.Coord]map[room.Side]bool` (as
portas de cada sala, acumuladas por lado, por coordenada) e uma lista
`placed []building.Coord` começando só com `{0, 0}`. Repita até `len(placed)
== roomCount` ou até **200 tentativas sem sucesso** (cap de segurança
contra loop infinito quando a grade fica cheia de vizinhos ocupados):
sorteie (via `rng`) uma coordenada já em `placed` e um dos 4 lados;
calcule a coordenada vizinha com `building.Neighbor`; se essa vizinha já
está em `placed`, essa tentativa não conta como sucesso — tente de novo.
Senão: adicione a vizinha a `placed`, registre a porta sorteada na
coordenada de origem e a porta oposta (`room.Opposite`) na coordenada
vizinha, e zere o contador de tentativas sem sucesso. Com `roomCount ==
1`, o laço nunca encontra uma vizinha nova pra adicionar — ele bate o cap
de 200 tentativas e para, devolvendo só o spawn, sem porta nenhuma, o que
é o resultado correto.

Depois de montar `doors`, construa cada `*room.Room` com `Map:
&maptile.Map{Tiles: room.NewBlankRoomTiles(doors[coord]), Images:
images}`, `Doors: doors[coord]`, `Cleared: true`, monte o
`building.Building{Rooms: ..., Current: building.Coord{0, 0}}` e
retorne.

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package generator_test

import (
	"app/internal/modules/building"
	"app/internal/modules/generator"
	"app/internal/modules/room"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateProducesRequestedRoomCount(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(6, rng, nil)
	require.Len(t, b.Rooms, 6)
}

func TestGenerateStartsAtOriginAndSpawnIsCurrent(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(6, rng, nil)
	require.Contains(t, b.Rooms, building.Coord{X: 0, Y: 0})
	require.Equal(t, building.Coord{X: 0, Y: 0}, b.Current)
}

func TestGenerateSingleRoomDoesNotHang(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(1, rng, nil)
	require.Len(t, b.Rooms, 1)
}

func TestGenerateAllRoomsAreClearedByDefault(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(6, rng, nil)
	for _, r := range b.Rooms {
		require.True(t, r.Cleared)
	}
}

func allSides() []room.Side {
	return []room.Side{room.North, room.South, room.East, room.West}
}

func TestGenerateDoorsAreReciprocal(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(10, rng, nil)
	for coord, r := range b.Rooms {
		for _, side := range allSides() {
			if !r.HasDoor(side) {
				continue
			}
			neighborCoord := building.Neighbor(coord, side)
			neighbor, exists := b.Rooms[neighborCoord]
			require.True(t, exists, "room at %+v has a door %v with no neighbor room", coord, side)
			require.True(t, neighbor.HasDoor(room.Opposite(side)), "neighbor at %+v does not reciprocate the door", neighborCoord)
		}
	}
}

func TestGenerateAllRoomsReachableFromSpawn(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(10, rng, nil)

	visited := map[building.Coord]bool{{X: 0, Y: 0}: true}
	queue := []building.Coord{{X: 0, Y: 0}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		r := b.Rooms[current]
		for _, side := range allSides() {
			if !r.HasDoor(side) {
				continue
			}
			next := building.Neighbor(current, side)
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	require.Len(t, visited, len(b.Rooms), "every generated room must be reachable from spawn")
}

func TestGenerateIsDeterministicForTheSameSeed(t *testing.T) {
	a := generator.Generate(10, rand.New(rand.NewSource(42)), nil)
	b := generator.Generate(10, rand.New(rand.NewSource(42)), nil)
	require.Equal(t, a.Rooms, b.Rooms)
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/generator/... -v`
Esperado: FAIL — pacote `generator` ainda não existe.

- [ ] **Passo 3: Implementar `Generate` em `internal/modules/generator/generator.go`** seguindo o algoritmo descrito acima.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/generator/... -v`
Esperado: PASS (7 testes)

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/generator/generator.go internal/modules/generator/generator_test.go
git commit -m "Add procedural room-graph generator"
```

---

### Task 2: Trocar as salas fixas pelo gerador em `cmd/main.go`

**Arquivos:**
- Modificar: `cmd/main.go`

**Interfaces:**
- Consome: `generator.Generate` (Task 1).

- [ ] **Passo 1: Trocar o building de 3 salas fixas (plano 2a) pelo gerador**

Carregar as imagens de tile uma vez com `maptile.New()` (só pra
reaproveitar `.Images` — o `.Tiles` desse retorno é descartado), e usar
`generator.Generate(6, rand.New(rand.NewSource(time.Now().UnixNano())),
images)` para montar o `building.Building` passado pro resto de
`initializeGame`. `6` é um tamanho de demonstração razoável — a
quantidade real de salas por season é decisão de um plano futuro
(Inimigos & bosses).

- [ ] **Passo 2: Verificar que o projeto inteiro compila e os testes passam**

Rodar: `go build ./...`
Esperado: sem erros.

Rodar: `go test ./...`
Esperado: todos os pacotes passam.

- [ ] **Passo 3: Verificação manual**

Rodar: `go run cmd/main.go` (ou `make run`). Confirmar que o jogo abre
numa sala de spawn com pelo menos uma porta, que dá pra explorar várias
salas diferentes andando pelas portas, e que nenhuma porta leva a uma
tela vazia/trava o jogo.

- [ ] **Passo 4: Commit**

```bash
git add cmd/main.go
git commit -m "Generate the building procedurally instead of 3 fixed rooms"
```
