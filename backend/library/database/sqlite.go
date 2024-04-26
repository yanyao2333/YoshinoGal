package database

import (
	"YoshinoGal/backend/logging"
	"YoshinoGal/backend/models"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"path/filepath"
)

var log = logging.GetLogger()

type SqliteGameLibrary struct {
	db         *sql.DB
	LibraryDir string
}

func NewSqliteGameLibrary(db *sql.DB, libraryDir string) *SqliteGameLibrary {
	return &SqliteGameLibrary{db: db, LibraryDir: filepath.Clean(libraryDir)}
}

func (s *SqliteGameLibrary) RemoveGame(id int) error {
	log.Debugf("删除游戏 数据库ID：%v", id)
	stmt, err := s.db.Prepare("DELETE FROM galgames_metadata WHERE id = ?")
	if err != nil {
		return errors.Wrap(err, "准备删除数据时发生错误")
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return errors.Wrap(err, "删除数据时发生错误")
	}

	return nil
}

func (s *SqliteGameLibrary) GetGamePath(id int) (string, error) {
	var path string
	log.Debugf("通过ID获取游戏路径 数据库ID：%v", id)
	err := s.db.QueryRow("SELECT game_dir_path FROM galgames_metadata WHERE id = ?", id).Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.GameNotFoundInDatabase
		}
		return "", errors.Wrap(err, "查询数据库时发生错误")
	}
	return path, nil
}

func (s *SqliteGameLibrary) GetGameScreenshots(id int) ([]string, error) {
	var screenshotsPaths string
	log.Debugf("获取游戏截图 数据库ID：%v", id)
	err := s.db.QueryRow("SELECT local_screenshots_paths FROM galgames_metadata WHERE id = ?", id).Scan(&screenshotsPaths)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.GameNotFoundInDatabase
		}
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}

	var screenshots []string
	err = json.Unmarshal([]byte(screenshotsPaths), &screenshots)
	if err != nil {
		return nil, errors.Wrap(err, "解析数据时发生错误")
	}

	return screenshots, nil

}

//func (s *SqliteGameLibrary) GetGameNamePathMapping() (map[string]string, error) {
//	log.Debugf("获取游戏映射...")
//	rows, err := s.db.Query("SELECT name, game_dir_path FROM galgames_metadata")
//	if err != nil {
//		return nil, errors.Wrap(err, "查询数据库时发生错误")
//	}
//	defer rows.Close()
//
//	gameIndex := make(map[string]string)
//	for rows.Next() {
//		var name string
//		var gameDirPath string
//		err = rows.Scan(&name, &gameDirPath)
//		if err != nil {
//			return nil, errors.Wrap(err, "读取数据库时发生错误")
//		}
//		gameIndex[name] = gameDirPath
//	}
//	log.Debugf("游戏映射共%d条", len(gameIndex))
//
//	return gameIndex, nil
//}

// GetPosterWallMapping 获取游戏海报墙映射（格式为 map[id]路径 ）
func (s *SqliteGameLibrary) GetPosterWallMapping() ([]models.PosterWallGameShow, error) {
	log.Debugf("获取海报墙所使用的游戏数据映射...")
	log.Debugf("%v", s.db)
	rows, err := s.db.Query("SELECT id, game_dir_path, name FROM galgames_metadata")
	if err != nil {
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}
	defer rows.Close()

	gameList := make([]models.PosterWallGameShow, 10)
	for rows.Next() {
		var id int
		var gameDirPath string
		var name string
		err = rows.Scan(&id, &gameDirPath, &name)
		if err != nil {
			return nil, errors.Wrap(err, "读取数据库时发生错误")
		}
		gameList = append(gameList, models.PosterWallGameShow{
			GameId:     id,
			PosterPath: gameDirPath + "/metadata/poster.jpg",
			GameName:   name,
		})
	}
	log.Debugf("游戏映射共%d条", len(gameList))

	return gameList, nil

}

