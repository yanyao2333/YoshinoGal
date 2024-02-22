package main

import "YoshinoGal/internal/routers"

//import "YoshinoGal/internal/game_play_time_monitor"
//
//func main() {
//	game_play_time_monitor.GamePlayTimeMonitor("E:\\GalGames", "E:\\GalGames\\playTime.json")
//	//	var a, _ = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
//	//		AppName:            "cool",
//	//		BaseDirectoryPath:  "./",
//	//		VersionAstilectron: "0.30.0",
//	//		VersionElectron:    "28.2.2",
//	//	})
//	//	defer a.Close()
//	//
//	//	// Start astilectron
//	//	a.Start()
//	//
//	//	var w, _ = a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
//	//		Center: astikit.BoolPtr(true),
//	//		Height: astikit.IntPtr(600),
//	//		Width:  astikit.IntPtr(600),
//	//	})
//	//	w.Create()
//	//}
//}

func main() {
	router := routers.SetupRouter()
	router.Run(":8080")
}
