package scraper

import (
	"YoshinoGal/internal/scraper/sources"
	"YoshinoGal/internal/scraper/types"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
)

var ScrapeAllStatus = 0 // ScanGamesAndScrape 运行状态 0: 未运行 1: 运行中 2: 错误

var GamesScrapeStatusMap = map[string]int{} // 游戏的刮削进度（键为游戏文件夹路径） 0: 已完成 1: 运行中 2: 错误 3: 未开始

// readMetadata 从给定路径读取元数据
func readMetadata(path string) (*types.Galgame, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "读取元数据失败")
	}
	var metadata types.Galgame
	err = json.Unmarshal(file, &metadata)
	if err != nil {
		return nil, errors.Wrap(err, "解析元数据失败")
	}
	return &metadata, nil
}

// writeMetadata 将元数据写入给定路径
func writeMetadata(path string, metadata *types.Galgame) error {
	file, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return errors.Wrap(err, "序列化元数据失败")
	}
	return os.WriteFile(path, file, 0777)
}

func downloadImage(url string, path string) error {
	log.Infof("开始下载图片 %s 到 %s", url, path)
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "下载图片失败")
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "创建文件失败")
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "写入文件失败")
	}

	log.Debugf("下载图片 %s 到 %s 成功", url, path)
	return nil
}

// ScrapOneGame 刮削一个游戏（包含元数据、封面图、截图等所有步骤） directlyRun 为false时表示是ScanGamesAndScrape调用的，为true时表示是直接调用的
func ScrapOneGame(gameName string, priority []string, gameDir string) error {
	if log == nil {
		InitLogger()
	}
	var wg sync.WaitGroup
	var downloadErrors []error
	concLimiter := make(chan struct{}, 10) // 控制并发数量为10

	log.Infof("开始为游戏 %s 刮削元数据", gameName)
	log.Debugf("开始搜索游戏 %s 的元数据", gameName)
	log.Debugf("搜索优先级：%v", priority)
	GamesScrapeStatusMap[gameDir] = 1
	// 只有当元数据和封面图都存在时才认为已刮削过
	_, errMetadata := os.Stat(gameDir + "/metadata/metadata.json")
	_, errPoster := os.Stat(gameDir + "/metadata/poster.jpg")
	if !os.IsNotExist(errMetadata) && !os.IsNotExist(errPoster) {
		log.Infof("游戏 %s 已经存在元数据，跳过", gameName)
		GamesScrapeStatusMap[gameDir] = 0
		return nil
	}

	// 搜索元数据
	metadata, err := searchAndBuildMetadata(gameName, priority)
	if err != nil {
		GamesScrapeStatusMap[gameDir] = 2
		return errors.Wrap(err, "搜索游戏元数据时发生错误")
	}

	errsChan := make(chan error, len(metadata.ScreenshotsUrls)+1) // 用于收集下载错误

	// 写入元数据
	log.Debugf("为游戏 %s 搜索元数据成功，开始写入元数据到 %s", gameName, gameDir)
	metadataPath := gameDir + "/metadata/metadata.json"
	err = writeMetadata(metadataPath, &metadata)
	if err != nil {
		GamesScrapeStatusMap[gameDir] = 2
		return errors.Wrap(err, "写入游戏元数据时发生错误")
	}
	log.Debugf("为游戏 %s 写入元数据成功", gameName)

	// 下载封面图
	wg.Add(1)
	go func() {
		defer wg.Done()
		concLimiter <- struct{}{} // 占用一个并发位
		err := downloadImage(metadata.PosterUrl, gameDir+"/metadata/poster.jpg")
		if err != nil {
			errsChan <- errors.Wrap(err, "下载封面图时发生错误")
		}
		<-concLimiter // 释放并发位
	}()

	// 下载截图
	for i, url := range metadata.ScreenshotsUrls {
		wg.Add(1)
		go func(url string, i int) {
			defer wg.Done()
			concLimiter <- struct{}{} // 占用一个并发位
			err := downloadImage(url, gameDir+"/metadata/screenshot_"+strconv.Itoa(i)+".jpg")
			if err != nil {
				errsChan <- errors.Wrap(err, "下载截图时发生错误")
			}
			<-concLimiter // 释放并发位
		}(url, i)
	}

	// 等待所有下载完成
	go func() {
		wg.Wait()
		close(errsChan)
	}()

	for err := range errsChan {
		downloadErrors = append(downloadErrors, err)
	}
	if len(downloadErrors) > 0 {
		GamesScrapeStatusMap[gameDir] = 2
		return downloadErrors[0]
	}

	log.Infof("为游戏 %s 刮削元数据成功，好耶！", gameName)
	GamesScrapeStatusMap[gameDir] = 0
	return nil
}

