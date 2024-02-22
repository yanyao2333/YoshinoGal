package routers

// GameScraperAll /game/scrape/all/run
type GameScraperAll struct {
	Directory string   `json:"directory"`
	Priority  []string `json:"priority"`
}

// GamePlayTimeMonitorStart /game/playtime/monitor/start
type GamePlayTimeMonitorStart struct {
	GameBaseFolder       string `json:"game_base_folder"`
	GamePlayTimeFilePath string `json:"game_play_time_file_path"`
}
