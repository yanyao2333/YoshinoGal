package models

import (
	"YoshinoGal/backend/library/database"
	"YoshinoGal/backend/library/playtime"
	"YoshinoGal/backend/library/watchdog"
)

type Library struct {
	Database       *database.SqliteGameLibrary
	GameLibraryDir string
	Watchdog       *watchdog.Interface
	Monitor        *playtime.Interface
}
