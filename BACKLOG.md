# Backlog — tiktok-go-sdk

Sorgente unica dei todo. Id stabili `TT-n`; spuntare, non cancellare.

> **Sync automatico issue & milestone**: questo file è la fonte di verità. Ogni `TT-n` diventa una issue GitHub (label `backlog`, milestone per sezione `## Mn — …`) via `cmd/backlog-sync` + workflow `.github/workflows/backlog-sync.yml`. Spuntare un item (`[x]`) chiude la issue al prossimo push su `main`; toglierlo la riapre. Quando **tutti** gli item di una milestone sono spuntati la milestone viene chiusa (e riaperta se ne riapri uno). Idempotente. Non aprire/chiudere issue o milestone a mano: edita qui.

Roadmap. **M1**: mettere in ordine tooling e rilascio (dove siamo ora). **M2**: OAuth utente completo — oggi c'è solo il grant `client_credentials`, quindi l'SDK non può agire *per conto di un creator* senza che il token gli venga passato da fuori. **M3**: rendere il client robusto (context, errori tipati, retry). **M4**: coprire il resto delle API. **M5**: qualità e test senza rete. **M6**: i difetti puntuali emersi dall'audit del codice (v0.3.1) — è la milestone da fare **per prima**, perché TT-25 significa che oggi una chiamata fallita può sembrare riuscita. Le versioni sono indicative: ogni item è comunque una release taggata a sé.

## M1 — Tooling & release (~v0.3)

- [x] **TT-1 — Regole di lavoro e CHANGELOG**: `CLAUDE.md` + `AGENTS.md` con le regole operative, `CHANGELOG.md` (Keep a Changelog, italiano) e la convenzione "ogni commit = release taggata `vX.Y.Z`". _(v0.3.0)_
- [x] **TT-2 — Sync BACKLOG → issue/milestone**: `internal/backlog` (parser) + `cmd/backlog-sync` (issue per `TT-n`, milestone per sezione, chiusura/riapertura idempotente, chiusura automatica della milestone completa) + workflow su push di `BACKLOG.md`. _(v0.3.0)_
- [x] **TT-3 — Tag automatizzato**: `scripts/release.sh` legge la versione dalla cima del `CHANGELOG.md`, esegue i gate (`go vet`, `go test`), verifica che il tag non esista e che la versione salga, crea `vX.Y.Z` annotato. `--check` per la CI. _(v0.3.0)_
- [x] **TT-4 — Release GitHub sui tag**: workflow `release.yml` che su `v*` verifica la coerenza tag↔CHANGELOG e pubblica la release con la sezione di changelog come note. _(v0.3.0)_
- [ ] **TT-5 — README corretto e utile**: l'esempio di import punta alla radice del modulo (`github.com/HiWay-Media/tiktok-go-sdk`) ma il package è in `/tiktok`: non compila. Correggere e aggiungere un esempio end-to-end reale (client → token → post).
- [ ] **TT-6 — Gate di qualità in CI**: aggiungere `go vet ./...` e `golangci-lint` al workflow di test (oggi la CI compila e testa ma non lint-a), più `govulncheck`.

## M2 — OAuth utente (~v0.4)

- [ ] **TT-7 — Scambio authorization code → access token**: `ExchangeCode(code string)` con grant `authorization_code` verso `/v2/oauth/token/`; oggi esiste solo `client_credentials`, quindi ogni operazione "per conto del creator" richiede un token ottenuto fuori dall'SDK.
- [ ] **TT-8 — Refresh token**: `RefreshAccessToken(refreshToken string)` + struct di risposta con `refresh_token`/`refresh_expires_in`. Senza, ogni integrazione deve rifare il login utente ogni 24h.
- [ ] **TT-9 — Scope e redirect URL configurabili**: oggi `RedirectURL` è `""` e gli scope sono commentati in `NewTikTok`, quindi `CodeAuthUrl()` produce un URL che TikTok rifiuta. Renderli parametri del costruttore (o opzioni funzionali, vedi TT-13).
- [ ] **TT-10 — `state` anti-CSRF**: `CodeAuthUrl()` hardcoda `"state-token"`. Accettare uno `state` dal chiamante e documentare che va verificato al callback.

## M3 — Robustezza client (~v0.5)

