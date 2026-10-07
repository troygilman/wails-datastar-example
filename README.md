# Desktop shell

A Wails v2 window that paints a local frame immediately, then fills `#app` from a remote Datastar server. The server is not packaged with the app, and this process does not proxy its traffic.

The design notes are in [datastar-desktop-shell-mvp.md](datastar-desktop-shell-mvp.md).

## Run

Install the Wails v2 CLI and the WebKitGTK 4.1 headers. This repo sets the `webkit2_41` build tag because Ubuntu 26.04 ships WebKitGTK 4.1.

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0
sudo apt install libgtk-3-dev libwebkit2gtk-4.1-dev
```

Package the binary with:

```bash
wails build
```

The binary is `build/bin/desktop-shell`. `wails build` is the package step. A plain `go build` does not open a window: Wails only links the desktop webview when the `production` tag is set, which `wails build` adds. The equivalent Go command is:

```bash
go build -tags production,webkit2_41 -o desktop-shell .
```

`wails dev` rebuilds on change and still reads `DESKTOP_DEV_TOKEN`. A production build ignores that variable.

## Server

The desktop process does not include this server. `make run` starts the server, then the shell, and stops the server when the shell exits. The token defaults to `dev` and the address to `127.0.0.1:8080`.

```bash
make run
make server   # server only
make dev      # shell only, against that server
```

`make run TOKEN=secret ADDR=127.0.0.1:9090` changes both sides together. The same commands without Make:

```bash
go run ./cmd/server
DESKTOP_DEV_TOKEN=dev wails dev
```

`DESKTOP_TOKEN` changes the bearer token the server accepts. `DESKTOP_ADDR` or `-addr` changes the listen address. Set `DESKTOP_API_BASE` on the shell when the server is not on `http://127.0.0.1:8080`, and set `DESKTOP_DEV_TOKEN` to the same value as `DESKTOP_TOKEN`.

`GET /desktop/bootstrap` patches `#app` and sets `ready`, then closes the stream. `POST /desktop/ping` is the follow-up action on the Ping button. Both require `Authorization: Bearer`. A missing or wrong token is 401. The response reflects `Origin` only for `http://wails.localhost`, `https://wails.localhost`, and `wails://wails`.

## Configure

| Variable | When it applies | Meaning |
| --- | --- | --- |
| `DESKTOP_API_BASE` | Shell | Absolute `http` or `https` URL of the Datastar server. Default `http://127.0.0.1:8080`. |
| `DESKTOP_DEV_TOKEN` | `wails dev`, and any shell build without the `production` tag | Bearer token used when the keychain has none. |
| `DESKTOP_TOKEN` | Server | Bearer token the server accepts. Default `dev`. |
| `DESKTOP_ADDR` | Server | Listen address. Default `:8080`. |

If the keychain already holds a token, that value wins. Paste a token into the shell to store a new one. It is written to the OS keychain under service `wails-datastar-example`, user `bearer`.

## What the window does

1. The embedded page in `frontend/dist` paints at once: a status line, and empty `#panel-main` and `#panel-side` inside `#app`.
2. `Session()` returns the API base and the bearer token. The page stores the base in the `api` signal. The token stays on `window.desktop` and is sent only as `Authorization: Bearer`.
3. With a token, the page calls `GET {api}/desktop/bootstrap` (`openWhenHidden`, at most five retries). The server patches `#app` and sets `ready`.
4. With no token, or a 401, the shell shows a paste field and stops. If the server cannot be reached, the connecting state remains until the retry cap, then Retry starts another bounded attempt.

Later buttons in the server frame should call the same host and header:

```html
<button data-on:click="@post($api + '/some/action', {headers: desktop.headers()})">
```

`frontend/dist/datastar.js` is the vendored Datastar 1.0.4 bundle. There is no CDN and no Node build.

## Server contract

`GET /desktop/bootstrap` returns `text/event-stream`. Patch the `#app` that is already on the page, set the `ready` signal, and close or hold the stream. A patch for an id the shell did not render is a bug.

Allow the page origin the webview actually sends. Wails v2 uses `wails://wails` on Linux and macOS, and `http://wails.localhost` on Windows. Reflect that origin. Also allow headers `Content-Type`, `Datastar-Request`, and `Authorization`, and methods `GET` and `POST`. Send `X-Accel-Buffering: no` on the stream.

## Layout

| Path | Role |
| --- | --- |
| `main.go` | `wails.Run` and the embedded assets |
| `app.go` | Bound `Session` and `SaveToken` |
| `cmd/server` | Remote Datastar process |
| `internal/server` | Bootstrap and ping handlers |
| `frontend/dist/` | First frame, vendored `datastar.js`, shell CSS and script |
| `internal/token` | Keychain get/set, with an in-memory store for tests |
| `internal/session` | API base and token resolution |

```bash
go test ./internal/...
go test -tags production ./internal/...
```
