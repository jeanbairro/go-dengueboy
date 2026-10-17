# Dengue Boy — Redesign estilo Isaac (Rebirth)

Status: aprovado pelo autor para virar plano de implementação.

## Intenção e porquê

O jogo hoje (Go + Ebitengine v2) é um protótipo mínimo: um boneco anda
livremente por um mapa fixo 25x20 com colisão contra parede. Não existe
objetivo, inimigo, dano, progressão, ou interação com o cenário. Os tiles
`WaterTire`/`WaterVase`/`EmptyTire` já estão declarados no enum de mapa,
mas não são usados em nenhum lugar — nem renderizados, nem na colisão. A
ação `Interact` já está mapeada no input (tecla Espaço/Enter/botão A) mas
nunca é lida.

O objetivo deste redesign é transformar isso num jogo de verdade, na
pegada de **The Binding of Isaac: Rebirth** (~80% da intensidade visual e
de tom), onde o jogador limpa focos de dengue dentro de uma
construção/prédio que vai se revelando sala por sala, enfrenta um chefe
temático a cada "season", e evolui com poderes ao longo de cada run. A
mensagem de prevenção à dengue deve emergir inteiramente do ato de jogar
— nunca de texto explicativo.

## Objetivos

- **Loop principal**: o jogador se move livremente (360°, sem grid) por
  salas de uma construção/obra abandonada. Cada sala tem focos de dengue
  (pneu com água, vaso com prato cheio, poça parada, garrafa virada pra
  cima, calha entupida). Interagir (`Interact`) limpa o foco. Sala 100%
  limpa destranca a(s) porta(s) para a(s) próxima(s) sala(s).
- **Progressão por season**, cada uma terminando num chefe:
  1. Larva gigante — lenta, tutorial de combate, corpo segmentado
     (estilo "minhoca gigante").
  2. Mosquito macho — rápido, mas tematicamente não pica (não causa dano
     de picada — é obstáculo/evasão, não ameaça direta).
  3. Mosquito fêmea — ataques à distância/picada, dano real + efeito de
     status.
  4. Enxame de larvas — várias larvas pequenas simultâneas na mesma
     arena.
  5. Chefão final — combina padrões das seasons anteriores.
- **Chefes e inimigos "bem feitos"**: padrões de projétil legíveis e
  esquiváveis (referência explícita do autor: "bichos que atiravam
  bolinhas"), inimigos grandes/segmentados com possíveis pontos fracos
  (referência: "minhocas gigantes").
- **Evolução do jogador**: o personagem ganha poderes ativos ao longo de
  uma season (estilo itens do Isaac), não só upgrades passivos de stat.
- **Morte reinicia só a season atual** — o jogador perde os poderes
  coletados *naquela* season, mas o status "season concluída" das
  seasons anteriores persiste pelo resto daquela sessão de jogo; não há
  reinício do jogo inteiro.
- **Tom visual/narrativo**: ~80% da intensidade horror-cute do Isaac —
  pode ser pesado (bichos grotescos, chefes ameaçadores), mas o
  personagem principal e o cenário continuam legíveis, não é gore
  gratuito.
- **Movimento livre 360°** — mudança técnica relevante em relação ao
  código atual, que é "anda célula por célula".
- **Ataque à distância com mira independente do movimento** — igual as
  lágrimas do Isaac: o jogador anda numa direção e mira/atira noutra ao
  mesmo tempo. Tematizado como um jato de repelente/inseticida em vez de
  lágrima. Os poderes coletados modificam esse ataque (dano, projéteis
  extras, efeitos).
- **Geração procedural de salas por season** — igual Isaac: cada season
  sorteia sua disposição de salas a partir de um conjunto de salas-
  modelo pré-desenhadas. Reiniciar a season (por morte) gera uma
  disposição nova, não repete o layout anterior.
- **Primeira leva de sprites placeholder**: gerar uma primeira versão de
  pixel art (personagem, inimigos, objetos de foco, tileset) na pegada
  Isaac para o autor refinar depois — ele mesmo vai ajustar/finalizar os
  assets (ex: no Aseprite).

## Não-objetivos / anti-objetivos (nesta iteração)

- **Nenhum texto educativo, dica, pop-up ou tela de tutorial explicando
  o tema de prevenção da dengue.** Isso foi confirmado duas vezes pelo
  autor como anti-objetivo crítico: a mensagem deve vir só da
  experiência de jogar. Intro curta, direto para o gameplay.
- Não construir uma cidade/overworld conectando vários prédios nesta
  iteração — escopo é uma única construção por enquanto.
- Sem multiplayer/co-op.
- Sem mudança de plataforma — continua desktop via Ebitengine/Go.
- **Critério de falha mesmo que tecnicamente funcione**: se o jogador
  precisar ler qualquer texto explicativo para entender o recado de
  saúde pública do jogo, o redesign falhou no que importa, mesmo que
  todas as mecânicas estejam funcionando.

## Restrições

- Stack atual deve ser mantida: Go + `github.com/hajimehoshi/ebiten/v2` +
  `github.com/quasilyte/ebitengine-input`.
- Pipeline de assets atual (spritesheet PNG fatiado via
  `ebitenutil.NewImageFromFile` + `SubImage`) pode evoluir, mas não há
  motivo para trocar de engine/framework.
- O build precisa continuar funcionando (`make run` / `go build ./...`).

## Decisões já tomadas (pelo autor, durante a conversa)

- Movimento livre (contínuo, 360°), substituindo o movimento atual
  travado em célula.
- Morte reinicia **apenas a season atual**; poderes daquela season são
  perdidos; seasons já batidas continuam marcadas como concluídas,
  permanentemente (não é um full-reset estilo Isaac raiz).
- Ordem da escada de chefes: Larva gigante → Mosquito macho (não pica) →
  Mosquito fêmea (pica, dano real) → Enxame de larvas → Chefão final.
- Padrão de design de chefe/inimigo: projéteis legíveis e esquiváveis +
  inimigos grandes/segmentados com pontos fracos.
- Tom: ~80% da intensidade horror-cute do Isaac Rebirth — pesado mas
  legível, sem recair em gore gratuito.
- Zero texto educativo na tela, em qualquer momento do jogo — anti-
  objetivo inegociável.
- Intro curta, sem tutorial longo, o jogador cai rápido para andar.
- Sprites: autor vai refinar; ele topa uma primeira passada de
  placeholders gerados para iterar em cima.
- Referências explícitas do Isaac citadas pelo autor como objetivo de
  qualidade: fluidez de movimento, tensão do permadeath por run, vilões
  bem desenhados (minhocas gigantes, atiradores de bolinha), evolução de
  poderes dentro da run, qualidade geral dos sprites.
- Ataque à distância com mira separada do movimento (não corpo a corpo),
  seguindo minha sugestão, confirmada pelo autor.
- Geração procedural de salas por season a partir de um banco de salas-
  modelo, seguindo minha sugestão, confirmada pelo autor.

## O que fica a cargo de quem construir (decisões minhas, meramente
sugeridas — não confirmadas linha a linha pelo autor)

