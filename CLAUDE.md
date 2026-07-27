# CLAUDE.md — tiktok-go-sdk

**tiktok-go-sdk** (`github.com/HiWay-Media/tiktok-go-sdk`, MIT): SDK Go per le **TikTok Developer API** — OAuth2, Content Posting (video e photo, `PULL_FROM_URL`), User Info, Video List, Research/Commercial Content API. È una **libreria**, non un binario: l'API pubblica è il contratto, ogni rename esportato è un breaking change. Filosofia: wrapper sottile e leggibile sopra l'API TikTok (resty + oauth2), nessun framework, nessuna magia.

## Regole di lavoro (SEMPRE)

- **Ogni commit = release taggata `vX.Y.Z`**: nuova sezione in `CHANGELOG.md` (Keep a Changelog, in italiano) + tag creato con `./scripts/release.sh`. Bump `minor` per novità (nuovi endpoint, nuovi metodi dell'interfaccia), `patch` per fix. Senza chiederlo. **Esenti**: auto-commit su `.claude/settings*.json` e commit di solo CI/report.
- **MAI `git push`** — lo fa sempre l'utente (`git push && git push --tags`). MAI `Co-Authored-By` nei commit.
- **Gate prima di chiudere**: `go vet ./...` + `go test ./...` verdi (gli stessi check della CI, che gira su Go 1.19→1.22).
- **Todo → `BACKLOG.md`** (sorgente unica, id stabili `TT-n`). Non sparpagliare TODO nei commenti. Le milestone (`## Mn — …`) diventano milestone GitHub via `cmd/backlog-sync`: spuntare un item chiude la issue, spuntare l'ultimo item chiude la milestone.
- **Niente segreti** in codice, test, doc, esempi o output: `CLIENT_KEY`/`CLIENT_SECRET`/access token **solo da env**. Un token in un log o in una fixture è un bug di sicurezza.
- **Mai rete reale nei test**: i nuovi test girano contro `httptest` con `SetBaseURL` sul server finto. Un test che chiama `open.tiktokapis.com` è un test rotto (fallisce in CI senza credenziali e brucia quota).
- **Compatibilità Go 1.19**: è il minimo dichiarato in `go.mod` e la prima entry della matrice CI. Niente `errors.Join` (1.20), `min`/`max` builtin o `slices`/`maps` (1.21), `for i := range N` (1.22).
- **Lingua = inglese**: codice, commenti, nomi, messaggi d'errore, godoc. **Eccezione: `CHANGELOG.md` e `BACKLOG.md` restano in italiano** (convenzione di progetto).

## Comandi

```bash
go build ./...
go vet ./...
go test ./...                      # test/ + internal/
./scripts/release.sh               # gate + tag vX.Y.Z letto dalla cima del CHANGELOG
./scripts/release.sh --check       # solo verifica coerenza CHANGELOG/tag (usato dalla CI)
go run ./cmd/backlog-sync -dry-run # anteprima del sync BACKLOG.md → issue/milestone
```

## Architettura

- `tiktok/tiktok.go` — interfaccia pubblica `ITiktok` + costruttore `NewTikTok(clientKey, clientSecret, isDebug)`; lo struct `tiktok` è privato, si espone solo l'interfaccia. **Ogni nuovo metodo va aggiunto a `ITiktok`**, altrimenti è irraggiungibile da fuori (è il bug di `ResearchAdQuery`).
- `tiktok/resty.go` — trasporto: `restyPost`, `restyPostFormUrlEncoded`, `restyPostWithQueryParams`, `restyGet` (auth `Bearer` + header JSON) e `debugPrint`. Ogni endpoint nuovo passa da qui, non crea client suoi.
- `tiktok/constants.go` — path `/v2/...` + costanti `API_*` con URL assoluto (`BASE_URL` + path).
- `tiktok/oauth2.go` — `Endpoint` OAuth2 TikTok + `GetClientAccessTokenManagement` (grant `client_credentials`, per Research/Commercial). Il flusso **utente** (authorization_code + refresh) non c'è ancora: M2.
- `tiktok/content.go` — Content Posting: `CreatorInfo`, `PostVideoInit`, `PublishVideo` (status fetch), `PostPhotoInit`, `GetVideoList`.
- `tiktok/user.go` — `UserInfo` (`/v2/user/info/`).
- `tiktok/commercial.go` — Research/Commercial Content API (`ResearchAdQuery`, incompleto — vedi trappole).
- `tiktok/model*.go`, `request_*.go` — struct di request/response con tag JSON; `ErrorObject` è la busta d'errore TikTok (oggi non tipizzata negli errori restituiti: M3).
- `tiktok/privacy.go`, `post_mode.go` — enum stringa + `CheckPrivacyLevel`/`CheckPostMode`, validazione **prima** della chiamata HTTP.
- `internal/backlog/`, `cmd/backlog-sync/` — tooling: parsing di `BACKLOG.md` e sync con issue/milestone GitHub. `internal/` = non importabile da fuori, non è API pubblica.
- `test/` — package `test` separato, consuma l'SDK come farebbe un utente (`GetTikTok()` da env).

## Trappole note / regole tecniche

- **L'import path è `github.com/HiWay-Media/tiktok-go-sdk/tiktok`**, non la radice del modulo: il package si chiama `tiktok` e sta in una sottodirectory. Ogni esempio che importa la radice non compila.
- **`debugPrint` stampa gli oggetti interi**, incluso l'`AccessTokenManagement` con il token in chiaro: `isDebug=true` è per lo sviluppo locale, mai in CI, mai in produzione, mai incollato in una issue.
- **I nomi esportati in SCREAMING_SNAKE** (`BASE_URL`, `API_USER_INFO`, `PUBLIC_TO_EVERYONE`, `DIRECT_POST`) non sono idiomatici ma **sono API pubblica**: rinominarli rompe chi importa. Eventuale rinomina = nuovi nomi + vecchi mantenuti come alias deprecati, mai una sostituzione secca.
- **`PostPhotoInit` usa la costante relativa `POST_PUBLISH_CONTENT_INIT`** mentre gli altri metodi usano le `API_*` assolute: funziona solo perché il client resty ha `SetBaseURL(BASE_URL)`. Uniformare è giusto, ma solo con un test che copra il path effettivo.
- **`ResearchAdQuery` è incompleto**: non è in `ITiktok`, ignora l'errore e non ritorna nulla. Non usarlo come modello per un endpoint nuovo (tracked in `BACKLOG.md`).
- **I test attuali chiamano l'API reale** (`TestGetClientAccessTokenManagement`) e dipendono da `CLIENT_KEY`/`CLIENT_SECRET` come secret di repo: sono legacy. Ogni test nuovo usa `httptest` e non deve poter fallire per una credenziale scaduta.
- **Dipendenze: solo `go-resty/resty/v2` e `golang.org/x/oauth2`.** Aggiungerne altre solo con forte motivazione — un SDK con tante dipendenze è un peso per chi lo importa.
- Siamo **pre-1.0** (`v0.2.x`): breaking change ammessi con bump `minor` e nota esplicita nel CHANGELOG.
- Un endpoint TikTok nuovo = costante path + `API_*`, struct request/response, metodo che valida gli input, chiama l'helper resty, controlla `resp.IsError()`, fa `json.Unmarshal`, **e viene aggiunto a `ITiktok`** + test con `httptest`.

## Puntatori

- Backlog: `BACKLOG.md` · Changelog: `CHANGELOG.md` · Release: `scripts/release.sh` + `.github/workflows/release.yml`
- Sync issue/milestone: `cmd/backlog-sync` + `.github/workflows/backlog-sync.yml`
- CI: `.github/workflows/go-test.yml` (matrice 1.19→1.22), `go-build.yml`
- Doc TikTok: <https://developers.tiktok.com/doc/overview>
