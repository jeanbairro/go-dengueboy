# Plano de Implementação: Movimento Livre & Fundação de Entidades

> **Para workers agênticos:** SUB-SKILL OBRIGATÓRIA: use superpowers:subagent-driven-development (recomendado) ou superpowers:executing-plans para implementar este plano task por task. Os passos usam a sintaxe de checkbox (`- [ ]`) para acompanhamento.

**Objetivo:** Substituir o movimento atual do jogador, travado em grid, por movimento livre e contínuo em 360° (com normalização diagonal e colisão com deslize nas paredes), e introduzir uma abstração genérica de `Entity` para que o game loop possa atualizar/desenhar o jogador — e, depois, inimigos e objetos interativos — de forma uniforme.

**Arquitetura:** Três pacotes pequenos e testáveis isoladamente mudam: `geom` ganha uma função pura de matemática vetorial pra normalização diagonal, um novo pacote `entity` define a interface `Entity` e um `Manager` que controla o loop de update/draw, e `player` é reescrito pra guardar sua posição em espaço de pixel (não mais em tile), se mover continuamente, e satisfazer `entity.Entity`. O `game.Game` passa então a conduzir o jogador através de um `entity.Manager`, em vez de chamá-lo diretamente. Uma câmera por sala fica explicitamente **fora de escopo** aqui: o mapa único atual tem exatamente o tamanho da tela (25×20 tiles = 800×640px), então uma câmera não teria efeito visível nenhum até o sistema de salas (um plano futuro) introduzir salas menores que a viewport. Construir isso agora significaria código não testado e não integrado; fica postergado pra aquele plano.

**Por que `Entity`/`Manager`, já que hoje só existe o `Player`?** Hoje o `Game` só tem uma "coisa com vida própria" na tela, então `g.Player.Update()` / `g.Player.Draw(screen)` chamados por nome funcionam bem. Mas o spec já prevê dezenas dessas coisas chegando nos próximos planos: mosquito macho, mosquito fêmea, larvas do enxame, o chefão, pneu infestado, vaso, poça — cada uma, todo frame, precisa se atualizar e se desenhar. Sem uma abstração comum, `Game.Update()`/`Draw()` teria que crescer uma linha a mais pra cada inimigo/objeto novo (`g.Mosquito1.Update()`, `g.Mosquito2.Update()`, ...), e o `Game` precisaria ser editado toda vez que um boss ou objeto novo nascesse. `Entity` é só o contrato "sei me atualizar (`Update()`) e sei me desenhar (`Draw(screen)`)"; `Manager` é a lista que guarda tudo que cumpre esse contrato e chama `Update()`/`Draw()` em lote, sem saber (nem precisar saber) o que cada coisa é por dentro. Construir isso agora, mesmo só com o `Player` dentro da lista, é o que permite que os próximos planos (objetos de limpeza, inimigos, bosses) só precisem dar um `Entities.Add(nova_coisa)` — sem nunca mais tocar em `game.go`.

**Stack técnica:** Go, `github.com/hajimehoshi/ebiten/v2`, `github.com/quasilyte/ebitengine-input`, `github.com/stretchr/testify`.

**Spec:** [docs/superpowers/specs/2026-10-17-dengue-boy-isaac-redesign-design.md](../specs/2026-10-17-dengue-boy-isaac-redesign-design.md)

## Restrições globais

- A stack permanece como está: Go + `github.com/hajimehoshi/ebiten/v2` + `github.com/quasilyte/ebitengine-input`. Nenhuma dependência nova.
- O projeto precisa continuar buildando: `go build ./...` tem que funcionar ao final do plano.
- O movimento precisa ser livre e contínuo (360°, com suporte a diagonal) — nenhum "encaixe" em tile/grid em código novo.

