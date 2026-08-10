# CLAUDE.md — tiktok-go-sdk

**tiktok-go-sdk** (`github.com/HiWay-Media/tiktok-go-sdk`, MIT): SDK Go per le **TikTok Developer API** — OAuth2, Content Posting (video e photo, `PULL_FROM_URL`), User Info, Video List, Research/Commercial Content API. È una **libreria**, non un binario: l'API pubblica è il contratto, ogni rename esportato è un breaking change. Filosofia: wrapper sottile e leggibile sopra l'API TikTok (resty + oauth2), nessun framework, nessuna magia.

## Regole di lavoro (SEMPRE)

- **Ogni commit = release taggata `vX.Y.Z`**: nuova sezione in `CHANGELOG.md` (Keep a Changelog, in italiano) + tag creato con `./scripts/release.sh`. Bump `minor` per novità (nuovi endpoint, nuovi metodi dell'interfaccia), `patch` per fix. Senza chiederlo. **Esenti**: auto-commit su `.claude/settings*.json` e commit di solo CI/report.
- **MAI `git push`** — lo fa sempre l'utente (`git push && git push --tags`). MAI `Co-Authored-By` nei commit.
- **TDD**: il test si scrive **prima** del codice e lo si vede fallire *per il motivo giusto*. Per ogni fix di bug la **controprova è obbligatoria**: rimettere il difetto (patch temporanea), verificare che il test fallisca, ripristinare. Un test scritto dopo il fix e mai visto rosso non dimostra niente — potrebbe passare anche senza il fix. I test che devono passare in **entrambi** gli stati (guardie contro i falsi positivi, tipo «una risposta di successo non è un errore») vanno detti tali nel commento, altrimenti la controprova sembra riuscita a metà.
- **Gate prima di chiudere**: `go vet ./...` + `go test -race ./...` verdi (gli stessi check della CI, che gira su Go 1.19→1.22).
- **Todo → `BACKLOG.md`** (sorgente unica, id stabili `TT-n`). Non sparpagliare TODO nei commenti. Le milestone (`## Mn — …`) diventano milestone GitHub via `cmd/backlog-sync`: spuntare un item chiude la issue, spuntare l'ultimo item chiude la milestone.
- **Niente segreti** in codice, test, doc, esempi o output: `CLIENT_KEY`/`CLIENT_SECRET`/access token **solo da env**. Un token in un log o in una fixture è un bug di sicurezza.
- **Mai rete reale nei test**: i nuovi test girano contro `httptest` con `WithBaseURL` sul server finto; le chiamate reali stanno dietro `//go:build integration`. Un test che chiama `open.tiktokapis.com` in `go test ./...` è un test rotto (fallisce senza credenziali e brucia quota).
- **Compatibilità Go 1.19**: è il minimo dichiarato in `go.mod` e la prima entry della matrice CI. Niente `errors.Join` (1.20), `min`/`max` builtin o `slices`/`maps` (1.21), `for i := range N` (1.22).
- **Lingua = inglese**: codice, commenti, nomi, messaggi d'errore, godoc. **Eccezione: `CHANGELOG.md` e `BACKLOG.md` restano in italiano** (convenzione di progetto).

## Comandi

```bash
go build ./...
go vet ./...
go test ./...                      # tiktok/ + test/ + internal/, nessuna rete
go test -race ./...                # come la CI
go test -tags integration ./test/  # chiamate reali a TikTok (richiede le credenziali)
./scripts/release.sh               # gate + tag vX.Y.Z letto dalla cima del CHANGELOG
./scripts/release.sh --check       # solo verifica coerenza CHANGELOG/tag (usato dalla CI)
go run ./cmd/backlog-sync -dry-run # anteprima del sync BACKLOG.md → issue/milestone
```

## Architettura

- `tiktok/tiktok.go` — interfaccia pubblica `ITiktok` + costruttore `NewTikTok(clientKey, clientSecret, isDebug, opts ...Option)`; lo struct `tiktok` è privato, si espone solo l'interfaccia. L'access token è protetto da mutex: il client si condivide tra goroutine.
  **Ogni nuovo metodo va aggiunto a `ITiktok`**, altrimenti è irraggiungibile da fuori (è il bug di `ResearchAdQuery`).
- `tiktok/options.go` — opzioni funzionali (`WithBaseURL`, `WithTimeout`, `WithHTTPClient`, `WithDebug`, `WithAccessToken`) e `DefaultTimeout` (30s).
- `tiktok/resty.go` — trasporto: `restyPost`, `restyPostFormUrlEncoded`, `restyPostWithQueryParams`, `restyGet` (auth `Bearer` + header JSON) e `debugPrint`. Ogni endpoint nuovo passa da qui, non crea client suoi.
- `tiktok/constants.go` — path relativi `/v2/...` (quelli usati dai metodi) + costanti `API_*` con URL assoluto, tenute solo per compatibilità.
- `tiktok/oauth2.go` — `Endpoint` OAuth2 TikTok + `GetClientAccessTokenManagement` (grant `client_credentials`, per Research/Commercial). Il flusso **utente** (authorization_code + refresh) non c'è ancora: M2.
- `tiktok/content.go` — Content Posting: `CreatorInfo`, `PostVideo`/`PostVideoInit`, `PublishVideo` (status fetch), `PostPhoto`/`PostPhotoInit`, `GetVideoList`. Le forme `*Init` sono le firme posizionali storiche: delegano a quelle con lo struct, tenendo i default di allora.
- `tiktok/user.go` — `UserInfo` (`/v2/user/info/`).
- `tiktok/commercial.go` — Research/Commercial Content API (`ResearchAdQuery`, incompleto — vedi trappole).
- `tiktok/apierror.go` — `APIError` (errore tipizzato, `errors.As`) e `checkResponse`, l'unico punto che decide se una risposta TikTok è un fallimento.
- `tiktok/model*.go`, `request_*.go` — struct di request/response con tag JSON; `ErrorObject` (endpoint `/v2/...`) e `OAuthErrorObject` (token endpoint) sono le **due** buste d'errore.
- `tiktok/privacy.go`, `post_mode.go` — enum stringa + `CheckPrivacyLevel`/`CheckPostMode`, validazione **prima** della chiamata HTTP.
- `internal/backlog/`, `cmd/backlog-sync/` — tooling: parsing di `BACKLOG.md` e sync con issue/milestone GitHub. `internal/` = non importabile da fuori, non è API pubblica.
- `test/` — package `test` separato, consuma l'SDK come farebbe un utente: `newClient(t)` con credenziali placeholder per i test offline, `GetTikTok(t)` da env (con skip) per quelli dietro il tag `integration`.

## Trappole note / regole tecniche

- **L'import path è `github.com/HiWay-Media/tiktok-go-sdk/tiktok`**, non la radice del modulo: il package si chiama `tiktok` e sta in una sottodirectory. Ogni esempio che importa la radice non compila.
- **`debugPrint` stampa gli oggetti interi**, incluso l'`AccessTokenManagement` con il token in chiaro: `isDebug=true` è per lo sviluppo locale, mai in CI, mai in produzione, mai incollato in una issue.
- **I nomi esportati in SCREAMING_SNAKE** (`BASE_URL`, `API_USER_INFO`, `PUBLIC_TO_EVERYONE`, `DIRECT_POST`) non sono idiomatici ma **sono API pubblica**: rinominarli rompe chi importa. Eventuale rinomina = nuovi nomi + vecchi mantenuti come alias deprecati, mai una sostituzione secca.
- **I metodi usano i path relativi** (`QUERY_CREATOR_INFO`, `USER_INFO`, `OAUTH_TOKEN`, …), mai le `API_*` assolute: resty ignora la base URL se riceve un URL assoluto, quindi una `API_*` in un metodo nuovo **fa sfuggire i test alla rete reale**. Le `API_*` restano esportate solo per compatibilità. `TestEndpointsUseTheBaseURL` in `tiktok/tiktok_test.go` è il guardrail: aggiungi lì ogni endpoint nuovo.
- **TikTok risponde `HTTP 200` anche quando rifiuta la richiesta**, mettendo l'esito nel body: guardare solo `resp.IsError()` fa passare un fallimento per un successo. Ogni metodo **deve** chiamare `checkResponse(resp, "<operazione>")` prima di deserializzare — è il bug che ha motivato M6.
- **`ResearchAdQuery` è incompleto**: non è in `ITiktok`, ignora l'errore e non ritorna nulla. Non usarlo come modello per un endpoint nuovo (tracked in `BACKLOG.md`).
- **Le chiamate reali all'API stanno dietro `//go:build integration`** (`go test -tags integration ./test/`): `go test ./...` non tocca la rete. Le credenziali mancanti fanno `t.Skip`, non `t.Fatal`. Ogni test nuovo usa `httptest` + `WithBaseURL`.
- **Dipendenze: solo `go-resty/resty/v2` e `golang.org/x/oauth2`.** Aggiungerne altre solo con forte motivazione — un SDK con tante dipendenze è un peso per chi lo importa.
- Siamo **pre-1.0** (`v0.4.x`): breaking change ammessi con bump `minor` e nota esplicita nel CHANGELOG.
- Un endpoint TikTok nuovo = costante path relativa + `API_*` (per compatibilità), struct request/response, metodo che valida gli input, chiama l'helper resty **con il path relativo**, controlla `resp.IsError()`, fa `json.Unmarshal`, **e viene aggiunto a `ITiktok`** + test con `httptest`.

## Puntatori

- Backlog: `BACKLOG.md` · Changelog: `CHANGELOG.md` · Release: `scripts/release.sh` + `.github/workflows/release.yml`
- Sync issue/milestone: `cmd/backlog-sync` + `.github/workflows/backlog-sync.yml`
- CI: `.github/workflows/go-test.yml` (matrice 1.19→1.22), `go-build.yml`
- Doc TikTok: <https://developers.tiktok.com/doc/overview>