- [ ] **TT-11 — `context.Context` su tutti i metodi**: senza context non si può impostare un deadline né cancellare una chiamata; per una libreria è il difetto più visibile. Aggiungere varianti `...WithContext` (o cambiare le firme con bump minor documentato).
- [ ] **TT-12 — Errori tipati**: oggi ogni errore è `fmt.Errorf("... %s", resp.String())`, cioè il body grezzo. Introdurre `APIError` che deserializza `ErrorObject` (`code`, `message`, `log_id`) + `errors.As`, così il chiamante può distinguere token scaduto, rate limit e input non valido.
- [ ] **TT-13 — Opzioni funzionali nel costruttore**: `NewTikTok(key, secret, opts ...Option)` con `WithBaseURL` (indispensabile per i test con `httptest`), `WithHTTPClient`, `WithTimeout`, `WithDebug`. Mantenere la firma attuale come wrapper deprecato.
- [ ] **TT-14 — Retry/backoff**: ritentare `429` e `5xx` con backoff esponenziale e rispetto di `Retry-After` (resty lo supporta nativamente). Numero di tentativi configurabile, default conservativo.
- [ ] **TT-15 — `HealthCheck` sensato**: oggi fa `POST /` sull'host TikTok, che non è un endpoint di health: o lo si aggancia a una chiamata reale a basso costo o lo si deprecata.

## M4 — Copertura API (~v0.6)

- [ ] **TT-16 — `ResearchAdQuery` completo**: oggi non è in `ITiktok`, ignora l'errore e non ritorna niente. Firma con filtri tipizzati, response struct, gestione errori, aggiunta all'interfaccia.
- [ ] **TT-17 — Research Video Query**: la costante `API_RESEARCH_VIDEO_QUERY` esiste ma non c'è il metodo. Query con filtri + paginazione via cursor.
- [ ] **TT-18 — Upload `FILE_UPLOAD`**: oggi si supporta solo `PULL_FROM_URL`; aggiungere l'init con `source=FILE_UPLOAD` e l'upload del file (chunk + `Content-Range`), che è il flusso richiesto quando il video non è su un dominio verificato.
- [ ] **TT-19 — Campi configurabili per `UserInfo`**: `/v2/user/info/` richiede il query param `fields`, che oggi non viene mai passato (`restyGet(API_USER_INFO, nil)`): la chiamata torna errore. Passare i campi richiesti, con un default sensato.
- [ ] **TT-20 — Paginazione `GetVideoList`**: esporre `cursor`/`has_more` e i `fields` richiesti invece dei tre hardcoded, così si può scorrere oltre la prima pagina.

## M5 — Qualità & test (~v0.7)