**Uma nota sobre a ordem do build:** a Task 3 muda a assinatura de `Player.Update()`, que `game.go` ainda chama com a assinatura antiga até a Task 4. `go build ./...` e `go test ./...` vão falhar entre a Task 3 e a Task 4 — isso é esperado. Rode a verificação de cada task restrita ao pacote que ela toca (ex: `go test ./internal/modules/player/...`) até o passo final da Task 4, que restaura um `go build ./...` completo e verde.

## Foco de revisão

- Teclas de direção opostas pressionadas ao mesmo tempo (Esquerda+Direita, ou Cima+Baixo) precisam cancelar o movimento líquido naquele eixo para zero, sem tremedeira nem escolher uma arbitrariamente — responsabilidade da Task 3 (`movementIntent`).
- Um movimento diagonal bloqueado em *ambos* os eixos por uma parede precisa deixar a posição do jogador totalmente inalterada nos dois eixos, nunca atravessando nenhuma das paredes — responsabilidade da Task 3.
- Um movimento diagonal bloqueado em apenas *um* eixo ainda precisa permitir que o jogador deslize pelo eixo livre, respeitando a fluidez que o design pede, em vez de parar seco — responsabilidade da Task 3.
- Sem nenhuma direção pressionada, o jogador precisa se estabilizar num sprite parado legível, não congelar no último frame de caminhada — responsabilidade da Task 3.
- Um `entity.Manager` sem nada registrado ainda (verdade pra qualquer cena antes do primeiro inimigo/objeto existir) precisa ter `Update()`/`Draw()` como no-ops seguros, sem panic — responsabilidade da Task 2.

---

### Task 1: Matemática do vetor de movimento diagonal

**Arquivos:**
- Modificar: `internal/modules/geom/geom.go`
- Teste: `internal/modules/geom/geom_test.go` (criar)

**Interfaces:**
- Produz: `func NormalizeDiagonal(dx, dy, speed float64) Position` — função pura; dada a intenção bruta por eixo (cada uma esperada como -1, 0 ou 1) e uma magnitude de velocidade desejada, retorna o vetor de movimento escalado, normalizado para que o movimento diagonal não fique mais rápido que o movimento nos eixos cardeais. Consumida pela Task 3 (pacote `player`).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package geom_test

