package playtime

import (
	"YoshinoGal/internal/library/database"
	"YoshinoGal/internal/library/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var GamePlayManager = NewGamePlayTimeManager()

var (
	user32                        = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow       = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId  = user32.NewProc("GetWindowThreadProcessId")
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
)

var MonitorRunningStatusFlag = false // 游戏时长监控器运行状态
var MonitorStopFlag = false          // 游戏时长监控器停止标志

func StartMonitor(libraryDB *database.SqliteGameLibrary) {
	InitLogger()
	defer log.Sync()
	if MonitorRunningStatusFlag {
		log.Warnf("游戏时长监控器已经在运行中了！")
		return
	}
	MonitorRunningStatusFlag = true
	go GamePlayTimeMonitor(libraryDB)
}

type gameFolders []string

type gamePlayTimeInfo struct {
	FolderPath     string `json:"folder_path"`
	GameName       string `json:"game_name"`
	LatestPlayTime int64  `json:"latest_play_time"`
	Date           string `json:"date"`
	TotalPlayTime  int64  `json:"total_play_time"`
}

type gameFolderPlayTime struct {
	DailyPlayTime  map[string]int64 `json:"daily_play_timePlayTimes"` // 日期 -> 游戏时间
	TotalPlayTime  int64            `json:"total_play_time"`          // 文件夹总游戏时间
	LatestPlayTime int64            `json:"latest_play_time"`         // 最近一次游玩时间
}

type GamePlayTimeManager struct {
	PlayTimeMap map[string]*gameFolderPlayTime `json:"play_time_map"` // 文件夹 -> gameFolderPlayTime
}

func NewGamePlayTimeManager() *GamePlayTimeManager {
	return &GamePlayTimeManager{
		PlayTimeMap: make(map[string]*gameFolderPlayTime),
	}
}

func (manager *GamePlayTimeManager) addGamePlayTime(info gamePlayTimeInfo) {
	folderPlayTime, exists := manager.PlayTimeMap[info.FolderPath]
	if !exists {
		folderPlayTime = &gameFolderPlayTime{
			DailyPlayTime:  make(map[string]int64),
			LatestPlayTime: 0,
			TotalPlayTime:  0,
		}
		manager.PlayTimeMap[info.FolderPath] = folderPlayTime
	}

	folderPlayTime.TotalPlayTime += info.TotalPlayTime // 更新总游戏时间
}

// AddPlayTime 增加特定游戏在特定日期的聚焦时间
func (manager *GamePlayTimeManager) AddPlayTime(folderPath, exePath, date string, newPlayTime int64) {
	if _, ok := manager.PlayTimeMap[folderPath]; !ok {
		manager.PlayTimeMap[folderPath] = &gameFolderPlayTime{}
	}
	manager.PlayTimeMap[folderPath].DailyPlayTime[date] += newPlayTime
	manager.PlayTimeMap[folderPath].TotalPlayTime += newPlayTime
}

//// getPlayTime 获取特定游戏在特定日期的游玩时间
//func (manager *GamePlayTimeManager) getPlayTime(folderPath, exePath, date string) (int64, bool) {
//	if _, ok := manager.PlayTimeMap[folderPath]; !ok {
//		return 0, false
//	}
//	if _, ok := manager.PlayTimeMap[folderPath].ExePlayTimes[exePath]; !ok {
//		return 0, false
//	}
//	return manager.PlayTimeMap[folderPath].ExePlayTimes[exePath][date], true
//}
//
//// genGamesFoldersSlice 根据一个总游戏文件夹生成一个包含下面所有单个游戏文件夹的切片
//func genGamesFoldersSlice(baseGameFolder string) (gameFolders, error) {
//	log.Debugf("开始扫描目录 %s", baseGameFolder)
//	var gameFolders gameFolders
//	files, err := os.ReadDir(baseGameFolder)
//	if err != nil {
//		return nil, errors.Wrap(err, "读取目录失败")
//	}
//	for _, file := range files {
//		if file.IsDir() {
//			if strings.HasPrefix(file.Name(), ".") {
//				continue
//			}
//			log.Debugf("找到游戏文件夹 %s", filepath.Join(baseGameFolder, file.Name()))
//			gameFolders = append(gameFolders, filepath.Join(baseGameFolder, file.Name()))
//		}
//	}
//	return gameFolders, nil
//}
//
//func readGamePlayTimeFromFile(filePath string) (*GamePlayTimeManager, error) {
//	file, err := os.Open(filePath)
//	if err != nil {
//		return nil, errors.Wrap(err, "无法打开文件")
//	}
//	defer file.Close()
//	decoder := json.NewDecoder(file)
//	manager := NewGamePlayTimeManager()
//	err = decoder.Decode(manager)
//	if err != nil {
//		return nil, errors.Wrap(err, "无法解析json")
//	}
//	log.Debugf("成功从%s读取游戏时长数据", filePath)
//	return manager, nil
//}
//
//func writeGamePlayTimeToFile(manager *GamePlayTimeManager, filePath string) error {
//	err := os.MkdirAll(filepath.Dir(filePath), 0777)
//	if err != nil {
//		return errors.Wrap(err, "无法创建文件夹")
//	}
//	file, err := os.Create(filePath)
//	if err != nil {
//		return errors.Wrap(err, "无法创建文件")
//	}
//	defer file.Close()
//	// 以json格式写入，添加缩进
//	encoder := json.NewEncoder(file)
//	encoder.SetIndent("", "    ")
//	err = encoder.Encode(manager)
//	if err != nil {
//		return errors.Wrap(err, "无法写入json")
//	}
//	log.Debugf("成功保存游戏时长数据到%s", filePath)
//	return nil
//}

func WritePlayTimeToDB(libraryDB *database.SqliteGameLibrary, manager *GamePlayTimeManager) error {
	for folderPath, folderPlayTime := range manager.PlayTimeMap {
		playTimeMeta := types.GalgamePlayTime{
			TotalTime:   folderPlayTime.TotalPlayTime,
			LastTime:    folderPlayTime.LatestPlayTime,
			EachDayTime: folderPlayTime.DailyPlayTime,
		}
		err := libraryDB.InsertGamePlayTime("", folderPath, playTimeMeta)
		if err != nil {
			return errors.WithMessage(err, "写入数据库时发生错误")
		}
	}
	return nil
}

// GamePlayTimeMonitor 监控当前活动窗口，记录每个游戏的聚焦时间
func GamePlayTimeMonitor(libraryDB *database.SqliteGameLibrary) {
	log.Infof("开始监控活动窗口，游戏文件夹路径：%s", libraryDB.LibraryDir)
	var lastActiveWindowFolderPath string // 注意 是exe文件所属的文件夹路径
	var lastActiveWindowStartTime time.Time
	index, err := libraryDB.GetGameIndex()
	if err != nil {
		logrus.Errorf("无法获取游戏索引：%v", err)
		MonitorRunningStatusFlag = false
		return
	}
	var gamesFolders gameFolders
	for _, v := range index {
		gamesFolders = append(gamesFolders, v)
	}
	checkDateTimeTicker := time.NewTicker(30 * time.Second)
	saveGamePlayTimeTicker := time.NewTicker(15 * time.Second)
	today := time.Now().Format("2006-01-02")

	for {
		if MonitorStopFlag {
			MonitorRunningStatusFlag = false
			MonitorStopFlag = false
			log.Infof("游戏时长监控器停止运行！")
			return
		}
		select {
		case <-checkDateTimeTicker.C:
			if today != time.Now().Format("2006-01-02") {
				today = time.Now().Format("2006-01-02")
				log.Debugf("欢迎来到全新的一天：%s ！开始记录新的游戏时长了喵~", today)
			}
		case <-saveGamePlayTimeTicker.C:
			log.Debugf("开始保存游戏时长数据到数据库")
			err := WritePlayTimeToDB(libraryDB, GamePlayManager)
			if err != nil {
				log.Errorf("无法保存游戏时长数据：%v", err)
			}
		default:
		}
		hwnd, _ := getForegroundWindow()
		if hwnd == 0 {
			continue
		}

		pid, err := getWindowThreadProcessId(hwnd)
		exePath, err := getExecutablePath(pid)

		isInFolder, folderPath, err := isExePathInGamesFolder(exePath, gamesFolders)
		if err != nil {
			log.Errorf("无法检查可执行文件路径是否在游戏文件夹中：%v", err)
			continue
		}

		if isInFolder {
			// 如果可执行文件所属文件夹路径作为key变了，更新聚焦时间
			if folderPath != lastActiveWindowFolderPath {
				if lastActiveWindowFolderPath != "" {
					totalFocusTime := int64(time.Since(lastActiveWindowStartTime).Seconds())
					GamePlayManager.addGamePlayTime(gamePlayTimeInfo{
						FolderPath:     lastActiveWindowFolderPath,
						GameName:       filepath.Base(lastActiveWindowFolderPath),
						Date:           today,
						LatestPlayTime: time.Now().Unix(),
						TotalPlayTime:  totalFocusTime,
					})
					log.Debugf("游戏文件夹路径: %s, 今天共游玩时间: %d s", lastActiveWindowFolderPath, GamePlayManager.PlayTimeMap[lastActiveWindowFolderPath].TotalPlayTime)
					//logrus.Debugf("all gameFocusTimeMap: %v", gameFocusTimeMap)
				}
				lastActiveWindowFolderPath = folderPath
				lastActiveWindowStartTime = time.Now()
			}
		} else if lastActiveWindowFolderPath != "" {
			totalFocusTime := int64(time.Since(lastActiveWindowStartTime).Seconds())
			GamePlayManager.addGamePlayTime(gamePlayTimeInfo{
				FolderPath:     lastActiveWindowFolderPath,
				GameName:       filepath.Base(lastActiveWindowFolderPath),
				Date:           today,
				LatestPlayTime: time.Now().Unix(),
				TotalPlayTime:  totalFocusTime,
			})
			log.Debugf("游戏文件夹路径: %s, 今天共游玩时间: %d s", lastActiveWindowFolderPath, GamePlayManager.PlayTimeMap[lastActiveWindowFolderPath].TotalPlayTime)
			//logrus.Debugf("all gameFocusTimeMap: %v", gameFocusTimeMap)
			lastActiveWindowFolderPath = ""
		}

		time.Sleep(1 * time.Second)
	}
}

// isExePathInGamesFolder 检查给定的可执行文件路径是否位于GameFolders中的任意一个文件夹下
func isExePathInGamesFolder(exePath string, gameFolders gameFolders) (isInit bool, folderPath string, err error) {
	// 将exePath转换为绝对路径
	absExePath, err := filepath.Abs(exePath)
	if err != nil {
		return false, "", errors.Wrap(err, "无法获取可执行文件的绝对路径")
	}

	for _, gamesFolder := range gameFolders {
		// 确保游戏文件夹路径也是绝对路径
		absGamesFolder, _ := filepath.Abs(gamesFolder)

		// 使用filepath.Rel获取exePath相对于gamesFolder的相对路径
		relPath, _ := filepath.Rel(absGamesFolder, absExePath)
		if relPath == "" {
			continue
		}

		// 如果相对路径不以".."开始，且不是"."，则表示exePath在该gamesFolder或其子目录下
		if !strings.HasPrefix(relPath, "..") && relPath != "." && !strings.HasPrefix(relPath, "/") {
			return true, absGamesFolder, nil
		}
	}

	// 如果所有检查都未通过，则可执行文件不在任何一个游戏文件夹中
	return false, "", nil
}

// getForegroundWindow 使用Windows API获取当前前景窗口的句柄。
func getForegroundWindow() (uintptr, error) {
	ret, _, _ := procGetForegroundWindow.Call()
	return ret, nil
}

// getWindowThreadProcessId 使用Windows API获取给定窗口的进程ID。
func getWindowThreadProcessId(hwnd uintptr) (uint32, error) {
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid, nil
}

// getExecutablePath 使用Windows API获取进程的可执行文件路径。
func getExecutablePath(pid uint32) (string, error) {
	handle, _ := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
	defer syscall.CloseHandle(handle)

	var buf [syscall.MAX_PATH]uint16
	var size uint32 = syscall.MAX_PATH
	procQueryFullProcessImageName.Call(
		uintptr(handle),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	return syscall.UTF16ToString(buf[:]), nil
}