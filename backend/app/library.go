package app

import (
	"YoshinoGal/backend/library/database"
	"YoshinoGal/backend/library/playtime"
	"YoshinoGal/backend/library/scraper"
	"YoshinoGal/backend/library/watchdog"
	"YoshinoGal/backend/logging"
	"YoshinoGal/backend/models"
	"context"
	"github.com/pkg/errors"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

type Library struct {
	Database       *database.SqliteGameLibrary
	GameLibraryDir string
	ScrapePriority []string
	Watchdog       *watchdog.Interface
	Monitor        *playtime.Interface
	CTX            context.Context
}

func NewLibrary() *Library {
	return &Library{}
}

var log = logging.GetLogger()

// GameLibraryInterface 游戏库接口 所有操作都被封装在这里
type GameLibraryInterface interface {
	GetPosterWall() (map[int]string, error)
	InitGameLibrary(gameDir string) error
	ManualScrapeLibrary()
}

func (a *Library) ManualScrapeLibrary() {
	go func() {
		err := scraper.ScanGamesAndScrape(a.GameLibraryDir, a.ScrapePriority, a.Database)
		if err != nil {
			wailsRuntime.EventsEmit(a.CTX, "GlobalRuntimeError", map[string]string{"errorMessage": err.Error(), "errorName": "ManualScrapeError"})
			log.Errorf("启动刮削后失败！%v", err)
		}
	}()
	log.Infof("手动启动刮削成功！")
}

// GetPosterWall 获取游戏海报墙
func (a *Library) GetPosterWall() ([]models.PosterWallGameShow, error) {
	mapping, err := a.Database.GetPosterWallMapping()
	if err != nil {
		return nil, errors.WithMessage(err, "获取游戏海报失败")
	}
	return mapping, nil
}

// InitGameLibrary 初始化游戏库，启动所有服务
func (a *Library) InitGameLibrary(gameDir string, scraperPriority []string, ctx context.Context) error {

	// 初始化游戏库
	log.Infof("初始化游戏库...")
	log.Infof("游戏目录: %s", gameDir)
	dbDir := gameDir + "/.YoshinoGal"
	_ = os.MkdirAll(dbDir, 0777)
	db, err := database.InitSQLiteDB(dbDir + "/library.db")
	if err != nil {
		return errors.Wrap(err, "初始化数据库失败")
	}
	library := database.NewSqliteGameLibrary(db, gameDir)
	log.Infof("初始化数据库 %s 成功", dbDir+"/library.db")

	// 启动游戏库相关服务
	log.Infof("启动游戏库相关服务...")
	log.Infof("启动游戏时长监控服务...")
	monitor, err := playtime.StartMonitor(library)
	if err != nil {
		return errors.Wrap(err, "启动监控服务失败")
	}
	log.Infof("启动游戏目录监控服务...")
	//dog, err := watchdog.WatchGame(gameDir, library, scraperPriority)
	if err != nil {
		return errors.WithMessage(err, "启动监控服务失败")
	}
	log.Infof("启动游戏库相关服务成功")

	a.Database = library
	//a.Watchdog = &dog
	a.Monitor = &monitor
	a.GameLibraryDir = gameDir
	a.ScrapePriority = scraperPriority
	a.CTX = ctx

	return nil
}