func (s *SqliteGameLibrary) GetGameIdFromPath(path string) (int, error) {
	var id int
	log.Debugf("通过路径获取游戏ID 路径：%v", path)
	err := s.db.QueryRow("SELECT id FROM galgames_metadata WHERE game_dir_path = ?", filepath.Clean(path)).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.CannotMatchGameIDFromPathInDatabase
		}
		return 0, errors.Wrap(err, "查询数据库时发生错误")
	}
	return id, nil
}

// GetGameNameByPath 通过游戏路径获取游戏名
//func (s *SqliteGameLibrary) GetGameNameByPath(path string) (string, error) {
//	var name string
//	err := s.db.QueryRow("SELECT name FROM galgames_metadata WHERE game_dir_path = ?", path).Scan(&name)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return "", errors.New("游戏不存在于数据库中")
//		}
//		return "", errors.Wrap(err, "查询数据库时发生错误")
//	}
//	return name, nil
//}

// IfHaveGame 检查游戏是否已存在于数据库中
//func (s *SqliteGameLibrary) IfHaveGame(name, path string) (bool, error) {
//	var count int
//	if name != "" {
//		err := s.db.QueryRow("SELECT COUNT(*) FROM galgames_metadata WHERE name = ?", name).Scan(&count)
//		if err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				return false, errors.New("游戏不存在于数据库中")
//			}
//			return false, errors.Wrap(err, "查询数据库时发生错误")
//		}
//		return count > 0, nil
//
//	}
//	if path != "" {
//		err := s.db.QueryRow("SELECT COUNT(*) FROM galgames_metadata WHERE game_dir_path = ?", path).Scan(&count)
//		if err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				return false, errors.New("游戏不存在于数据库中")
//			}
//			return false, errors.Wrap(err, "查询数据库时发生错误")
//		}
//		return count > 0, nil
//	}
//	return false, errors.New("参数错误")
//}

// GetAllGamePlayTime 获取数据库中所有游戏的游戏时长（用于在monitor初始化时保证数据一致）
func (s *SqliteGameLibrary) GetAllGamePlayTime() (map[int]models.GalgamePlayTime, error) {
	rows, err := s.db.Query("SELECT id, play_time FROM galgames_metadata")
	if err != nil {
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}
	defer rows.Close()

	playTimeIndex := make(map[int]models.GalgamePlayTime)
	for rows.Next() {
		var id int
		var playTime string
		err = rows.Scan(&id, &playTime)
		if err != nil {
			return nil, errors.Wrap(err, "读取数据库时发生错误")
		}
		var playTimeData models.GalgamePlayTime
		err = json.Unmarshal([]byte(playTime), &playTimeData)
		if err != nil {
			return nil, errors.Wrap(err, "解析数据时发生错误")
		}
		playTimeIndex[id] = playTimeData
	}

	return playTimeIndex, nil
}

func (s *SqliteGameLibrary) GetGamePlayTime(id int) (*models.GalgamePlayTime, error) {
	var playTime string
	log.Debugf("获取游戏时长 数据库ID：%v", id)
	err := s.db.QueryRow("SELECT play_time FROM galgames_metadata WHERE id = ?", id).Scan(&playTime)
	log.Debugf("获取游戏时长 数据库ID：%v 数据：%v", id, playTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.GameNotFoundInDatabase
		}
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}

	var playTimeData models.GalgamePlayTime
	err = json.Unmarshal([]byte(playTime), &playTimeData)
	if err != nil {
		log.Errorf("解析数据时发生错误：%s 使用空数据代替", err)
		return &models.GalgamePlayTime{EachDayTime: make(map[string]int64)}, nil
	}
	if playTimeData.EachDayTime == nil {
		playTimeData.EachDayTime = make(map[string]int64)
	}

	return &playTimeData, nil
}

