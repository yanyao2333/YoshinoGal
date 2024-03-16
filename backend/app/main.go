package app

import (
	"YoshinoGal/internal/library"
	"YoshinoGal/internal/library/database"
	"YoshinoGal/internal/library/playtime"
	"YoshinoGal/internal/library/scraper"
	"YoshinoGal/internal/library/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var (
	SUCCESS = 0
	FAIL    = 1
)

var gameLibrary *database.SqliteGameLibrary

func EnsureLibraryInitialized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/library/load" {
			c.Next()
			return
		}
		if gameLibrary == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "游戏库未初始化！请先通过 /library/load 初始化游戏库",
				"code":    FAIL,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	InitLogger()
	router := gin.Default()
	router.Use(gin.Recovery())

	libraryGroup := router.Group("/library").Use(EnsureLibraryInitialized())

	// 基础路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Ciallo～(∠・ω< )",
		})
	})

	libraryGroup.POST("/load", func(c *gin.Context) {
		json := LoadLibrary{}
		err := c.BindJSON(&json)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		gameLibrary, err = library.InitGameLibrary(json.Directory)
		if err != nil {
			log.Errorf("初始化游戏库失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "初始化游戏库失败",
				"code":    FAIL,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "初始化游戏库成功",
			"code":    SUCCESS,
		})
	})

	//app.GET("/img/*path", func(c *gin.Context) {
	//	imagePath := c.Param("path")
	//	fileName := path.Base(imagePath)
	//	if !strings.Contains(fileName, "screenshot") && !strings.Contains(fileName, "poster") {
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"message": "非法请求",
	//			"code":    FAIL,
	//		})
	//		return
	//	}
	//	log.Infof("请求图片: %s", imagePath)
	//	c.File(imagePath[1:])
	//})

	libraryGroup.GET("/game/screenshots", func(c *gin.Context) {
		gameId := c.Query("gid")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		gamePath, err := gameLibrary.GetGameScreenshots(gameIdInt)
		if err != nil {
			if errors.Is(err, types.GameScreenshotsNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    FAIL,
					"message": "找不到游戏截图",
				})
				return
			}
			log.Errorf("获取游戏路径失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取游戏路径失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取游戏截图成功",
			"data":    gamePath,
		})
		return
	})

	libraryGroup.GET("/game/poster", func(c *gin.Context) {
		gameId := c.Query("gid")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		gamePath, err := gameLibrary.GetGamePathFromId(gameIdInt)
		if err != nil {
			if errors.Is(err, types.GamePathNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    FAIL,
					"message": "游戏不存在",
				})
				return
			}
			log.Errorf("获取游戏路径失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取游戏路径失败",
			})
			return
		}
		c.File(gamePath + "/metadata/poster.jpg")
	})

	// 执行ScanGamesAndScrape 识别一个目录下的所有游戏并进行刮削
	libraryGroup.POST("/scrape/all/run", func(c *gin.Context) {
		if scraper.ScrapeAllStatus == 1 {
			c.JSON(http.StatusOK, gin.H{
				"code":    FAIL,
				"message": "咱已经在运行中了！别再重复请求了！",
			})
			return
		}
		go func() {
			pri := []string{"VNDB"}
			err := scraper.ScanGamesAndScrape(gameLibrary.LibraryDir, pri, gameLibrary)
			if err != nil {
				log.Errorf("刮削失败: %s", err)
			}
		}()
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "咱收到刮削请求了喵！正在处理~",
		})
	})

	libraryGroup.GET("/list", func(c *gin.Context) {
		GamesIndex, err := gameLibrary.GetGameIdPathMapping()
		if err != nil {
			log.Errorf("获取游戏列表失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取游戏列表失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取游戏列表成功",
			"data":    GamesIndex,
		})
	})

	//libraryGroup.POST("/index/posterwall", func(c *gin.Context) {
	//	posterwallIndex := map[string]string{}
	//	gamesIndex, err := gameLibrary.GetGameNamePathMapping()
	//	if err != nil {
	//		log.Errorf("获取索引失败: %s", err)
	//		c.JSON(http.StatusInternalServerError, gin.H{
	//			"code":    FAIL,
	//			"message": "获取索引失败",
	//		})
	//		return
	//	}
	//	for gameName, gameDir := range gamesIndex {
	//		posterwallIndex[gameName] = gameDir + "/metadata/poster.jpg"
	//	}
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    SUCCESS,
	//		"message": "获取索引成功",
	//		"data":    posterwallIndex,
	//	})
	//})

	// 获取单个游戏的总游戏时长
	libraryGroup.POST("/game/playtime/total", func(c *gin.Context) {
		gameId := c.Query("gid")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		playTime, err := gameLibrary.GetGamePlayTime(gameIdInt)
		if err != nil {
			log.Errorf("获取游戏时长失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取游戏时长失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取游戏时长成功",
			"data":    playTime.TotalTime,
		})
	})

	libraryGroup.GET("/game/metadata", func(c *gin.Context) {
		gameId := c.Query("gid")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			log.Errorf("请求格式错误: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "请求格式错误",
				"code":    FAIL,
			})
			return
		}
		metadata, err := gameLibrary.GetGameDataById(gameIdInt)
		if err != nil {
			log.Errorf("获取元数据失败: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    FAIL,
				"message": "获取元数据失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "获取元数据成功",
			"data":    metadata,
		})
	})

	// 后期会替换为watchdog
	//app.POST("/library/index/refresh", func(c *gin.Context) {
	//	json := LibraryIndexRefresh{}
	//	err := c.BindJSON(&json)
	//	if err != nil {
	//		log.Errorf("请求格式错误: %s", err)
	//		c.JSON(http.StatusBadRequest, gin.H{
	//			"message": "请求格式错误",
	//			"code":    FAIL,
	//		})
	//		return
	//	}
	//	err = os.MkdirAll(json.Directory+"/.YoshinoGal", 0777)
	//	if err != nil {
	//		log.Errorf("创建.YoshinoGal目录失败: %s", err)
	//		c.JSON(http.StatusInternalServerError, gin.H{
	//			"code":    FAIL,
	//			"message": "创建" + json.Directory + "/.YoshinoGal目录失败",
	//		})
	//		return
	//	}
	//	err = scraper.RefreshGamesIndex(json.Directory)
	//	if err != nil {
	//		log.Errorf("刷新失败: %s", err)
	//	}
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    SUCCESS,
	//		"message": "咱收到刷新索引请求了喵！正在处理~",
	//	})
	//})

	libraryGroup.GET("/scrape/all/status", func(c *gin.Context) {
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

	libraryGroup.POST("/playtime/monitor/start", func(c *gin.Context) {
		if playtime.MonitorRunningStatusFlag == true {
			c.JSON(http.StatusOK, gin.H{
				"code":    FAIL,
				"message": "游戏时长监控器已经在运行中了！别再重复请求了！",
			})
			return
		}
		go func() {
			playtime.StartMonitor(gameLibrary)
		}()
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "咱收到请求了喵！已启动游戏时长监控器~",
		})
	})

	libraryGroup.GET("/playtime/monitor/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "游戏时长监控器状态获取成功",
			"status":  playtime.MonitorRunningStatusFlag,
		})
	})

	libraryGroup.POST("/playtime/monitor/stop", func(c *gin.Context) {
		playtime.MonitorStopFlag = true
		c.JSON(http.StatusOK, gin.H{
			"code":    SUCCESS,
			"message": "已停止游戏时长监控器",
		})
	})

	return router
}
