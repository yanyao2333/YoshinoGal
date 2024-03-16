package app

// LibraryScraperAll /library/scrape/all/run
type LibraryScraperAll struct {
	Directory string   `json:"directory" binding:"required"`
	Priority  []string `json:"priority" binding:"required"`
}

// INeedGameName 只需要游戏名
type INeedGameId struct {
	GameId int `json:"game_id" binding:"required"`
}

// LoadLibrary /library/load
type LoadLibrary struct {
	Directory string `json:"directory" binding:"required"`
}