//func (s *SqliteGameLibrary) GetGameDataByPath(path string) (*models.GalgameMetadata, error) {
//	var name string
//	err := s.db.QueryRow("SELECT name FROM galgames_metadata WHERE game_dir_path = ?", path).Scan(&name)
//	if err != nil {
//		return nil, errors.Wrap(err, "查询数据库时发生错误")
//	}
//
//	return s.GetGameMetadata(name)
//
//}

// InsertGamePlayTime 插入游戏时长数据
func (s *SqliteGameLibrary) InsertGamePlayTime(id int, playTime models.GalgamePlayTime) error {
	log.Debugf("插入游戏时长 数据库ID：%v", id)
	stmt, err := s.db.Prepare(`
		UPDATE galgames_metadata
		SET play_time = ?
		WHERE id = ?
	`)
	if err != nil {
		return errors.Wrap(err, "准备插入数据时发生错误")
	}
	defer stmt.Close()

	var playTimeStr, _ = json.Marshal(playTime)
	_, err = stmt.Exec(
		playTimeStr,
		id,
	)
	if err != nil {
		return errors.Wrap(err, "插入数据时发生错误")
	}

	return nil

}

// InsertGameLocalInfo 插入游戏本地信息
func (s *SqliteGameLibrary) InsertGameLocalInfo(id int, localInfo models.GalgameLocalInfo) error {
	log.Debugf("插入游戏本地信息 数据库ID：%v", id)
	stmt, err := s.db.Prepare(`
		UPDATE galgames_metadata
		SET local_poster_path = ?,
			local_screenshots_paths = ?,
			game_dir_path = ?,
			play_time = ?
		WHERE id = ?
	`)
	if err != nil {
		return errors.Wrap(err, "准备插入数据时发生错误")
	}
	defer stmt.Close()

	var playTime, _ = json.Marshal(localInfo.PlayTime)
	var newScreenshotsPaths []string
	for _, url := range localInfo.LocalScreenshotsPaths {
		newScreenshotsPaths = append(newScreenshotsPaths, filepath.Clean(url))
	}
	var localScreenshotsPaths, _ = json.Marshal(newScreenshotsPaths)
	_, err = stmt.Exec(
		filepath.Clean(localInfo.LocalPosterPath),
		localScreenshotsPaths,
		filepath.Clean(localInfo.GameDirPath),
		playTime,
		id,
	)
	if err != nil {
		return errors.Wrap(err, "插入数据时发生错误")
	}

	return nil
}

