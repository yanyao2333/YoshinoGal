package library

import (
	"YoshinoGal/internal/library/database"
	"github.com/pkg/errors"
	"os"
)

// InitGameLibrary 初始化游戏库
func InitGameLibrary(gameDir string) (*database.SqliteGameLibrary, error) {
	InitLogger()
	log.Infof("初始化游戏库...")
	log.Infof("游戏目录: %s", gameDir)
	dbDir := gameDir + "/.YoshinoGal"
	os.MkdirAll(dbDir, 0777)
	db, err := database.InitSQLiteDB(dbDir + "/library.db")
	if err != nil {
		return nil, errors.Wrap(err, "初始化数据库失败")
	}
	library := database.NewSqliteGameLibrary(db, gameDir)
	log.Infof("初始化数据库 %s 成功", dbDir+"/library.db")
	return library, nil
}