- **Implementação do sistema de salas/portas/minimapa**: a abordagem
  (geração procedural por season) já está decidida acima; os detalhes
  ficam a cargo de quem construir. `maptile.Map` hoje é um mapa único
  hardcoded — proponho uma abstração `Room` reutilizável + um
  gerenciador de construção (`Building`/`Floor`) que guarda o grafo de
  salas sorteado, conexões, estado de porta trancada/destrancada, e um
  banco inicial pequeno de salas-modelo (pode começar com poucas,
  crescer depois), com minimapa simples no HUD.
- **Esquema de mira**: a mecânica (ataque à distância, mira independente
  do movimento) já está decidida acima; o mapeamento exato de teclas/
  botões (ex: setas para mirar no teclado, analógico direito no
  controle) fica a cargo de quem construir, seguindo o padrão do
  próprio Isaac no PC+controle.
- **Câmera**: hoje não existe nenhuma. Sugiro câmera fixa por sala
  (estilo Isaac — cada sala cabe inteira na tela) em vez de câmera livre
  seguindo o jogador, por ser mais simples e combinar com o sistema de
  salas. Decisão reversível depois.
- **Sistema de entidades**: hoje só `Player` tem `Update`/`Draw`.
  Precisa de uma abstração genérica (ou ECS leve) para inimigos, objetos
  interativos e projéteis, sem duplicar a lógica do `Player`.
- **Lista exata de focos de dengue e mecânica de interação** (segurar
  botão vs. ação instantânea; lista além de pneu/vaso — poça, garrafa,
  calha): segue a proposta das notas anteriores, não foi reconfirmada
  item a item pelo autor; tratar como default razoável, ajustável.
- **Dados de sala/mapa**: considerar migrar de literais Go hardcoded
  para arquivos de dados (JSON ou similar) assim que existir mais de
  1-2 salas — timing fica a critério de quem implementar.
- **Lista exata de poderes/habilidades ativas**: o autor pediu "evoluir
  com poderes" mas não especificou quais. Sugiro começar com um conjunto
  mínimo (2-3 poderes ativos + alguns itens passivos de stat), seguindo
  o padrão de itens do Isaac, e expandir depois.
- **Progressão meta-persistente além do flag de "season concluída"**
  (ex: desbloqueios permanentes entre playthroughs completos) não foi
  discutida — fora de escopo a menos que o autor levante isso de novo.
