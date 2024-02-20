package game_play_time_monitor

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                        = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow       = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId  = user32.NewProc("GetWindowThreadProcessId")
	kernel32                      = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
)

type gameFolders []string

// genGamesFoldersSlice 根据一个总游戏文件夹生成一个包含下面所有单个游戏文件夹的切片
func genGamesFoldersSlice(baseGameFolder string) (gameFolders, error) {
	logrus.Debugf("开始扫描目录 %s", baseGameFolder)
	var gameFolders gameFolders
	files, err := os.ReadDir(baseGameFolder)
	if err != nil {
		return nil, errors.Wrap(err, "读取目录失败")
	}
	for _, file := range files {
		if file.IsDir() {
			logrus.Debugf("找到游戏文件夹 %s", filepath.Join(baseGameFolder, file.Name()))
			gameFolders = append(gameFolders, filepath.Join(baseGameFolder, file.Name()))
		}
	}
	return gameFolders, nil
}

// monitorActiveWindows 监控当前活动窗口，记录每个游戏的聚焦时间
func monitorActiveWindows(gamesFolders gameFolders) {
	activeWindowTimeMap := make(map[string]int64)
	var lastActiveWindowExePath string
	var lastActiveWindowStartTime time.Time

	for {
		hwnd, _ := getForegroundWindow()
		if hwnd == 0 {
			continue
		}

		pid, err := getWindowThreadProcessId(hwnd)
		exePath, err := getExecutablePath(pid)

		isInFolder, err := isExePathInGamesFolder(exePath, gamesFolders)
		if err != nil {
			logrus.Errorf("无法检查可执行文件路径是否在游戏文件夹中：%v", err)
			continue
		}

		if isInFolder {
			// 如果可执行文件路径作为key变了，更新聚焦时间
			if exePath != lastActiveWindowExePath {
				if lastActiveWindowExePath != "" {
					totalFocusTime := int64(time.Since(lastActiveWindowStartTime).Seconds())
					activeWindowTimeMap[lastActiveWindowExePath] += totalFocusTime
					logrus.Debugf("ExePath: %s, Total Focus Time: %d seconds", lastActiveWindowExePath, activeWindowTimeMap[lastActiveWindowExePath])
					logrus.Debugf("all activeWindowTimeMap: %v", activeWindowTimeMap)
				}
				lastActiveWindowExePath = exePath
				lastActiveWindowStartTime = time.Now()
			}
		} else if lastActiveWindowExePath != "" {
			totalFocusTime := int64(time.Since(lastActiveWindowStartTime).Seconds())
			activeWindowTimeMap[lastActiveWindowExePath] += totalFocusTime
			logrus.Debugf("ExePath: %s, Total Focus Time: %d seconds", lastActiveWindowExePath, activeWindowTimeMap[lastActiveWindowExePath])
			logrus.Debugf("all activeWindowTimeMap: %v", activeWindowTimeMap)
			lastActiveWindowExePath = ""
		}

		time.Sleep(1 * time.Second)
	}
}

// isExePathInGamesFolder 检查给定的可执行文件路径是否位于GameFolders中的任意一个文件夹下
func isExePathInGamesFolder(exePath string, gameFolders gameFolders) (bool, error) {
	// 将exePath转换为绝对路径
	absExePath, err := filepath.Abs(exePath)
	if err != nil {
		return false, errors.Wrap(err, "无法获取可执行文件的绝对路径")
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
			return true, nil
		}
	}

	// 如果所有检查都未通过，则可执行文件不在任何一个游戏文件夹中
	return false, nil
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
