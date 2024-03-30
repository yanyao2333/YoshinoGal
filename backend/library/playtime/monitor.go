package playtime

import (
	"YoshinoGal/backend/library/database"
	"YoshinoGal/backend/logging"
	"YoshinoGal/backend/models"
	"github.com/pkg/errors"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var log = logging.GetLogger()

var Manager *GamePlayTimeManager

var MonitorRunningStatus = false // 游戏时长监控器运行状态
var MonitorStopFlag = false      // 游戏时长监控器停止标志

func StartMonitor(libraryDB *database.SqliteGameLibrary) (Interface, error) {
	if Manager == nil {
		Manager = NewGamePlayTimeManager(libraryDB)
	}
	if MonitorRunningStatus {
		log.Warnf("游戏时长监控器已经在运行中了！")
		return Manager, nil
	}
	MonitorRunningStatus = true
	go Manager.Monitor()
	go Manager.saveGamePlayTime()
	return Manager, nil
}

// updatePlayTimeModel 游戏时长更新信息
type updatePlayTimeModel struct {
	FolderPath     string `json:"folder_path"`
	LatestPlayTime int64  `json:"latest_play_time"`
	Date           string `json:"date"`
	AddPlayTime    int64  `json:"add_play_time"`
}

type GamePlayTimeManager struct {
	UpdateChan chan updatePlayTimeModel
	libraryDB  *database.SqliteGameLibrary
}

func NewGamePlayTimeManager(libraryDB *database.SqliteGameLibrary) *GamePlayTimeManager {
	return &GamePlayTimeManager{
		UpdateChan: make(chan updatePlayTimeModel, 100),
		libraryDB:  libraryDB,
	}
}

func updateGamePlayTime(info updatePlayTimeModel, db *database.SqliteGameLibrary) error {
	id, err := db.GetGameIdFromPath(info.FolderPath)
	if err != nil {
		if errors.Is(err, models.CannotMatchGameIDFromPathInDatabase) {
			log.Warnf("游戏文件夹路径 %s 不在数据库中", info.FolderPath)
			return nil
		}
		return errors.WithMessage(err, "查询数据库时发生错误")
	}
	playTimeMeta, err := db.GetGamePlayTime(id)
	if err != nil {
		if errors.Is(err, models.GameNotFoundInDatabase) {
			log.Warnf("%v", err)
			return nil
		}
		return errors.WithMessage(err, "查询数据库时发生错误")
	}
	playTimeMeta.TotalTime += info.AddPlayTime
	playTimeMeta.LatestPlayTime = info.LatestPlayTime
	playTimeMeta.EachDayTime[info.Date] += info.AddPlayTime
	err = db.InsertGamePlayTime(id, *playTimeMeta)
	if err != nil {
		return errors.WithMessage(err, "写入数据库时发生错误")
	}
	return nil
}

func (manager *GamePlayTimeManager) StopMonitor() {
	MonitorStopFlag = true
}

type Interface interface {
	Monitor()
	saveGamePlayTime()
	StopMonitor()
}

func (manager *GamePlayTimeManager) saveGamePlayTime() {
	for {
		if MonitorStopFlag {
			MonitorRunningStatus = false
			MonitorStopFlag = false
			log.Debugf("游戏时长监控器停止运行，正在保存队列中剩余的游戏时长数据")
			for update := range manager.UpdateChan {
				err := updateGamePlayTime(update, manager.libraryDB)
				if err != nil {
					log.Errorf("无法更新游戏时长数据：%v", err)
				}
			}
		}
		select {
		case update := <-manager.UpdateChan:
			err := updateGamePlayTime(update, manager.libraryDB)
			if err != nil {
				log.Errorf("无法更新游戏时长数据：%v", err)
			}
		}
	}
}

// Monitor 监控当前活动窗口，记录游戏时长
func (manager *GamePlayTimeManager) Monitor() {
	var libraryDB = manager.libraryDB
	log.Infof("开始监控活动窗口，游戏文件夹路径：%s", libraryDB.LibraryDir)
	var lastActiveWindowFolderPath string
	var lastActiveWindowStartTime time.Time
	checkDateTimeTicker := time.NewTicker(30 * time.Second)
	today := time.Now().Format("2006-01-02")

	for {
		if MonitorStopFlag {
			MonitorRunningStatus = false
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
		default:
		}
		hwnd, _ := getForegroundWindow()
		if hwnd == 0 {
			continue
		}

		pid, err := getWindowThreadProcessId(hwnd)
		exePath, err := getExecutablePath(pid)

		isSubPath, err := IsSubpath(libraryDB.LibraryDir, exePath)
		if err != nil {
			log.Errorf("无法判断文件夹是否为library的子文件夹：%v", err)
			continue
		}
		folderPath := filepath.Dir(exePath)

		if isSubPath {
			if folderPath != lastActiveWindowFolderPath {
				if lastActiveWindowFolderPath != "" {
					totalFocusTime := int64(time.Since(lastActiveWindowStartTime).Seconds()) // 计算上一个游戏的时间
					manager.UpdateChan <- updatePlayTimeModel{
						FolderPath:     lastActiveWindowFolderPath,
						LatestPlayTime: time.Now().Unix(),
						Date:           today,
						AddPlayTime:    totalFocusTime,
					}
					if err != nil {
						log.Errorf("无法更新游戏时长数据：%v", err)
					}
					log.Debugf("游戏时长更新 游戏路径：%s", lastActiveWindowFolderPath)
				}
				lastActiveWindowFolderPath = folderPath
				lastActiveWindowStartTime = time.Now()
			}
		} else if lastActiveWindowFolderPath != "" {
			totalFocusTime := int64(time.Since(lastActiveWindowStartTime).Seconds())
			manager.UpdateChan <- updatePlayTimeModel{
				FolderPath:     lastActiveWindowFolderPath,
				LatestPlayTime: time.Now().Unix(),
				Date:           today,
				AddPlayTime:    totalFocusTime,
			}
			if err != nil {
				log.Errorf("无法更新游戏时长数据：%v", err)
			}
			log.Debugf("游戏时长更新 游戏路径：%s", lastActiveWindowFolderPath)
			lastActiveWindowFolderPath = ""
		}

		time.Sleep(1 * time.Second)
	}
}

// IsSubpath 检查subpath是否是path的子路径
func IsSubpath(path, subpath string) (bool, error) {
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return false, err
	}
	absSubpath, err := filepath.Abs(filepath.Clean(subpath))
	if err != nil {
		return false, err
	}

	if !strings.HasPrefix(absSubpath, absPath) {
		return false, nil
	}

	rel, err := filepath.Rel(absPath, absSubpath)
	if err != nil {
		return false, err
	}

	if strings.HasPrefix(rel, "..") {
		return false, nil
	}

	return true, nil
}