- [ ] **TT-21 — Unit test con `httptest`**: un test per endpoint contro un server finto (richiesta attesa, risposta di esempio, caso d'errore TikTok). Richiede TT-13 per iniettare la base URL.
- [ ] **TT-22 — Test legacy senza rete**: i test in `test/` dipendono da `CLIENT_KEY`/`CLIENT_SECRET` e `TestGetClientAccessTokenManagement` chiama davvero TikTok: spostare le chiamate reali dietro un build tag `integration`, fuori da `go test ./...`.
- [ ] **TT-23 — Coverage in CI**: `go test -coverprofile` + riepilogo nel job summary, per vedere il progresso di TT-21.
- [ ] **TT-24 — Esempi e godoc**: cartella `examples/` compilabile (auth, post video, post photo, user info) e godoc sui tipi pubblici — per una libreria la doc è parte del prodotto.

## M6 — Audit: correttezza & igiene (~v0.4) — da fare per prima

Esito dell'audit del codice del 2026-07-28 (v0.3.0). Gli item sono ordinati per gravità: TT-25 e TT-26 sono difetti di correttezza che un utente dell'SDK non può aggirare dall'esterno.

- [ ] **TT-25 — L'errore TikTok in-band non viene rilevato**: TikTok risponde **HTTP 200 anche quando la chiamata fallisce**, mettendo l'esito in `error.code`; l'SDK guarda solo `resp.IsError()`, quindi deserializza una risposta vuota e ritorna `err == nil`. **Un post fallito è indistinguibile da uno riuscito.** Verificato in audit: `POST /v2/oauth/token/` con credenziali vuote risponde `HTTP 200` + `{"error":"invalid_request",…}` e `GetClientAccessTokenManagement` ritorna un `AccessTokenManagement` con token vuoto e nessun errore — poi ogni chiamata successiva parte con `Authorization: Bearer ` vuoto. Fix: dopo l'unmarshal, controllare la busta d'errore e ritornare errore se `code` non è `ok`/vuoto. È il prerequisito di TT-12 (tipizzazione).
- [ ] **TT-26 — Le buste d'errore sono due, non una**: gli endpoint `/v2/...` usano `ErrorObject` (`code`, `message`, `log_id`), l'endpoint OAuth usa una busta piatta e diversa (`error`, `error_description`, `log_id`) che **nessuno struct dell'SDK modella** — `AccessTokenManagement` non ha campi d'errore, quindi l'informazione viene buttata via anche volendo. Servono due tipi (o uno con entrambe le forme) e un unico punto che decide "questa risposta è un errore".
- [ ] **TT-27 — Igiene degli errori restituiti**: i messaggi sono copiati male e diagnosticano la chiamata sbagliata — `PostPhotoInit` e `GetVideoList` dicono entrambi `"post video init error"`, `UserInfo` dice `"creator info error"`, `PublishVideo` ha il typo `"publis video error"`. In più le sentinelle esportate non sono idiomatiche (`PrivacyLevelWrong`/`PhotoModeWrong` invece di `ErrPrivacyLevel`/`ErrPostMode`, messaggi con maiuscola e punto esclamativo, che staticcheck segnala con ST1005). Rinominare solo con alias deprecati: sono API pubblica.
- [ ] **TT-28 — Typo in campi esportati e tag JSON**: `DataPublishVideo.PubblishId` (due `b`) è API pubblica e va corretto con alias deprecato; `PublishStatusFetch.PublicalyAvailablePostId` ha lo stesso typo **anche nel tag JSON** (`publicaly_available_post_id`) — va verificato contro la doc TikTok, perché se il campo reale è `publicly_…` quel campo resta silenziosamente sempre vuoto.
- [ ] **TT-29 — Campi mancanti nelle response**: `DataQueryCreatorInfo` non espone `max_video_post_duration_sec`, che compare nell'esempio in cima al file ed è **il motivo per cui TikTok obbliga a chiamare creator_info** prima di postare (validare la durata del video); `UserInfo` modella solo `open_id`/`union_id`/`avatar_url` mentre l'endpoint ne ritorna molti altri (username, display_name, follower_count, …). Da allineare, insieme a TT-19.
- [ ] **TT-30 — Parametri decisi al posto del chiamante**: `PostPhotoInit` forza `disable_comment=false`, `auto_add_music=true` e `photo_cover_index=1` senza che si possano cambiare, e non valida che `photoUrls` non sia vuoto; `PostVideoInit` accetta un `description` che nella doc del video init **non esiste** (il campo documentato è solo `title`) e non valorizza mai `video_cover_timestamp_ms`, che pure è nello struct. Verificare contro la doc e portare i parametri nella firma (o nelle opzioni).
- [ ] **TT-31 — Nessun timeout di default**: il client resty è creato con `resty.New()` senza `SetTimeout`, quindi una chiamata a TikTok può restare appesa a tempo indefinito e bloccare la goroutine del chiamante. Per una libreria usata lato server è un difetto di robustezza: default sensato (es. 30s) + override via opzione (TT-13).
- [ ] **TT-32 — Data race sull'access token**: `SetAccessToken` scrive `o.accessToken` senza sincronizzazione mentre altre goroutine lo leggono in ogni `resty*`; il client è per sua natura condiviso (una chiamata per richiesta HTTP). Proteggere con un `sync.RWMutex` o documentare esplicitamente che il client non è sicuro da mutare dopo il primo uso — meglio la prima.
- [ ] **TT-33 — `NewTikTok` non valida e promette un errore che non arriva mai**: ritorna sempre `(client, nil)` e accetta `clientKey`/`clientSecret` vuoti, così l'errore si manifesta molto più tardi come un 200-con-errore da TikTok (TT-25). Validare gli input o togliere l'errore dalla firma — non entrambe le cose a metà.
- [x] **TT-34 — `gofmt` e gate di formato**: 14 file su 21 non erano formattati (indentazione mista tab/spazi in `constants.go`, `model_content.go`, `request_content.go`, `tiktok.go`, …). `gofmt -w` su tutto l'albero in un commit dedicato + job `fmt` in `go-test.yml`, pinnato a una singola versione di Go perché l'output di gofmt è cambiato tra le release. _(v0.3.2)_
- [x] **TT-35 — Igiene della CI**: `go mod tidy` sostituito da `go mod download` + `go mod verify` (in CI `tidy` riparava un `go.mod` incoerente invece di farlo fallire), secret `CLIENT_KEY`/`CLIENT_SECRET` spostati dal livello workflow al solo job che li usa, `go vet` aggiunto al gate, action aggiornate a checkout@v4/setup-go@v5. Rimossa da `dependabot.yml` la entry per `/tests` (directory inesistente) e da `go-build.yml` il `cache-dependency-path: subdir/go.sum`, path inesistente su cui setup-go fallisce. _(v0.3.3; lint e govulncheck restano in TT-6)_
- [ ] **TT-36 — Dipendenze ferme**: `go-resty/resty/v2` è alla v2.7.0 (2022), molte minor indietro. Valutare l'aggiornamento **verificando** che la versione scelta resti compatibile con Go 1.19 (è il minimo dichiarato e la prima entry della matrice CI): se non lo è, la scelta va fatta esplicitamente, non subita da un merge di Dependabot.
