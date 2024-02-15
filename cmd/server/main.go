package main

func main() {
	//	var a, _ = astilectron.New(log.New(os.Stderr, "", 0), astilectron.Options{
	//		AppName:            "cool",
	//		BaseDirectoryPath:  "./",
	//		VersionAstilectron: "0.30.0",
	//		VersionElectron:    "28.2.2",
	//	})
	//	defer a.Close()
	//
	//	// Start astilectron
	//	a.Start()
	//
	//	var w, _ = a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
	//		Center: astikit.BoolPtr(true),
	//		Height: astikit.IntPtr(600),
	//		Width:  astikit.IntPtr(600),
	//	})
	//	w.Create()
	//}
}

//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//func main() {
//	router := gin.Default()
//
//	// 基础路由
//	router.GET("/", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{
//			"message": "欢迎使用Galgame信息刮削与展示API服务",
//		})
//	})
//
//	// 手动触发识别的API示例
//	router.POST("/trigger-recognition", func(c *gin.Context) {
//		// 这里将来会实现具体的识别逻辑
//		c.JSON(http.StatusOK, gin.H{
//			"message": "手动触发识别请求已接收，处理中...",
//		})
//	})
//
//	// 修正元数据的API示例
//	router.PATCH("/update-metadata", func(c *gin.Context) {
//		// 这里将来会实现具体的更新元数据逻辑
//		c.JSON(http.StatusOK, gin.H{
//			"message": "元数据更新请求已接收，处理中...",
//		})
//	})
//
//	// 提供游戏元数据的API示例
//	router.GET("/game-metadata", func(c *gin.Context) {
//		// 这里将来会实现具体的提供元数据逻辑
//		c.JSON(http.StatusOK, gin.H{
//			"message": "游戏元数据获取请求已接收，处理中...",
//		})
//	})
//
//	// 运行服务器
//	router.Run(":8080")
//}
