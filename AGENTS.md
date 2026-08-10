# AGENTS.md — tiktok-go-sdk

**tiktok-go-sdk**: SDK Go (MIT) per le TikTok Developer API — OAuth2, Content Posting (video/photo), User Info, Video List, Research/Commercial. È una libreria: l'API pubblica e' il contratto.

Questo file definisce le regole operative per gli agent (Copilot, Claude, altri tool AI) quando lavorano in questo repository.

## Regole di lavoro (SEMPRE)

- **Ogni commit = release taggata `vX.Y.Z`**: nuova sezione in `CHANGELOG.md` (Keep a Changelog, in italiano) + `./scripts/release.sh` per il tag. Bump `minor` per novita' (nuovi endpoint/metodi), `patch` per fix. **Esenti**: auto-commit su `.claude/settings*.json` e commit di solo CI.
- **MAI `git push`**: lo fa sempre l'utente. MAI `Co-Authored-By` nei commit.
- **Gate prima di chiudere**: `go vet ./...` + `go test ./...` verdi.
- **Todo -> `BACKLOG.md`** (id stabili `TT-n`, milestone `## Mn — ...`): sincronizzati in issue/milestone GitHub da `cmd/backlog-sync`. Niente TODO sparsi nel codice.
- **Niente segreti**: `CLIENT_KEY`/`CLIENT_SECRET`/token solo da env, mai in codice, test, doc, output.
- **Mai rete reale nei test**: nuovi test con `httptest` + `WithBaseURL`. Le chiamate reali stanno dietro `//go:build integration`. Un test che chiama `open.tiktokapis.com` in `go test ./...` e' un bug.
- **Go 1.19 minimo** (matrice CI 1.19->1.22): niente `errors.Join`, `slices`/`maps`, `min`/`max`, `for i := range N`.
- **Lingua = inglese** per codice, commenti, godoc, errori. `CHANGELOG.md` e `BACKLOG.md` restano in italiano.

## Comandi

- `go build ./...` - `go vet ./...` - `go test -race ./...` - integrazione: `go test -tags integration ./test/`
- Release: `./scripts/release.sh` (tag dalla cima del CHANGELOG) - verifica: `./scripts/release.sh --check`
- Backlog: `go run ./cmd/backlog-sync -dry-run`

## Trappole note

- Import path: `github.com/HiWay-Media/tiktok-go-sdk/tiktok` (sottodirectory, non la radice del modulo).
- Un metodo nuovo va aggiunto a `ITiktok` in `tiktok/tiktok.go`, altrimenti e' irraggiungibile (vedi `ResearchAdQuery`, incompleto: fuori interfaccia e ignora l'errore).
- **TikTok risponde HTTP 200 anche sugli errori** (esito nel body): ogni metodo deve passare da `checkResponse(resp, "<operazione>")`, mai dal solo `resp.IsError()`.
- `debugPrint` stampa gli oggetti interi, **access token compreso**: `isDebug=true` solo in locale.
- Nomi esportati SCREAMING_SNAKE (`BASE_URL`, `PUBLIC_TO_EVERYONE`, ...) sono API pubblica: rinominare solo con alias deprecati.
- I metodi usano i **path relativi**, mai le `API_*` assolute: con un URL assoluto resty ignora la base URL e il test sfugge alla rete vera. Ogni endpoint nuovo va aggiunto a `TestEndpointsUseTheBaseURL`.
- `NewTikTok` valida le credenziali, applica un timeout di default (30s) e accetta opzioni variadiche; l'access token e' protetto da mutex.
- Dipendenze: solo `go-resty/resty/v2` e `golang.org/x/oauth2`. Non aggiungerne senza forte motivazione.
- Pre-1.0 (`v0.4.x`): breaking change ammessi con bump minor e nota nel CHANGELOG.

## Puntatori

- Backlog: `BACKLOG.md` - Changelog: `CHANGELOG.md` - Release: `scripts/release.sh` - Sync issue: `cmd/backlog-sync`
- CI: `.github/workflows/go-test.yml`, `go-build.yml`, `release.yml`, `backlog-sync.yml`
