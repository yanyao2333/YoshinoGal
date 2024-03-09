package main

import (
	"YoshinoGal/internal/router"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//import "YoshinoGal/internal/playtime"
//
//func main() {
//	playtime.GamePlayTimeMonitor("E:\\GalGames", "E:\\GalGames\\playTime.json")
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
	//gin.SetMode(gin.ReleaseMode)
	port := flag.Int("port", 8080, "a port to listen")
	flag.Parse()
	router := router.SetupRouter()
	err := router.Run(":" + strconv.Itoa(*port))
	portWasUsedErrString := "Only one usage of each socket address (protocol/network address/port) is normally permitted"
	if strings.Contains(err.Error(), portWasUsedErrString) {
		fmt.Println("端口被占用！返回码：114514")
		os.Exit(114514)
	}
	if err != nil {
		fmt.Println("启动失败！返回码：1919810")
		os.Exit(1919810)
	}
}
