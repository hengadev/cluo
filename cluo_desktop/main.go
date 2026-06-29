package main

import (
	"context"
	"embed"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"cluo_desktop/updater"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

// newAPIProxy returns a handler that proxies requests to the production API.
// Cookies are managed in Go's cookiejar rather than WebView2's cookie engine,
// which does not reliably persist Set-Cookie headers from intercepted responses.
func newAPIProxy() http.Handler {
	target, _ := url.Parse("https://api.clientvault.fr")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outURL := *target
		outURL.Path = r.URL.Path
		outURL.RawQuery = r.URL.RawQuery

		outReq, err := http.NewRequestWithContext(r.Context(), r.Method, outURL.String(), r.Body)
		if err != nil {
			http.Error(w, "proxy error", http.StatusBadGateway)
			return
		}
		for k, vv := range r.Header {
			if k == "Host" {
				continue
			}
			outReq.Header[k] = vv
		}

		resp, err := client.Do(outReq)
		if err != nil {
			http.Error(w, "proxy error", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// The jar has already captured Set-Cookie; omit them from the response
		// so WebView2 does not attempt to manage a parallel cookie store.
		for k, vv := range resp.Header {
			if k == "Set-Cookie" {
				continue
			}
			w.Header()[k] = vv
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body) //nolint:errcheck
	})
}

func main() {
	app := NewApp()
	upd := updater.NewUpdater()

	err := wails.Run(&options.App{
		Title:  "Cluo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: newAPIProxy(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			upd.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			upd,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
