package database

import "YoshinoGal/backend/models"

type GameLibraryDatabase interface {
	InsertGameMetadata(game *models.GalgameMetadata) error
	GetGameMetadata(id int) (*models.GalgameMetadata, error)
	InsertGameLocalInfo(id int, localInfo models.GalgameLocalInfo) error
	InsertGamePlayTime(id int, playTime models.GalgamePlayTime) error
	GetPosterWallMapping() (map[int]string, error)
	GetGamePlayTime(id int) (*models.GalgamePlayTime, error)
	GetGameIdFromPath(path string) (int, error)
	GetGamePath(id int) (string, error)
	GetGameScreenshots(id int) ([]string, error)
	RemoveGame(id int) error
}