// 以下为Windows API调用部分

var (
	user32                        = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow       = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId  = user32.NewProc("GetWindowThreadProcessId")
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
)

// getForegroundWindow 使用Windows API获取当前前景窗口的句柄。
func getForegroundWindow() (uintptr, error) {
	ret, _, _ := procGetForegroundWindow.Call()
	return ret, nil
}

// getWindowThreadProcessId 使用Windows API获取给定窗口的进程ID。
func getWindowThreadProcessId(hwnd uintptr) (uint32, error) {
	var pid uint32
	_, _, _ = procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid, nil
}

// getExecutablePath 使用Windows API获取进程的可执行文件路径。
func getExecutablePath(pid uint32) (string, error) {
	handle, _ := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
	defer syscall.CloseHandle(handle)

	var buf [syscall.MAX_PATH]uint16
	var size uint32 = syscall.MAX_PATH
	_, _, _ = procQueryFullProcessImageName.Call(
		uintptr(handle),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	return syscall.UTF16ToString(buf[:]), nil
}

// Windows API调用部分结束

//type gameFolders []string
//
//type gamePlayTimeInfo struct {
//	FolderPath     string `json:"folder_path"`
//	GameName       string `json:"game_name"`
//	LatestPlayTime int64  `json:"latest_play_time"`
//	Date           string `json:"date"`
//	TotalPlayTime  int64  `json:"total_play_time"`
//}
//
//type gameFolderPlayTime struct {
//	DailyPlayTime  map[string]int64 `json:"daily_play_time"`  // 日期 -> 游戏时间
//	TotalPlayTime  int64            `json:"total_play_time"`  // 该游戏总游戏时间
//	LatestPlayTime int64            `json:"latest_play_time"` // 最近一次游玩时间
//}
//

//
//func NewGamePlayTimeManager() *GamePlayTimeManager {
//	return &GamePlayTimeManager{
//		PlayTimeMap: make(map[string]*gameFolderPlayTime),
//	}
//}

//func (manager *GamePlayTimeManager) addGamePlayTime(info gamePlayTimeInfo, libraryDB *database.SqliteGameLibrary) error {
//	folderPlayTime, exists := manager.PlayTimeMap[info.FolderPath]
//	if !exists {
//		err := LoadOneGamePlayTimeFromDB(libraryDB, info.FolderPath)
//		if err != nil {
//			return errors.WithMessage(err, "无法从数据库加载游戏时长数据")
//		}
//		folderPlayTime, exists := manager.PlayTimeMap[info.FolderPath]
//		if !exists {
//			folderPlayTime = &gameFolderPlayTime{
//				DailyPlayTime:  make(map[string]int64),
//				LatestPlayTime: 0,
//				TotalPlayTime:  0,
//			}
//		}
//		manager.PlayTimeMap[info.FolderPath] = folderPlayTime
//		manager.PlayTimeMap[info.FolderPath].DailyPlayTime[info.Date] += info.TotalPlayTime
//		manager.PlayTimeMap[info.FolderPath].TotalPlayTime += info.TotalPlayTime
//		manager.PlayTimeMap[info.FolderPath].LatestPlayTime = info.LatestPlayTime
//		return nil
//	}
//
//	folderPlayTime.TotalPlayTime += info.TotalPlayTime            // 更新总游戏时间
//	folderPlayTime.LatestPlayTime = info.LatestPlayTime           // 更新最近一次游玩时间
//	folderPlayTime.DailyPlayTime[info.Date] += info.TotalPlayTime // 更新当天游戏时间
//	manager.PlayTimeMap[info.FolderPath] = folderPlayTime
//	return nil
//}

// AddPlayTime 增加特定游戏在特定日期的聚焦时间
//func (manager *GamePlayTimeManager) AddPlayTime(folderPath, date string, newPlayTime int64) {
//	if _, ok := manager.PlayTimeMap[folderPath]; !ok {
//		manager.PlayTimeMap[folderPath] = &gameFolderPlayTime{}
//	}
//	manager.PlayTimeMap[folderPath].DailyPlayTime[date] += newPlayTime
//	manager.PlayTimeMap[folderPath].TotalPlayTime += newPlayTime
//}

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

//func LoadOneGamePlayTimeFromDB(libraryDB *database.SqliteGameLibrary, folderPath string) error {
//	id, err := libraryDB.GetGameIdFromPath(folderPath)
//	if err != nil {
//		if errors.Is(err, models.CannotMatchGameIDFromPath) {
//			log.Warnf("游戏文件夹路径%s不在数据库中", folderPath)
//			return nil
//		}
//		return errors.WithMessage(err, "查询数据库时发生错误")
//	}
//	playTimeMeta, err := libraryDB.GetGamePlayTime(id)
//	if err != nil {
//		return errors.WithMessage(err, "查询数据库时发生错误")
//	}
//	var gameFolderPlayTime = &gameFolderPlayTime{
//		DailyPlayTime:  playTimeMeta.EachDayTime,
//		TotalPlayTime:  playTimeMeta.TotalTime,
//		LatestPlayTime: playTimeMeta.LatestPlayTime,
//	}
//	GamePlayManager.PlayTimeMap[folderPath] = gameFolderPlayTime
//	log.Debugf("成功从数据库加载游戏时长数据：%v", gameFolderPlayTime)
//	return nil
//}
//
//func WritePlayTimeToDB(libraryDB *database.SqliteGameLibrary, manager *GamePlayTimeManager) error {
//	for folderPath, folderPlayTime := range manager.PlayTimeMap {
//		playTimeMeta := models.GalgamePlayTime{
//			TotalTime:   folderPlayTime.TotalPlayTime,
//			LatestPlayTime:    folderPlayTime.LatestPlayTime,
//			EachDayTime: folderPlayTime.DailyPlayTime,
//		}
//		id, err := libraryDB.GetGameIdFromPath(folderPath)
//		if err != nil {
//			if errors.Is(err, models.CannotMatchGameIDFromPath) {
//				log.Warnf("游戏文件夹路径%s不在数据库中", folderPath)
//				continue
//			}
//			return errors.WithMessage(err, "查询数据库时发生错误")
//		}
//		err = libraryDB.InsertGamePlayTime(id, playTimeMeta)
//		if err != nil {
//			return errors.WithMessage(err, "写入数据库时发生错误")
//		}
//	}
//	return nil
//}
