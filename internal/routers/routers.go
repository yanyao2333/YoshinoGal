package routers

import (
	"YoshinoGal/internal/game_play_time_monitor"
	"YoshinoGal/internal/scraper"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	SUCCESS = 0
	FAIL    = 1
)

func SetupRouter() *gin.Engine {
	InitLogger()
	router := gin.Default()

	// 基础路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ciallo～(∠・ω< )",
		})
	})

	router.GET("/img/*path", func(c *gin.Context) {
		imagePath := c.Param("path")
		fileName := path.Base(imagePath)
		if !strings.Contains(fileName, "screenshot") && !strings.Contains(fileName, "poster") {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "非法请求",
				"code":    FAIL,
			})
			return
		}
		log.Infof("请求图片: %s", imagePath)
		c.File(imagePath[1:])
	})

	// 执行ScanGamesAndScrape 识别一个目录下的所有游戏并进行刮削
	router.POST("/library/scrape/all/run", func(c *gin.Context) {
		json := LibraryScraperAll{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		err = os.MkdirAll(json.Directory+"/.YoshinoGal", 0777)
		if err != nil {
			log.Errorf("创建.YoshinoGal目录失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "创建" + json.Directory + "/.YoshinoGal目录失败",
			})
			return
		}
		if scraper.ScrapeAllStatus == 1 {
			c.JSON(http.StatusOK, gin.H{
				"code":    FAIL,
				"message": "咱已经在运行中了！别再重复请求了！",
			})
			return
		}
		go func() {
			err := scraper.ScanGamesAndScrape(json.Directory, json.Priority)
			if err != nil {
				log.Errorf("刮削失败: %s", err)
			}
		}()
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "咱收到刮削请求了喵！正在处理~",
		})
	})

	router.POST("/library/index/get", func(c *gin.Context) {
		json := GetLibraryIndex{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		GamesIndex, err := scraper.GetGamesIndex(json.Directory)
		if err != nil {
			log.Errorf("获取索引失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取索引失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取索引成功",
			"data":    GamesIndex,
		})
	})

	router.POST("/library/index/posterwall", func(c *gin.Context) {
		json := PosterWallIndex{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		posterwallIndex := map[string]string{}
		gamesIndex, err := scraper.GetGamesIndex(json.Directory)
		if err != nil {
			log.Errorf("获取索引失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取索引失败",
			})
			return
		}
		for gameName, gameDir := range gamesIndex {
			posterwallIndex[gameName] = gameDir + "/metadata/poster.jpg"
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取索引成功",
			"data":    posterwallIndex,
		})
	})

	router.POST("/playtime/get/one", func(c *gin.Context) {
		json := GetOneGamePlayTime{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		playTime := game_play_time_monitor.GetOneGamePlayTime(json.Directory)
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取游戏时长成功",
			"data":    playTime,
		})
	})

	router.POST("/library/metadata/get/one", func(c *gin.Context) {
		jsonData := GetOneMetadata{}
		err := c.BindJSON(&jsonData)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		GamesIndex, err := scraper.GetGamesIndex(jsonData.Directory)
		if err != nil {
			log.Errorf("获取索引失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取索引失败",
			})
			return
		}
		gameDir, exists := GamesIndex[jsonData.GameName]
		if !exists {
			c.JSON(http.StatusOK, gin.H{
				"code":    FAIL,
				"message": "游戏不存在",
			})
			return
		}
		metadata, err := scraper.GetMetadata(gameDir)
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取索引成功",
			"data":    metadata,
		})
	})

	router.POST("/library/index/refresh", func(c *gin.Context) {
		json := LibraryIndexRefresh{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		err = os.MkdirAll(json.Directory+"/.YoshinoGal", 0777)
		if err != nil {
			log.Errorf("创建.YoshinoGal目录失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "创建" + json.Directory + "/.YoshinoGal目录失败",
			})
			return
		}
		err = scraper.RefreshGamesIndex(json.Directory)
		if err != nil {
			log.Errorf("刷新失败: %s", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "咱收到刷新索引请求了喵！正在处理~",
		})
	})

	router.GET("/library/scrape/all/status", func(c *gin.Context) {
		if scraper.ScrapeAllStatus == 1 {
			c.JSON(http.StatusOK, gin.H{
				"code":             SUCCESS,
				"message":          "刮削正在进行中",
				"status":           scraper.ScrapeAllStatus,
				"each_game_status": scraper.GamesScrapeStatusMap,
			})
		} else if scraper.ScrapeAllStatus == 2 {
			c.JSON(http.StatusOK, gin.H{
				"code":             SUCCESS,
				"message":          "刮削失败",
				"status":           scraper.ScrapeAllStatus,
				"each_game_status": scraper.GamesScrapeStatusMap,
			})
		} else if scraper.ScrapeAllStatus == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":             SUCCESS,
				"message":          "刮削完成或还未进行刮削",
				"status":           scraper.ScrapeAllStatus,
				"each_game_status": scraper.GamesScrapeStatusMap,
			})
		}
	})

	router.POST("/playtime/monitor/start", func(c *gin.Context) {
		json := PlayTimeMonitorStart{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		err = os.MkdirAll(json.GameBaseFolder+"/.YoshinoGal", 0777)
		if err != nil {
			log.Errorf("创建.YoshinoGal目录失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "创建" + json.GameBaseFolder + "/.YoshinoGal目录失败",
			})
			return
		}
		if game_play_time_monitor.MonitorRunningStatusFlag == true {
			c.JSON(http.StatusOK, gin.H{
				"code":    FAIL,
				"message": "游戏时长监控器已经在运行中了！别再重复请求了！",
			})
			return
		}
		go func() {
			game_play_time_monitor.StartMonitor(json.GameBaseFolder, json.GameBaseFolder+"/.YoshinoGal/playTime.json")
		}()
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "咱收到请求了喵！已启动游戏时长监控器~",
		})
	})

	router.GET("/playtime/monitor/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "游戏时长监控器状态获取成功",
			"status":  game_play_time_monitor.MonitorRunningStatusFlag,
		})
	})

	router.POST("/playtime/monitor/stop", func(c *gin.Context) {
		game_play_time_monitor.MonitorStopFlag = true
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "已停止游戏时长监控器",
		})
	})

	return router
}