- **Escopo da persistência "season concluída"**: para esta v1, esse
  estado vive só em memória, durante a sessão de jogo em execução — não
  há save em disco. Se o jogo for fechado, o progresso de seasons
  concluídas se perde. Um sistema de save persistente em disco fica como
  trabalho futuro, caso o autor queira depois; não é necessário para o
  que foi descrito nesta conversa.
- **Efeito de status da picada da fêmea**: lentidão temporária de
  movimento por alguns segundos (ex: "coceira"), implementada como um
  framework leve e genérico de status effect — não um hack específico —
  para poder ser reaproveitado por chefes futuros (enxame, chefão
  final) sem redesenho.
- **Ponto de partida no código**: a mudança para movimento livre e o
  sistema de salas são mudanças que quebram a estrutura atual de
  `internal/modules/player`, `maptile`, `collision` e `game`. Isso deve
  ser planejado como trabalho em branch nova, como o autor já pretende.

## Estado atual do código (referência)

- `internal/modules/player`: posição, sprites direcionais (4 direções x
  5 frames), movimento em grid, colisão AABB contra parede.
- `internal/modules/maptile`: mapa único 25x20 hardcoded em
  `loadMap()`; só os tiles `Empty` e `Wall` são usados de fato.
- `internal/modules/collision`: AABB simples, só contra `Wall`.
- `internal/modules/game`: loop `Update`/`Draw`/`Layout` sem conceito de
  estados (menu, jogando, boss, game over).
- `cmd/main.go`: registra a ação `Interact` no keymap, mas ela nunca é
  lida em lugar nenhum — gancho pronto, não implementado.
- Sem câmera, sem HUD, sem áudio, sem inimigos, sem sistema de vida/dano,
  sem save.
- Assets: `char.png` (personagem, proporções RPG Maker genéricas, não
  "Isaac"), `fly.png` (mosquitos com silhueta fraca), `jungle.png`
  (tileset de selva, tematicamente errado — tema é quintal/obra urbana),
  `objects.png` (pneus e vasos, sem estados visuais
  infestado/limpo).

## Roadmap de implementação

Esse spec cobre o redesign inteiro, mas o trabalho é grande demais pra um
plano só — é construído em fatias, cada uma com seu próprio plano de
implementação (skill `superpowers:writing-plans`) e seu próprio PR contra
a branch de integração `feature/dengue-redesign`. Ordem definida pela
dependência técnica entre as peças, não só pela ordem narrativa do jogo:

1. ✅ **Fundação** — movimento livre 360°, sistema de entidade genérico
   (`entity.Entity`/`Manager`). Plano:
   [2026-10-17-free-movement-entity-foundation.md](../plans/2026-10-17-free-movement-entity-foundation.md).
   PR #1.
2. ⬜ **Salas & construção** — `Room`/`Building`, geração procedural a
   partir de salas-modelo, portas, minimapa, e **câmera por sala**
   (adiada da Fundação de propósito — só faz sentido quando existem salas
   menores que a tela).
3. ⬜ **Combate** — ataque à distância com mira independente do
   movimento, projéteis, vida/dano do jogador, framework de status
   effect (reaproveitável pelos bosses).
4. ⬜ **Limpeza (foco de dengue)** — pneu/vaso/poça/garrafa/calha como
   entidades interativas, ação `Interact` ligada pela primeira vez,
   condição de "sala limpa" que destrava as portas da Salas (depende de
   #2).
5. ⬜ **Inimigos & bosses** — IA base de inimigo + a escada de bosses por
   season (larva gigante → mosquito macho → fêmea → enxame → chefão
   final), lógica de "morte reinicia só a season atual" (depende de #3).
6. ⬜ **Poderes** — itens que modificam o ataque/stats do jogador,
   coletados durante a season (depende de #3).
7. ⬜ **Sprites** — primeira leva de placeholders em pixel art pra todo o
   elenco (personagem, inimigos, objetos, tileset), pro autor refinar
   depois. Pode rodar em paralelo com qualquer item acima assim que o
   elenco de cada fatia estiver definido.

**Dois itens com decisão de design ainda em aberto**, marcados como "fica
a cargo de quem constrói" nas seções acima — merecem uma rodada curta de
brainstorm antes de virarem plano, em vez de decisão unilateral de quem
implementar:
- **#6 Poderes**: lista exata de habilidades ainda não foi definida.
- **#5 Inimigos & bosses**: o padrão geral de combate já está definido
  (projétil legível + corpo segmentado), mas o ataque específico de cada
  boss individual ainda não.

Os itens #2 (Salas) e #3 (Combate) já têm spec suficiente pra ir direto
pra um plano de implementação, sem brainstorm extra.
