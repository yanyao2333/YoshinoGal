package scraper

import (
	"YoshinoGal/internal/scraper/sources"
	"YoshinoGal/internal/scraper/structs"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

var PRIORITY = []string{"VNDB"} // 搜索结果优先级

// ReadMetadata 从给定路径读取元数据
func ReadMetadata(path string) (*structs.Galgame, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "读取元数据失败")
	}
	var metadata structs.Galgame
	err = json.Unmarshal(file, &metadata)
	if err != nil {
		return nil, errors.Wrap(err, "解析元数据失败")
	}
	return &metadata, nil
}

// WriteMetadata 将元数据写入给定路径
func WriteMetadata(path string, metadata *structs.Galgame) error {
	file, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return errors.Wrap(err, "序列化元数据失败")
	}
	return os.WriteFile(path, file, 0644)
}

// ScanGames 遍历指定目录，为每个游戏执行搜索并更新元数据
func ScanGames(directory string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return errors.Wrap(err, "读取目录失败")
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if file.IsDir() {
			gameDir := filepath.Join(directory, file.Name())
			gameName := file.Name() // 假设目录名即游戏名

			wg.Add(1)
			go func(gameName, gameDir string) {
				defer wg.Done()
				searchAndWriteMetadata(gameName)
			}(gameName, gameDir)
		}
	}

	wg.Wait()
	return nil
}

// // searchAndWriteMetadata 对单个游戏执行搜索并更新元数据
func searchAndWriteMetadata(gameName string) {
	var SearchFunctionsList = [1]func(gameName string, topR int) (map[string][]structs.Galgame, error){
		sources.VNDBSearch,
	}
	resultsChan := make(chan map[string][]structs.Galgame, len(SearchFunctionsList))

	// 为每个数据源启动一个 goroutine 来执行搜索
	var wg sync.WaitGroup
	for _, Search := range SearchFunctionsList {
		wg.Add(1)
		Search := Search
		go func() {
			defer wg.Done()
			galgames, err := Search(gameName, 1)
			if err != nil {
				logrus.Errorf("搜索游戏 %s 时发生错误：%s", gameName, err)
				return
			}
			resultsChan <- galgames
		}()
	}

	// 等待所有的搜索完成
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集所有的搜索结果
	var allResults map[string][]structs.Galgame
	for results := range resultsChan {
		for source, games := range results {
			allResults[source] = append(allResults[source], games...)
		}
	}
}
