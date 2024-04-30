package app

import (
	"YoshinoGal/backend"
	"YoshinoGal/backend/logging"
	"context"
	"encoding/json"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"os"
)

type LocalConfig struct {
	GameLibraryDir  string   `json:"game_library_dir"`
	ScraperPriority []string `json:"scraper_priority"`
	AppVersion      string   `json:"app_version"`
}

// YoshinoGalApplication 储存软件运行信息
type YoshinoGalApplication struct {
	Version     string             `json:"version"`
	LocalConfig *LocalConfig       `json:"local_config"`
	Logger      *zap.SugaredLogger `json:"logger"`
	CTX         context.Context    `json:"ctx"`
	Library     *Library           `json:"library"`
}

// NewApp creates a new App application struct
func NewApp(library *Library) *YoshinoGalApplication {
	return &YoshinoGalApplication{Library: library}
}

func (a *YoshinoGalApplication) SetLocalConfig() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		a.Logger.Warnf("配置文件不存在，一眼顶针，鉴定为第一次启动，加载引导页")
		wailsRuntime.EventsEmit(a.CTX, "FirstRunning")
		return
	}
	content, err := os.ReadFile("config.json")
	if err != nil {
		a.Logger.Errorf("读取配置文件失败")
		wailsRuntime.EventsEmit(a.CTX, "ConfigReadError", map[string]string{"errorMessage": err.Error()})
		return
	}
	a.LocalConfig = &LocalConfig{}
	err = json.Unmarshal(content, a.LocalConfig)
	if err != nil {
		a.Logger.Errorf("解析配置文件失败")
		wailsRuntime.EventsEmit(a.CTX, "ConfigReadError", map[string]string{"errorMessage": err.Error()})
		return
	}
}

// startup is called at application startup
func (a *YoshinoGalApplication) Startup(ctx context.Context) {
	// Perform your setup here
	a.CTX = ctx
	a.Logger = logging.GetLogger()
	a.Version = backend.VERSION
	a.SetLocalConfig()
	a.InitLibrary()
	wailsRuntime.EventsEmit(a.CTX, "BackendReady")
}

// domReady is called after front-end resources have been loaded
func (a *YoshinoGalApplication) domReady(ctx context.Context) {
	log.Debugf("dom ready")
}

// shutdown is called at application termination
func (a *YoshinoGalApplication) Shutdown(ctx context.Context) {
}

// InitLibrary 初始化游戏库 这个函数作为`InitGameLibraryFromConfig`事件的回调
func (a *YoshinoGalApplication) InitLibrary() {
	err := a.Library.InitGameLibrary(a.LocalConfig.GameLibraryDir, a.LocalConfig.ScraperPriority, a.CTX)
	if err != nil {
		wailsRuntime.EventsEmit(a.CTX, "LibraryInitError", map[string]string{"errorMessage": err.Error()})
	}
}
