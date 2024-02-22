package routers

import (
	"YoshinoGal/internal/game_play_time_monitor"
	"YoshinoGal/internal/scraper"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// 执行ScanGamesAndScrape 识别一个目录下的所有游戏并进行刮削
	router.POST("/game/scrape/all/run", func(c *gin.Context) {
		json := GameScraperAll{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
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

	router.GET("/game/scrape/all/status", func(c *gin.Context) {
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

	router.POST("/game/playtime/monitor/start", func(c *gin.Context) {
		json := GamePlayTimeMonitorStart{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
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
			game_play_time_monitor.StartMonitor(json.GameBaseFolder, json.GamePlayTimeFilePath)
		}()
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "咱收到请求了喵！已启动游戏时长监控器~",
		})
	})

	router.GET("/game/playtime/monitor/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "游戏时长监控器状态获取成功",
			"status":  game_play_time_monitor.MonitorRunningStatusFlag,
		})
	})

	router.POST("/game/playtime/monitor/stop", func(c *gin.Context) {
		game_play_time_monitor.MonitorStopFlag = true
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "已停止游戏时长监控器",
		})
	})

	return router
}