import (
	"app/internal/modules/geom"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeDiagonalCardinalUnchanged(t *testing.T) {
	result := geom.NormalizeDiagonal(1, 0, 2.56)
	require.Equal(t, geom.Position{X: 2.56, Y: 0}, result)
}

func TestNormalizeDiagonalScalesDownDiagonalMovement(t *testing.T) {
	result := geom.NormalizeDiagonal(1, -1, 2.56)
	require.InDelta(t, 1.8101933598, result.X, 0.0001)
	require.InDelta(t, -1.8101933598, result.Y, 0.0001)
}

func TestNormalizeDiagonalOpposingInputCancelsToZero(t *testing.T) {
	result := geom.NormalizeDiagonal(0, 0, 2.56)
	require.Equal(t, geom.Position{X: 0, Y: 0}, result)
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/geom/... -run TestNormalizeDiagonal -v`
Esperado: FAIL com "undefined: geom.NormalizeDiagonal"

- [ ] **Passo 3: Implementar `NormalizeDiagonal(dx, dy, speed float64) Position` em `internal/modules/geom/geom.go`**

Quando `dx` e `dy` são ambos diferentes de zero, escalar os dois por `1/√2` (`0.7071067811865476`) antes de multiplicar por `speed`, pra que o movimento diagonal tenha a mesma magnitude do movimento cardeal. Quando só um eixo (ou nenhum) é diferente de zero, não precisa de escala nenhuma.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/geom/... -run TestNormalizeDiagonal -v`
Esperado: PASS (3 testes)

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/geom/geom.go internal/modules/geom/geom_test.go
git commit -m "feat: add diagonal-normalized movement vector helper"
```

---

### Task 2: Interface genérica `Entity` e `Manager`

**Arquivos:**
- Criar: `internal/modules/entity/entity.go`
- Teste: `internal/modules/entity/manager_test.go` (criar)

**Interfaces:**
- Produz:
  - `type Entity interface { Update(); Draw(screen *ebiten.Image) }`
  - `type Manager struct { ... }`
  - `func NewManager() *Manager`
  - `func (m *Manager) Add(e Entity)`
  - `func (m *Manager) Update()` — chama `Update()` em cada entidade registrada, na ordem de registro.
  - `func (m *Manager) Draw(screen *ebiten.Image)` — chama `Draw(screen)` em cada entidade registrada, na ordem de registro.
  - `func (m *Manager) Len() int` — número de entidades registradas.
  - Consumida pela Task 3 (`player.Player` vai satisfazer `Entity`) e pela Task 4 (`game.Game` conecta o jogador através de um `Manager`).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package entity_test

import (
	"app/internal/modules/entity"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/require"
)

type fakeEntity struct {
	updateCalls int
	drawCalls   int
}

func (f *fakeEntity) Update()                        { f.updateCalls++ }
func (f *fakeEntity) Draw(screen *ebiten.Image) { f.drawCalls++ }

func TestManagerUpdateCallsEveryRegisteredEntity(t *testing.T) {
	m := entity.NewManager()
	a, b := &fakeEntity{}, &fakeEntity{}
	m.Add(a)
	m.Add(b)

	m.Update()

	require.Equal(t, 1, a.updateCalls)
	require.Equal(t, 1, b.updateCalls)
}

func TestManagerDrawCallsEveryRegisteredEntity(t *testing.T) {
	m := entity.NewManager()
	a, b := &fakeEntity{}, &fakeEntity{}
	m.Add(a)
	m.Add(b)

	m.Draw(nil)

	require.Equal(t, 1, a.drawCalls)
	require.Equal(t, 1, b.drawCalls)
}

func TestManagerWithNoEntitiesIsNoop(t *testing.T) {
	m := entity.NewManager()

	require.NotPanics(t, func() {
		m.Update()
		m.Draw(nil)
	})
	require.Equal(t, 0, m.Len())
}

func TestManagerLenReflectsRegistrations(t *testing.T) {
	m := entity.NewManager()
	m.Add(&fakeEntity{})
	m.Add(&fakeEntity{})

	require.Equal(t, 2, m.Len())
}
```

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/entity/... -v`
Esperado: FAIL com "no Go files in..." ou "undefined: entity.NewManager" (o pacote ainda não existe)

- [ ] **Passo 3: Implementar `Entity`, `Manager`, `NewManager`, `Add`, `Update`, `Draw`, `Len` em `internal/modules/entity/entity.go`**

`Manager` guarda um slice `[]Entity` não-exportado; cada método é um loop simples ou um `len()`.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/entity/... -v`
Esperado: PASS (4 testes)

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/entity/entity.go internal/modules/entity/manager_test.go
git commit -m "feat: add generic Entity interface and Manager"
```

---

### Task 3: Movimento livre contínuo no `Player`

**Arquivos:**
- Modificar: `internal/modules/player/player.go` (constantes, imports, `Update`, novo `Move`, novo `updateFacing`, novo `readMovementVector`, novo `movementIntent`, `toObject`)
- Modificar: `internal/modules/player/renderer.go:21-25` (`Draw` — parar de multiplicar a posição por `spriteSize`, já que a posição agora é em espaço de pixel)
- Teste: `internal/modules/player/player_test.go` (substituir inteiro — o `TestPlayerMovement` atual já falha hoje por causa de um caminho de asset relativo e testa o contrato antigo baseado em grid; fica totalmente substituído)
- Teste: `internal/modules/player/movement_test.go` (criar, white-box — `package player`, não `player_test`, pra poder exercitar o `movementIntent` não-exportado diretamente)

**Interfaces:**
- Consome: `geom.NormalizeDiagonal(dx, dy, speed float64) geom.Position` (Task 1); interface `entity.Entity` (Task 2, pra uma asserção em tempo de compilação).
- Produz:
  - `func (p *Player) Update()` — **mudança que quebra compatibilidade**: remove o parâmetro `mapTile *maptile.Map`, que não era usado.
  - `func (p *Player) Move(movement geom.Position)` — novo, exportado. Aplica `movement` em `p.Position` um eixo de cada vez (X depois Y), checando `p.CollisionSystem.HasCollision` independentemente por eixo, pra que o jogador deslize ao longo de uma parede em vez de parar seco quando só um eixo está bloqueado. Sempre atualiza `p.PreviousPosition` pra posição anterior ao movimento e `p.CurrentAction` pra refletir a direção implícita em `movement` (horizontal tem prioridade sobre vertical quando os dois eixos são diferentes de zero, já que o spritesheet só tem 4 linhas cardeais), independente do movimento ter sido bloqueado ou não.
  - `movementIntent(left, right, up, down bool) (dx, dy float64)` — novo, não-exportado, puro. Transforma os 4 booleanos de direção pressionada em intenção bruta por eixo em `{-1, 0, 1}`, somando direções opostas pra `0` quando as duas estão pressionadas.
  - `Player` agora satisfaz `entity.Entity` (asserção em tempo de compilação: `var _ entity.Entity = (*Player)(nil)`).
  - Consumida pela Task 4 (`game.Game`).

- [ ] **Passo 1: Escrever os testes que vão falhar**

```go
package player_test

import (
	"app/internal/modules/collision"
	"app/internal/modules/geom"
	"app/internal/modules/player"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayerMoveFreelyWithNoCollision(t *testing.T) {
	p := &player.Player{
		Position:        geom.Position{X: 100, Y: 100},
		CollisionSystem: &collision.CollisionSystem{},
	}

	p.Move(geom.Position{X: 2.56, Y: 0})

	require.Equal(t, geom.Position{X: 102.56, Y: 100}, p.Position)
	require.Equal(t, geom.Position{X: 100, Y: 100}, p.PreviousPosition)
	require.Equal(t, player.MoveRight, p.CurrentAction)
}

func TestPlayerMoveSlidesAlongWallWhenOneAxisBlocked(t *testing.T) {
	wall := collision.Object{Position: geom.Position{X: 132, Y: 100}, Width: 32, Height: 32}
	p := &player.Player{
		Position:        geom.Position{X: 100, Y: 100},
		CollisionSystem: &collision.CollisionSystem{Objects: []collision.Object{wall}},
	}

	p.Move(geom.Position{X: 5, Y: 5})

	require.Equal(t, geom.Position{X: 100, Y: 105}, p.Position)
}

func TestPlayerMoveBlockedOnBothAxesStaysInPlace(t *testing.T) {
	wall := collision.Object{Position: geom.Position{X: 100, Y: 116}, Width: 48, Height: 48}
	p := &player.Player{
		Position:        geom.Position{X: 100, Y: 100},
		CollisionSystem: &collision.CollisionSystem{Objects: []collision.Object{wall}},
	}

	p.Move(geom.Position{X: 5, Y: 5})

	require.Equal(t, geom.Position{X: 100, Y: 100}, p.Position)
}

func TestPlayerMoveZeroVectorSetsIdleFacing(t *testing.T) {
	p := &player.Player{
		Position:        geom.Position{X: 50, Y: 50},
		CurrentAction:   player.MoveRight,
		CollisionSystem: &collision.CollisionSystem{},
	}

	p.Move(geom.Position{X: 0, Y: 0})
	p.Move(geom.Position{X: 0, Y: 0})

	require.Equal(t, geom.Position{X: 50, Y: 50}, p.Position)
	require.Equal(t, player.NoAction, p.CurrentAction)
}
```

E, em `internal/modules/player/movement_test.go` (white-box, `package player`):

```go
package player

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovementIntentSingleDirection(t *testing.T) {
	dx, dy := movementIntent(true, false, false, false)
	require.Equal(t, -1.0, dx)
	require.Equal(t, 0.0, dy)
}

func TestMovementIntentOpposingHorizontalCancels(t *testing.T) {
	dx, dy := movementIntent(true, true, false, false)
	require.Equal(t, 0.0, dx)
	require.Equal(t, 0.0, dy)
}

func TestMovementIntentOpposingVerticalCancels(t *testing.T) {
	dx, dy := movementIntent(false, false, true, true)
	require.Equal(t, 0.0, dx)
	require.Equal(t, 0.0, dy)
}

func TestMovementIntentDiagonal(t *testing.T) {
	dx, dy := movementIntent(false, true, true, false)
	require.Equal(t, 1.0, dx)
	require.Equal(t, -1.0, dy)
}
```

(ordem `left, right, up, down`; `up`/`down` mapeiam pra `dy` com `up = -1`, seguindo a convenção já existente de `MoveUp` diminuir Y.)

- [ ] **Passo 2: Rodar os testes e confirmar que falham**

Rodar: `go test ./internal/modules/player/... -v`
Esperado: FAIL com "undefined: (*player.Player).Move" / "undefined: movementIntent" (ou erro de compilação parecido)

- [ ] **Passo 3: Implementar o novo modelo de movimento**

Em `internal/modules/player/player.go`:
- Mudar `initialPositionX`/`initialPositionY` de `2`/`2` (unidades de tile) pra `64`/`64` (unidades de pixel — o mesmo ponto de spawn, `2 * spriteSize`).
- Mudar `speed` de `.08` (tiles/frame) pra `2.56` (pixels/frame — o equivalente exato, `.08 * 32`).
- Remover o import de `maptile` (não é mais usado uma vez que o parâmetro de `Update` é removido).
- Substituir `getWantedPosition`/`setCurrentAction` por:
  - `movementIntent(left, right, up, down bool) (dx, dy float64)`: função pura, sem receiver. `dx` é `-1` só-esquerda, `+1` só-direita, `0` se nenhuma ou as duas estiverem pressionadas; `dy` é `-1` só-cima, `+1` só-baixo, `0` se nenhuma ou as duas estiverem pressionadas.
  - `readMovementVector() geom.Position`: lê as 4 direções de `p.InputHandler.ActionIsPressed`, chama `movementIntent(...)`, e retorna `geom.NormalizeDiagonal(dx, dy, speed)`.
  - `updateFacing(movement geom.Position)`: define `p.CurrentAction` a partir do sinal de `movement` — `MoveLeft`/`MoveRight` tem prioridade sobre `MoveUp`/`MoveDown` quando os dois eixos são diferentes de zero; `NoAction` quando os dois são zero.
- Adicionar `Move(movement geom.Position)` como descrito acima em Interfaces, chamando `updateFacing` primeiro, depois resolvendo X e então Y independentemente contra `p.CollisionSystem.HasCollision(toObject(...))`.
- Reescrever `Update()` (sem parâmetros) como: calcular `movement := p.readMovementVector()`, chamar `p.Move(movement)`, chamar `p.setSprite()`.
- Atualizar `toObject(position geom.Position) collision.Object` pra usar `position` direto (sem `* spriteSize`), já que `Position` agora já é espaço de pixel.
- Adicionar `var _ entity.Entity = (*Player)(nil)` (import novo: `app/internal/modules/entity`).

Em `internal/modules/player/renderer.go:21-25`: trocar `op.GeoM.Translate(float64(p.Position.X*spriteSize), float64(p.Position.Y*spriteSize))` por `op.GeoM.Translate(p.Position.X, p.Position.Y)`.

- [ ] **Passo 4: Rodar os testes e confirmar que passam**

Rodar: `go test ./internal/modules/player/... -v`
Esperado: PASS (8 testes: 4 em `player_test.go`, 4 em `movement_test.go`). Nota: `go build ./...` ainda vai falhar nesse ponto — `game.go` chama a assinatura antiga de `Update(mapTile)` até a Task 4.

- [ ] **Passo 5: Commit**

```bash
git add internal/modules/player/player.go internal/modules/player/renderer.go internal/modules/player/player_test.go internal/modules/player/movement_test.go
git commit -m "feat: free continuous movement with diagonal normalization and wall sliding"
```

---

### Task 4: Conectar o jogador ao `Game` através do `entity.Manager`

**Arquivos:**
- Modificar: `internal/modules/game/game.go` (struct, `New`, `Update`, `Draw`)
- Teste: `internal/modules/game/game_test.go` (criar)
- Sem mudanças necessárias: `cmd/main.go` — as assinaturas exportadas de `player.New(...)` e `game.New(...)` não mudam.

**Interfaces:**
- Consome: `entity.NewManager() *Manager`, `(*Manager).Add`, `.Update()`, `.Draw(screen)`, `.Len()` (Task 2); `Player` satisfazendo `entity.Entity` (Task 3).
- Produz: `Game.Entities *entity.Manager` (campo novo, exportado). `Game.Player` continua existindo (um plano futuro de câmera vai precisar de acesso direto a ele).

- [ ] **Passo 1: Escrever o teste que vai falhar**

```go
package game_test

import (
	"app/internal/modules/collision"
	"app/internal/modules/game"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"testing"

	input "github.com/quasilyte/ebitengine-input"
	"github.com/stretchr/testify/require"
)

func TestNewRegistersPlayerInEntities(t *testing.T) {
	p := &player.Player{
		Position:        geom.Position{X: 0, Y: 0},
		CollisionSystem: &collision.CollisionSystem{},
	}

	g := game.New(&maptile.Map{}, p, input.System{})

	require.Equal(t, 1, g.Entities.Len())
	require.Equal(t, p, g.Player)
}
```

- [ ] **Passo 2: Rodar o teste e confirmar que falha**

Rodar: `go test ./internal/modules/game/... -v`
Esperado: FAIL — erro de compilação (`g.Entities` indefinido), já que `Game` ainda não tem o campo `Entities`.

- [ ] **Passo 3: Adicionar `Entities *entity.Manager` ao `Game` e conectar em `internal/modules/game/game.go`**

- Adicionar o campo `Entities *entity.Manager` à struct `Game` (import novo: `app/internal/modules/entity`).
- Em `New(maptile *maptile.Map, player *player.Player, inputSystem input.System) *Game`: construir `entities := entity.NewManager()`, chamar `entities.Add(player)`, e definir isso no `Game` retornado.
- Reescrever `Update()` pra chamar `g.Entities.Update()` em vez de `g.updatePlayer()`; remover o método `updatePlayer`, agora não usado.
- Reescrever `Draw(screen)` pra chamar `g.Entities.Draw(screen)` em vez de `g.Player.Draw(screen)`.

- [ ] **Passo 4: Rodar o teste e confirmar que passa**

Rodar: `go test ./internal/modules/game/... -v`
Esperado: PASS

- [ ] **Passo 5: Confirmar que o build completo está verde de novo**

Rodar: `go build ./...`
Esperado: sem erros (esse é o primeiro ponto, desde a Task 3, onde o projeto inteiro compila).

Rodar: `go test ./...`
Esperado: todos os pacotes passam.

- [ ] **Passo 6: Verificação manual**

Rodar: `go run cmd/main.go` (ou `make run`). Confirmar que o personagem ainda se move e anima corretamente, agora conduzido através de `Entities.Update()`/`Draw()` — e confirmar que o movimento agora é contínuo/diagonal em vez de travado em grid, incluindo o deslize suave ao longo de uma parede quando o movimento é feito em ângulo contra ela.

- [ ] **Passo 7: Commit**

```bash
git add internal/modules/game/game.go internal/modules/game/game_test.go
git commit -m "feat: drive player through entity.Manager in Game"
```
