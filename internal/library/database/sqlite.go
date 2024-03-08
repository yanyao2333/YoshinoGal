package database

import (
	"YoshinoGal/internal/library/types"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type SqliteGameLibrary struct {
	db         *sql.DB
	LibraryDir string
}

func NewSqliteGameLibrary(db *sql.DB, libraryDir string) *SqliteGameLibrary {
	return &SqliteGameLibrary{db: db, LibraryDir: libraryDir}
}

// GetGameIndex 获取游戏索引（格式为 map[游戏名]路径 ）
func (s *SqliteGameLibrary) GetGameIndex() (map[string]string, error) {
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

	return gameIndex, nil

}

// GetGameNameByPath 通过游戏路径获取游戏名
func (s *SqliteGameLibrary) GetGameNameByPath(path string) (string, error) {
	var name string
	err := s.db.QueryRow("SELECT name FROM galgames_metadata WHERE game_dir_path = ?", path).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("游戏不存在于数据库中")
		}
		return "", errors.Wrap(err, "查询数据库时发生错误")
	}
	return name, nil
}

// IfHaveGame 检查游戏是否已存在于数据库中
func (s *SqliteGameLibrary) IfHaveGame(name, path string) (bool, error) {
	var count int
	if name != "" {
		err := s.db.QueryRow("SELECT COUNT(*) FROM galgames_metadata WHERE name = ?", name).Scan(&count)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, errors.New("游戏不存在于数据库中")
			}
			return false, errors.Wrap(err, "查询数据库时发生错误")
		}
		return count > 0, nil

	}
	if path != "" {
		err := s.db.QueryRow("SELECT COUNT(*) FROM galgames_metadata WHERE game_dir_path = ?", path).Scan(&count)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, errors.New("游戏不存在于数据库中")
			}
			return false, errors.Wrap(err, "查询数据库时发生错误")
		}
		return count > 0, nil
	}
	return false, errors.New("参数错误")
}

func (s *SqliteGameLibrary) GetGamePlayTime(name, path string) (*types.GalgamePlayTime, error) {
	//var playTime string
	//err := s.db.QueryRow("SELECT play_time FROM galgames_metadata WHERE name = ?", name).Scan(&playTime)
	//if err != nil {
	//	return nil, errors.Wrap(err, "查询数据库时发生错误")
	//}
	//
	//var playTimeData types.GalgamePlayTime
	//err = json.Unmarshal([]byte(playTime), &playTimeData)
	//if err != nil {
	//	return nil, errors.Wrap(err, "解析数据时发生错误")
	//}
	//
	//return &playTimeData, nil
	if name != "" {
		var playTime string
		err := s.db.QueryRow("SELECT play_time FROM galgames_metadata WHERE name = ?", name).Scan(&playTime)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.New("游戏不存在于数据库中")
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
	if path != "" {
		name, err := s.GetGameNameByPath(path)
		if err != nil {
			return nil, errors.WithMessage(err, "查询数据库时发生错误")
		}
		return s.GetGamePlayTime(name, "")
	}
	return nil, errors.New("参数错误")
}

func (s *SqliteGameLibrary) GetGameDataByPath(path string) (*types.GalgameMetadata, error) {
	var name string
	err := s.db.QueryRow("SELECT name FROM galgames_metadata WHERE game_dir_path = ?", path).Scan(&name)
	if err != nil {
		return nil, errors.Wrap(err, "查询数据库时发生错误")
	}

	return s.GetGameDataByName(name)

}

func (s *SqliteGameLibrary) InsertGamePlayTime(name, path string, playTime types.GalgamePlayTime) error {
	if name == "" {
		var err error
		name, err = s.GetGameNameByPath(path)
		if err != nil {
			return errors.WithMessage(err, "查询数据库时发生错误")
		}
	}
	// 插入数据
	stmt, err := s.db.Prepare(`
		UPDATE galgames_metadata
		SET play_time = ?
		WHERE name = ?
	`)
	if err != nil {
		return errors.Wrap(err, "准备插入数据时发生错误")
	}
	defer stmt.Close()

	var playTimeStr, _ = json.Marshal(playTime)
	_, err = stmt.Exec(
		playTimeStr,
		name,
	)
	if err != nil {
		return errors.Wrap(err, "插入数据时发生错误")
	}

	return nil

}

func (s *SqliteGameLibrary) InsertGameLocalInfo(name, path string, localInfo types.GalgameLocalInfo) error {
	if name == "" {
		var err error
		name, err = s.GetGameNameByPath(path)
		if err != nil {
			return errors.WithMessage(err, "查询数据库时发生错误")
		}
	}
	// 插入数据
	stmt, err := s.db.Prepare(`
		UPDATE galgames_metadata
		SET local_poster_path = ?,
			local_screenshots_paths = ?,
			game_dir_path = ?,
			play_time = ?
		WHERE name = ?
	`)
	if err != nil {
		return errors.Wrap(err, "准备插入数据时发生错误")
	}
	defer stmt.Close()

	var playTime, _ = json.Marshal(localInfo.PlayTime)
	var localScreenshotsPaths, _ = json.Marshal(localInfo.LocalScreenshotsPaths)
	_, err = stmt.Exec(
		localInfo.LocalPosterPath,
		localScreenshotsPaths,
		localInfo.GameDirPath,
		playTime,
		name,
	)
	if err != nil {
		return errors.Wrap(err, "插入数据时发生错误")
	}

	return nil
}

func (s *SqliteGameLibrary) InsertGameMetadata(game types.GalgameMetadata) error {
	var count int
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
	var localScreenshotsPaths, _ = json.Marshal(game.GalGameLocal.LocalScreenshotsPaths)
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
		game.GalGameLocal.LocalPosterPath,
		localScreenshotsPaths,
		game.GalGameLocal.GameDirPath,
		playTime,
	)
	if err != nil {
		return errors.Wrap(err, "插入数据时发生错误")
	}

	return nil
}

func (s *SqliteGameLibrary) GetGameDataByName(name string) (*types.GalgameMetadata, error) {
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
		FROM galgames_metadata WHERE name = ?
	`, name).Scan(
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
		return nil, errors.Wrap(err, "查询数据时发生错误")
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
