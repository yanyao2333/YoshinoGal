package database

import (
	"YoshinoGal/internal/library/types"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"path/filepath"
)

type SqliteGameLibrary struct {
	db         *sql.DB
	LibraryDir string
}

func NewSqliteGameLibrary(db *sql.DB, libraryDir string) *SqliteGameLibrary {
	return &SqliteGameLibrary{db: db, LibraryDir: libraryDir}
}

func (s *SqliteGameLibrary) GetGamePathFromId(id int) (string, error) {
	var path string
	log.Debugf("通过ID获取游戏路径 数据库ID：%v", id)
	err := s.db.QueryRow("SELECT game_dir_path FROM galgames_metadata WHERE id = ?", id).Scan(&path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", types.GamePathNotFound
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
			return nil, types.GameScreenshotsNotFound
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

func (s *SqliteGameLibrary) GetGameNamePathMapping() (map[string]string, error) {
	log.Debugf("获取游戏映射...")
	rows, err := s.db.Query("SELECT name, game_dir_path FROM galgames_metadata")
	if err != nil {
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}
	defer rows.Close()

	gameIndex := make(map[string]string)
	for rows.Next() {
		var name string
		var gameDirPath string
		err = rows.Scan(&name, &gameDirPath)
		if err != nil {
			return nil, errors.Wrap(err, "读取数据库时发生错误")
		}
		gameIndex[name] = gameDirPath
	}
	log.Debugf("游戏映射共%d条", len(gameIndex))

	return gameIndex, nil
}

// GetGameIdPathMapping 获取游戏映射（格式为 map[id]路径 ）
func (s *SqliteGameLibrary) GetGameIdPathMapping() (map[int]string, error) {
	log.Debugf("获取游戏映射...")
	rows, err := s.db.Query("SELECT id, game_dir_path FROM galgames_metadata")
	if err != nil {
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}
	defer rows.Close()

	gameIndex := make(map[int]string)
	for rows.Next() {
		var id int
		var gameDirPath string
		err = rows.Scan(&id, &gameDirPath)
		if err != nil {
			return nil, errors.Wrap(err, "读取数据库时发生错误")
		}
		gameIndex[id] = gameDirPath
	}
	log.Debugf("游戏映射共%d条", len(gameIndex))

	return gameIndex, nil

}

func (s *SqliteGameLibrary) GetGameIdFromPath(path string) (int, error) {
	var id int
	log.Debugf("通过路径获取游戏ID 路径：%v", path)
	err := s.db.QueryRow("SELECT id FROM galgames_metadata WHERE game_dir_path = ?", filepath.Clean(path)).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, types.GameIdNotFound
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

func (s *SqliteGameLibrary) GetGamePlayTime(id int) (*types.GalgamePlayTime, error) {
	var playTime string
	log.Debugf("获取游戏时长 数据库ID：%v", id)
	err := s.db.QueryRow("SELECT play_time FROM galgames_metadata WHERE id = ?", id).Scan(&playTime)
	log.Debugf("获取游戏时长 数据库ID：%v 数据：%v", id, playTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "游戏不存在于数据库中")
		}
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}

	var playTimeData types.GalgamePlayTime
	err = json.Unmarshal([]byte(playTime), &playTimeData)
	if err != nil {
		return nil, errors.Wrap(err, "解析数据时发生错误")
	}

	return &playTimeData, nil
}

//func (s *SqliteGameLibrary) GetGameDataByPath(path string) (*types.GalgameMetadata, error) {
//	var name string
//	err := s.db.QueryRow("SELECT name FROM galgames_metadata WHERE game_dir_path = ?", path).Scan(&name)
//	if err != nil {
//		return nil, errors.Wrap(err, "查询数据库时发生错误")
//	}
//
//	return s.GetGameDataById(name)
//
//}

// InsertGamePlayTime 插入游戏时长数据
func (s *SqliteGameLibrary) InsertGamePlayTime(id int, playTime types.GalgamePlayTime) error {
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
func (s *SqliteGameLibrary) InsertGameLocalInfo(id int, localInfo types.GalgameLocalInfo) error {
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

func (s *SqliteGameLibrary) InsertGameMetadata(game *types.GalgameMetadata) error {
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

func (s *SqliteGameLibrary) GetGameDataById(id int) (*types.GalgameMetadata, error) {
	log.Debugf("通过ID获取游戏元数据 数据库ID：%v", id)
	var game types.GalgameMetadata
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

	json.Unmarshal([]byte(names), &game.Names)
	json.Unmarshal([]byte(scores), &game.Score)
	json.Unmarshal([]byte(sources), &game.MetadataSources)
	json.Unmarshal([]byte(developers), &game.Developers)
	json.Unmarshal([]byte(tags), &game.Tags)
	json.Unmarshal([]byte(playTime), &game.GalGameLocal.PlayTime)
	json.Unmarshal([]byte(localScreenshotsPaths), &game.GalGameLocal.LocalScreenshotsPaths)
	json.Unmarshal([]byte(screenshots), &game.RemoteScreenshotsUrls)

	return &game, nil
}

func InitSQLiteDB(dbPath string) (*sql.DB, error) {
	InitLogger()
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