func (s *SqliteGameLibrary) InsertGameMetadata(game *models.GalgameMetadata) error {
	var count int
	log.Debugf("插入游戏元数据 游戏名：%v", game.Name)
	err := s.db.QueryRow("SELECT COUNT(*) FROM galgames_metadata WHERE name = ?", game.Name).Scan(&count)
	if err != nil {
		return errors.Wrap(err, "查询数据库时发生错误")
	}
	if count > 0 {
		log.Warnf("游戏 %s 已存在于数据库中，覆盖掉他", game.Name)
	}

	// 插入数据
	stmt, err := s.db.Prepare(`
        INSERT INTO galgames_metadata(
            name,
            names,
            release_date,
            score,
            remote_poster_url,
            description,
            metadata_sources,
            length_minutes,
            length,
            dev_status,
            remote_screenshots_urls,
            developers,
            is_r18,
            tags,
            local_poster_path,
            local_screenshots_paths,
			game_dir_path,
            play_time
        ) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return errors.Wrap(err, "准备插入数据时发生错误")
	}
	defer stmt.Close()

	var names, _ = json.Marshal(game.Names)
	var scores, _ = json.Marshal(game.Score)
	var sources, _ = json.Marshal(game.MetadataSources)
	var developers, _ = json.Marshal(game.Developers)
	var tags, _ = json.Marshal(game.Tags)
	var screenshots, _ = json.Marshal(game.RemoteScreenshotsUrls)
	var playTime, _ = json.Marshal(game.GalGameLocal.PlayTime)
	var newScreenshotsPaths []string
	for _, url := range game.GalGameLocal.LocalScreenshotsPaths {
		newScreenshotsPaths = append(newScreenshotsPaths, filepath.Clean(url))
	}
	var localScreenshotsPaths, _ = json.Marshal(newScreenshotsPaths)
	_, err = stmt.Exec(
		game.Name,
		names,
		game.ReleaseDate,
		scores,
		game.RemotePosterUrl,
		game.Description,
		sources,
		game.LengthMinutes,
		game.Length,
		game.DevStatus,
		screenshots,
		developers,
		game.IsR18,
		tags,
		filepath.Clean(game.GalGameLocal.LocalPosterPath),
		localScreenshotsPaths,
		filepath.Clean(game.GalGameLocal.GameDirPath),
		playTime,
	)
	if err != nil {
		return errors.Wrap(err, "插入数据时发生错误")
	}
	log.Infof("插入游戏元数据成功 游戏名：%v", game.Name)

	return nil
}

func (s *SqliteGameLibrary) GetGameMetadata(id int) (*models.GalgameMetadata, error) {
	log.Debugf("通过ID获取游戏元数据 数据库ID：%v", id)
	var game models.GalgameMetadata
	var names string
	var scores string
	var sources string
	var developers string
	var tags string
	var screenshots string
	var playTime string
	var localScreenshotsPaths string
	err := s.db.QueryRow(`
		SELECT
			name,
			names,
			release_date,
			score,
			remote_poster_url,
			description,
			metadata_sources,
			length_minutes,
			length,
			dev_status,
			remote_screenshots_urls,
			developers,
			is_r18,
			tags,
			local_poster_path,
			local_screenshots_paths,
			game_dir_path,
			play_time
		FROM galgames_metadata WHERE id = ?
	`, id).Scan(
		&game.Name,
		&names,
		&game.ReleaseDate,
		&scores,
		&game.RemotePosterUrl,
		&game.Description,
		&sources,
		&game.LengthMinutes,
		&game.Length,
		&game.DevStatus,
		&screenshots,
		&developers,
		&game.IsR18,
		&tags,
		&game.GalGameLocal.LocalPosterPath,
		&localScreenshotsPaths,
		&game.GalGameLocal.GameDirPath,
		&playTime,
	)
	if err != nil {
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}

	_ = json.Unmarshal([]byte(names), &game.Names)
	_ = json.Unmarshal([]byte(scores), &game.Score)
	_ = json.Unmarshal([]byte(sources), &game.MetadataSources)
	_ = json.Unmarshal([]byte(developers), &game.Developers)
	_ = json.Unmarshal([]byte(tags), &game.Tags)
	_ = json.Unmarshal([]byte(playTime), &game.GalGameLocal.PlayTime)
	_ = json.Unmarshal([]byte(localScreenshotsPaths), &game.GalGameLocal.LocalScreenshotsPaths)
	_ = json.Unmarshal([]byte(screenshots), &game.RemoteScreenshotsUrls)

	return &game, nil
}

func InitSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "打开数据库时发生错误")
	}

	// 创建表格
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS galgames_metadata (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            release_date TEXT,
            names TEXT,
            score TEXT,
            description TEXT,
            remote_poster_url TEXT,
            metadata_sources TEXT,
            length_minutes INTEGER,
            length INTEGER,
            dev_status INTEGER,
            remote_screenshots_urls TEXT,
            developers TEXT,
            is_r18 INTEGER,
            tags TEXT,
            local_poster_path TEXT,
            local_screenshots_paths LIST,
            game_dir_path TEXT,
            play_time TEXT
        )
    `)
	if err != nil {
		return nil, errors.Wrap(err, "创建表格时发生错误")
	}

	log.Infof("初始化数据库成功，数据库路径：%s", dbPath)

	return db, nil
}
