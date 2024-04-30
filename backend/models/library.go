package models

// PosterWallGameShow 海报墙使用的数据映射
type PosterWallGameShow struct {
	GameId     int    `json:"game_id"`
	PosterPath string `json:"poster_path"`
	GameName   string `json:"game_name"`
	PosterB64  string `json:"poster_b64"`
}
