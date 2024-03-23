package watchdog

import (
	"YoshinoGal/backend/library/database"
	"YoshinoGal/backend/library/scraper"
	"YoshinoGal/backend/logging"
	"YoshinoGal/backend/models"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

var log = logging.GetLogger()

type Interface interface {
	Watch()
	AddWatchDir(dir string)
	RemoveWatchDir(dir string)
	StartWatchGame(gameDir string, lib *database.SqliteGameLibrary, scraperPriority []string)
	Close()
}

type WatchDog struct {
	Watcher          *fsnotify.Watcher
	NewFolderChan    chan string
	RemoveFolderChan chan string
	RenameFolderChan chan string
}

func NewWatchDog() *WatchDog {
	return &WatchDog{
		NewFolderChan:    make(chan string),
		RemoveFolderChan: make(chan string),
		RenameFolderChan: make(chan string),
	}
}

func (w *WatchDog) Close() {
	w.Watcher.Close()
}

// Watch 启动fsnotify的监视，并将事件发送到相应的channel
func (w *WatchDog) Watch() {
	w.Watcher, _ = fsnotify.NewWatcher()
	defer w.Watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-w.Watcher.Events:
				if !ok {
					return
				}
				log.Infof("event: %v", event)
				if strings.Contains(event.Name, ".") {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Error(err)
					}
					if fileInfo.IsDir() {
						w.NewFolderChan <- event.Name
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Error(err)
					}
					if fileInfo.IsDir() {
						w.RemoveFolderChan <- event.Name
					}

				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Error(err)
					}
					if fileInfo.IsDir() {
						w.RenameFolderChan <- event.Name
					}
				}
			case err, ok := <-w.Watcher.Errors:
				if !ok {
					return
				}
				log.Infof("error: %v", err)
			}
		}
	}()
}

func (w *WatchDog) AddWatchDir(dir string) {
	w.Watcher.Add(dir)
}

func (w *WatchDog) RemoveWatchDir(dir string) {
	w.Watcher.Remove(dir)
}

// StartWatchGame 开始监视游戏目录
func (w *WatchDog) StartWatchGame(gameDir string, lib *database.SqliteGameLibrary, scraperPriority []string) {
	log := logging.GetLogger()
	log.Infof("开始监视游戏目录 %s 文件修改动态", gameDir)
	w.AddWatchDir(gameDir)
	go w.Watch()
	for {
		select {
		case newFolder := <-w.NewFolderChan:
			log.Infof("发现新游戏文件夹 %s", newFolder)
			gameName := filepath.Base(newFolder)
			go func() {
				err := scraper.ScrapOneGame(gameName, scraperPriority, gameDir, true, lib)
				if err != nil {
					log.Error(err)
				}
			}()
		case rmFolder := <-w.RemoveFolderChan:
			log.Infof("游戏文件夹 %s 被删除，同步删除数据库记录", rmFolder)
			gid, err := lib.GetGameIdFromPath(rmFolder)
			if err != nil {
				if errors.Is(err, models.CannotMatchGameIDFromPathInDatabase) {
					log.Infof("游戏 %s 不存在于数据库中", rmFolder)
				} else {
					log.Error(err)
				}
			}
			err = lib.RemoveGame(gid)
			if err != nil {
				log.Error(err)
			}
		case renameFolder := <-w.RenameFolderChan:
			log.Infof("游戏文件夹 %s 被重命名，开始更新数据库", renameFolder)
			// TODO 重命名操作会带来rename和create两个事件，但通过前者无法获取到新的文件夹名，因此需要通过后者获取新的文件夹名，但由于两者并不关联，所以现阶段无法实现，只能删除后重新生成
			gid, err := lib.GetGameIdFromPath(renameFolder)
			if err != nil {
				if errors.Is(err, models.CannotMatchGameIDFromPathInDatabase) {
					log.Infof("游戏 %s 不存在于数据库中", renameFolder)
				} else {
					log.Error(err)
				}
			}
			err = lib.RemoveGame(gid)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

// WatchGame 监视游戏目录
func WatchGame(gameDir string, lib *database.SqliteGameLibrary, scraperPriority []string) (Interface, error) {
	watchdog := NewWatchDog()
	go watchdog.StartWatchGame(gameDir, lib, scraperPriority)
	return watchdog, nil
}
