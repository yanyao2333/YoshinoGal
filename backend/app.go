package backend

import (
	"context"
	"go.uber.org/zap"
)

type LocalConfig struct {
	GameDir         string   `json:"game_dir"`
	ScraperPriority []string `json:"scraper_priority"`
}

// YoshinoGalApplication 储存软件运行信息
type YoshinoGalApplication struct {
	Version     string             `json:"version"`
	LocalConfig *LocalConfig       `json:"local_config"`
	Logger      *zap.SugaredLogger `json:"logger"`
	CTX         context.Context    `json:"ctx"`
}

// NewApp creates a new App application struct
func NewApp() *YoshinoGalApplication {
	return &YoshinoGalApplication{}
}

// startup is called at application startup
func (a *YoshinoGalApplication) startup(ctx context.Context) {
	// Perform your setup here
	a.CTX = ctx
}

// domReady is called after front-end resources have been loaded
func (a *YoshinoGalApplication) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *YoshinoGalApplication) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *YoshinoGalApplication) shutdown(ctx context.Context) {
	// Perform your teardown here
}
