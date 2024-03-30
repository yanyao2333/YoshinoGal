package library

import (
	"YoshinoGal/backend/library/database"
	"YoshinoGal/backend/library/playtime"
	"YoshinoGal/backend/library/watchdog"
	"YoshinoGal/backend/logging"
	"YoshinoGal/backend/models"
	"github.com/pkg/errors"
	"os"
)

// InitGameLibrary 初始化游戏库，启动所有服务
func InitGameLibrary(gameDir string, scraperPriority []string) (*models.Library, error) {
	var log = logging.GetLogger()

	// 初始化游戏库
	log.Infof("初始化游戏库...")
	log.Infof("游戏目录: %s", gameDir)
	dbDir := gameDir + "/.YoshinoGal"
	_ = os.MkdirAll(dbDir, 0777)
	db, err := database.InitSQLiteDB(dbDir + "/library.db")
	if err != nil {
		return nil, errors.Wrap(err, "初始化数据库失败")
	}
	library := database.NewSqliteGameLibrary(db, gameDir)
	log.Infof("初始化数据库 %s 成功", dbDir+"/library.db")

	// 启动游戏库相关服务
	log.Infof("启动游戏库相关服务...")
	log.Infof("启动游戏时长监控服务...")
	monitor, err := playtime.StartMonitor(library)
	if err != nil {
		return nil, errors.Wrap(err, "启动监控服务失败")
	}
	log.Infof("启动游戏目录监控服务...")
	dog, err := watchdog.WatchGame(gameDir, library, scraperPriority)
	if err != nil {
		return nil, errors.Wrap(err, "启动监控服务失败")
	}
	log.Infof("启动游戏库相关服务成功")

	return &models.Library{
		Database:       library,
		GameLibraryDir: gameDir,
		Watchdog:       &dog,
		Monitor:        &monitor,
	}, nil
}