// ScanGamesAndScrape 遍历指定目录，为每个还没刮削的游戏刮削
// 默认该目录下的所有一级子目录都是一个游戏目录 搜索时将以目录名作为游戏名
func ScanGamesAndScrape(directory string, priority []string) error {
	GamesScrapeStatusMap = map[string]int{}
	ScrapeAllStatus = 1
	if log == nil {
		InitLogger()
	}
	log.Infof("开始扫描目录 %s", directory)
	files, err := os.ReadDir(directory)
	if err != nil {
		ScrapeAllStatus = 2
		return errors.Wrap(err, "读取目录失败")
	}

	for _, file := range files {
		if file.IsDir() {
			GamesScrapeStatusMap[directory+"/"+file.Name()] = 3
		}
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		file := file
		gameName := file.Name()
		err = os.MkdirAll(directory+"/"+gameName+"/metadata", 0777)
		if err != nil {
			ScrapeAllStatus = 2
			return errors.Wrap(err, "为路径 "+directory+"/"+gameName+"/metadata 创建目录失败")
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := ScrapOneGame(gameName, priority, directory+"/"+gameName)
			if err != nil {
				log.Errorf("为游戏 %s 搜索元数据时发生错误：%s", gameName, err)
			}
		}()
	}

	wg.Wait()
	log.Infof("所有游戏刮削完成，好耶！")
	ScrapeAllStatus = 0
	return nil
}

// isZeroValue 判断一个值是否是零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan:
		return v.IsNil()
	default:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}

// mergeGalgameObjects 根据一个元数据来源优先级列表合并多个游戏元数据
func mergeGalgameObjects(priorityList []string, galgames map[string]types.Galgame) types.Galgame {
	var result types.Galgame
	for _, key := range priorityList {
		if galgame, exists := galgames[key]; exists {
			mergeGalgame(&result, galgame)
		}
	}
	return result
}

// mergeGalgame 将src中的非零值合并到dst中
func mergeGalgame(dst *types.Galgame, src types.Galgame) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src)
	for i := 0; i < dstVal.NumField(); i++ {
		dstField := dstVal.Field(i)
		srcField := srcVal.Field(i)
		if !isZeroValue(srcField) && (isZeroValue(dstField) || reflect.DeepEqual(dstField.Interface(), types.Galgame{})) {
			dstField.Set(srcField)
		}
	}
}

// // searchAndBuildMetadata 为给定游戏名搜索并构建元数据
func searchAndBuildMetadata(gameName string, priority []string) (types.Galgame, error) {
	var SearchFunctionsList = [1]func(gameName string) (map[string]types.Galgame, error){
		sources.SearchInVNDB,
	}
	resultsChan := make(chan map[string]types.Galgame, len(SearchFunctionsList))

	// 为每个数据源启动一个 goroutine 来执行搜索
	var wg sync.WaitGroup
	for _, Search := range SearchFunctionsList {
		wg.Add(1)
		Search := Search
		go func() {
			defer wg.Done()
			galgames, err := Search(gameName)
			if err != nil {
				log.Errorf("搜索游戏 %s 时发生错误：%s", gameName, err)
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
	var allResults = make(map[string]types.Galgame)
	for result := range resultsChan {
		for key, value := range result {
			allResults[key] = value
		}
	}

	// 合并搜索结果
	mergedResult := mergeGalgameObjects(priority, allResults)

	return mergedResult, nil
}
