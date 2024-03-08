package database

import "YoshinoGal/internal/library/types"

type GameLibraryDatabase interface {
	InsertGameMetadata(game types.GalgameMetadata) error
	GetGameDataByName(name string) (*types.GalgameMetadata, error)
	GetGameDataByPath(path string) (*types.GalgameMetadata, error)
	InsertGameLocalInfo(name, path string, localInfo types.GalgameLocalInfo) error
	InsertGamePlayTime(name, path string, playTime types.GalgamePlayTime) error
	GetGameIndex() (map[string]string, error)
	GetGamePlayTime(name, path string) (*types.GalgamePlayTime, error)
	IfHaveGame(name, path string) (bool, error) // 是否存在游戏
}
