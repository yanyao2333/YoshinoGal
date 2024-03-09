package database

import "YoshinoGal/internal/library/types"

type GameLibraryDatabase interface {
	InsertGameMetadata(game *types.GalgameMetadata) error
	GetGameDataById(id int) (*types.GalgameMetadata, error)
	InsertGameLocalInfo(id int, localInfo types.GalgameLocalInfo) error
	InsertGamePlayTime(id int, playTime types.GalgamePlayTime) error
	GetGameIdPathMapping() (map[int]string, error)
	GetGamePlayTime(id int) (*types.GalgamePlayTime, error)
	GetGameIdFromPath(path string) (int, error)
	GetGameNamePathMapping() (map[string]string, error)
	GetGamePathFromId(id int) (string, error)
	GetGameScreenshots(id int) ([]string, error)
}
