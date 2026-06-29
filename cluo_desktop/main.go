package main

import (
	"context"
	"embed"
	"net/http"
	"net/http/httputil"
	"net/url"

	"cluo_desktop/updater"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func newAPIProxy() http.Handler {
	target, _ := url.Parse("https://api.clientvault.fr")
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
		},
	}
}

func main() {
	app := NewApp()
	upd := updater.NewUpdater()

	err := wails.Run(&options.App{
		Title:  "cluo_desktop",
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
