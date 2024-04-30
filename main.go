package main

import (
	"YoshinoGal/backend"
	"YoshinoGal/backend/app"
	"YoshinoGal/backend/logging"
	"context"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"os"
	"strings"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
	println("Requesting file:", requestedFilename)
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	// Create an instance of the app structure
	logger := logging.GetLogger()
	library := app.NewLibrary()
	yoshino := app.NewApp(library)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "YoshinoGal",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			yoshino.CTX = ctx
			yoshino.Logger = logging.GetLogger()
			yoshino.Version = backend.VERSION
			yoshino.SetLocalConfig()
			yoshino.InitLibrary()
			logger.Debugf("%v", yoshino.Library)
			//library = yoshino.Library
			wailsRuntime.EventsEmit(yoshino.CTX, "BackendReady")
		},
		OnShutdown: func(ctx context.Context) {

		},
		Bind: []interface{}{
			yoshino,
			library,
		},
	})

	if err != nil {
		logger.Errorf("Error:%v", err.Error())
	}
}
