# Datastar desktop shell

MVP: local first frame, remote server, Wails v2.

Draft · 7 October 2026

## 1. Goal

Ship a desktop window that paints a usable frame the moment it opens, then fills that frame from the existing Datastar server. The server is not packaged with the app.

Done means: double-click the binary, see the shell immediately, watch the remote bootstrap morph into `#app`, and have later clicks keep working against that server. If the server is down, the shell stays on screen and retries.

## 2. Decision

Wails v2. It is the stable line. v3 is still beta and is out of scope for the MVP. The shell is Go because that is the language the rest of the stack is written in. The webview is the OS webview: WebView2 on Windows, WKWebView on macOS, WebKitGTK on Linux.

The page origin is `http://wails.localhost`. That origin is stable and is what the remote server must allow. Do not load the remote app as the window URL. That puts the server back on the boot path.

## 3. Shape

Three pieces. Only the first is new.

| Piece | Where it lives | MVP job |
| --- | --- | --- |
| Shell | Wails binary | Window, embedded HTML, token read, inject header |
| First frame | `frontend/dist`, `go:embed` | Layout, stable ids, offline state, vendored `datastar.js` |
| App | Existing server | Bootstrap SSE, later actions, auth check |

The Go process does not proxy Datastar traffic. The page calls the remote host directly. Proxying would hide the origin and couple boot to the server.

## 4. Boot

1. Process starts, asset server serves the embedded page, window opens. First paint is local.
2. A bound Go method returns the API base URL and the bearer token. The page writes both into signals. The token is not baked into the HTML.
3. `data-on:load` fires `@get` to `{api}/desktop/bootstrap` with `Authorization` and `openWhenHidden`. Retry is bounded.
4. Server replies with `datastar-patch-elements` on `#app` and `datastar-patch-signals` setting `ready`. Morph matches on id.
5. Later actions use the same absolute base URL and the same header. Closing the window drops the stream. The server treats that as a normal disconnect.

### Shell page

Vendored `datastar.js`. No CDN. Signals the page can know without the server are set in markup: `ready` false, a connecting message, empty panels with stable ids. The bootstrap patch replaces contents of `#app`. It does not create the page.

## 5. Auth

Bearer token, not cookies. Datastar’s `@get` has a `headers` option and no credentials option, and the page origin is not the server origin, so session cookies will not be sent.

- Token lives in the OS keychain. Go reads it. The page never writes it to disk.
- Passed on every action as `Authorization: Bearer`. Also send `Datastar-Request`, which the client adds itself.
- 401 clears the ready signal and shows a signed-out frame. MVP can use a one-time paste into the shell to store the token. Login UI is later.
- Dev token may come from an env var. Production build ignores that path.

## 6. Server contract

One new route is enough for the MVP.

| Route | Response | Notes |
| --- | --- | --- |
| `GET /desktop/bootstrap` | `text/event-stream` | Patch `#app`, set `ready`, then close or hold |
| Existing action routes | SSE or HTML fragment | Called with absolute URLs from the shell |

CORS, reflected not starred:

- `Access-Control-Allow-Origin: http://wails.localhost`
- `Access-Control-Allow-Headers: Content-Type, Datastar-Request, Authorization`
- `Access-Control-Allow-Methods: GET, POST`
- No proxy buffering on the stream. `X-Accel-Buffering: no`.

Patch elements must carry the same ids the shell already rendered. A patch with no matching id is a bug, not a fallback.

## 7. Repo layout

- `main.go` — `wails.Run`, embed, bind the token method.
- `frontend/dist/index.html` — shell markup.
- `frontend/dist/datastar.js` — vendored bundle, pinned version.
- `frontend/dist/shell.css` — only what the first frame needs.
- `internal/token` — keychain get/set. Interface so tests can fake it.

No Vite, no Node app, no local API server. `wails build` is the package step.

## 8. Out of scope

- Wails v3, and any mobile target.
- Packaging the Datastar server, or proxying it through the asset handler.
- Auto-update, code signing, tray icon, multi-window.
- Cookie or SSO login inside the webview.
- Offline use beyond the first frame and a retry loop.

## 9. Acceptance

1. Cold start shows the shell before the bootstrap request returns. Airplane-mode launch still shows the shell.
2. With a valid token, `#app` morphs to the server frame and a follow-up action round-trips.
3. With no token or a 401, the shell shows a signed-out state and does not loop forever.
4. Server down: connecting state remains, retries stop at the cap, manual retry works.
5. Binary does not contain the server. Killing the window does not leave a local listener.

## 10. Risks

| Risk | Why it bites | MVP handling |
| --- | --- | --- |
| Webview drift | Linux WebKitGTK is not WebView2 | Smoke-test the shell page on each OS before styling it |
| CORS origin | A custom scheme sends a null or odd Origin | Log Origin on the server; allow only the real one |
| Token in the page | Signals ship to the server on every request | Keep the token out of signals; pass it only as a header |
| SSE idle death | Proxies buffer or cut idle streams | Heartbeat comment from the server, or short bootstrap that closes |

First build target is the machine you develop on. The other two OS webviews are a follow-up smoke test, not a blocker for the first morph.
