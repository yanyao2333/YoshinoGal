package router

// LibraryScraperAll /library/scrape/all/run
type LibraryScraperAll struct {
	Directory string   `json:"directory" binding:"required"`
	Priority  []string `json:"priority" binding:"required"`
}

//// PlayTimeMonitorStart /playtime/monitor/start
//type PlayTimeMonitorStart struct {
//	GameBaseFolder string `json:"game_base_folder" binding:"required"`
//}

//// LibraryIndexRefresh /library/index/refresh
//type LibraryIndexRefresh struct {
//	Directory string `json:"directory" binding:"required"`
//}

//// GetLibraryIndex /library/index/get
//type GetLibraryIndex struct {
//	Directory string `json:"directory" binding:"required"`
//}
//
//// PosterWallIndex /library/index/posterwall
//type PosterWallIndex struct {
//	Directory string `json:"directory" binding:"required"`
//}

//// GetOneMetadata /library/metadata/get/one
//type GetOneMetadata struct {
//	GameName string `json:"game_name" binding:"required"`
//	//Directory string `json:"directory" binding:"required"`
//}

//// GetOneGamePlayTime /playtime/get/one
//type GetOneGamePlayTime struct {
//	Directory string `json:"directory" binding:"required"`
//}

// INeedGameName 只需要游戏名
type INeedGameId struct {
	GameId int `json:"game_id" binding:"required"`
}

// LoadLibrary /library/load
type LoadLibrary struct {
	Directory string `json:"directory" binding:"required"`
}
