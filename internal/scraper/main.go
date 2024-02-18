package scraper

import (
	"YoshinoGal/internal/scraper/sources"
	"YoshinoGal/internal/scraper/types"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
)

// ReadMetadata 从给定路径读取元数据
func ReadMetadata(path string) (*types.Galgame, error) {
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

// WriteMetadata 将元数据写入给定路径
func WriteMetadata(path string, metadata *types.Galgame) error {
	file, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return errors.Wrap(err, "序列化元数据失败")
	}
	return os.WriteFile(path, file, 0777)
}

func DownloadImage(url string, path string) error {
	logrus.Infof("开始下载图片 %s 到 %s", url, path)
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

	logrus.Debugf("下载图片 %s 到 %s 成功", url, path)
	return nil
}

// ScrapOneGame 刮削一个游戏（包含元数据、封面图、截图等所有步骤）
func ScrapOneGame(gameName string, priority []types.GalgameMetadataSource, gameDir string) error {
	var wg sync.WaitGroup
	var downloadErrors []error
	concLimiter := make(chan struct{}, 10) // 控制并发数量为10

	logrus.Infof("开始为游戏 %s 刮削元数据", gameName)
	logrus.Debugf("开始搜索游戏 %s 的元数据", gameName)
	logrus.Debugf("搜索优先级：%v", priority)
	// 只有当元数据和封面图都存在时才认为已刮削过
	_, errMetadata := os.Stat(gameDir + "/metadata/metadata.json")
	_, errPoster := os.Stat(gameDir + "/metadata/poster.jpg")
	if !os.IsNotExist(errMetadata) && !os.IsNotExist(errPoster) {
		logrus.Infof("游戏 %s 已经存在元数据，跳过", gameName)
		return nil
	}

	// 搜索元数据
	metadata, err := searchAndBuildMetadata(gameName, priority)
	if err != nil {
		return errors.Wrap(err, "搜索游戏元数据时发生错误")
	}

	errsChan := make(chan error, len(metadata.ScreenshotsUrls)+1) // 用于收集下载错误

	// 写入元数据
	logrus.Debugf("为游戏 %s 搜索元数据成功，开始写入元数据到 %s", gameName, gameDir)
	metadataPath := gameDir + "/metadata/metadata.json"
	err = WriteMetadata(metadataPath, &metadata)
	if err != nil {
		return errors.Wrap(err, "写入游戏元数据时发生错误")
	}
	logrus.Debugf("为游戏 %s 写入元数据成功", gameName)

	// 下载封面图
	wg.Add(1)
	go func() {
		defer wg.Done()
		concLimiter <- struct{}{} // 占用一个并发位
		err := DownloadImage(metadata.PosterUrl, gameDir+"/metadata/poster.jpg")
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
			err := DownloadImage(url, gameDir+"/metadata/screenshot"+strconv.Itoa(i)+".jpg")
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
		return downloadErrors[0]
	}

	logrus.Infof("为游戏 %s 刮削元数据成功，好耶！", gameName)
	return nil
}

// ScanGamesAndScrape 遍历指定目录，为每个还没刮削的游戏刮削
// 默认该目录下的所有一级子目录都是一个游戏目录 搜索时将以目录名作为游戏名
func ScanGamesAndScrape(directory string, priority []types.GalgameMetadataSource) error {
	logrus.Infof("开始扫描目录 %s", directory)
	files, err := os.ReadDir(directory)
	if err != nil {
		return errors.Wrap(err, "读取目录失败")
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
			return errors.Wrap(err, "为路径 "+directory+"/"+gameName+"/metadata 创建目录失败")
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := ScrapOneGame(gameName, priority, directory+"/"+gameName)
			if err != nil {
				logrus.Errorf("为游戏 %s 搜索元数据时发生错误：%s", gameName, err)
			}
		}()
	}

	wg.Wait()
	logrus.Infof("所有游戏刮削完成，好耶！")
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

// MergeGalgameObjects 根据一个元数据来源优先级列表合并多个游戏元数据
func MergeGalgameObjects(priorityList []types.GalgameMetadataSource, galgames map[types.GalgameMetadataSource]types.Galgame) types.Galgame {
	var result types.Galgame
	for _, key := range priorityList {
		if galgame, exists := galgames[key]; exists {
			MergeGalgame(&result, galgame)
		}
	}
	return result
}

// MergeGalgame 将src中的非零值合并到dst中
func MergeGalgame(dst *types.Galgame, src types.Galgame) {
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
func searchAndBuildMetadata(gameName string, priority []types.GalgameMetadataSource) (types.Galgame, error) {
	var SearchFunctionsList = [1]func(gameName string) (map[types.GalgameMetadataSource]types.Galgame, error){
		sources.VNDBSearch,
	}
	resultsChan := make(chan map[types.GalgameMetadataSource]types.Galgame, len(SearchFunctionsList))

	// 为每个数据源启动一个 goroutine 来执行搜索
	var wg sync.WaitGroup
	for _, Search := range SearchFunctionsList {
		wg.Add(1)
		Search := Search
		go func() {
			defer wg.Done()
			galgames, err := Search(gameName)
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
	var allResults = make(map[types.GalgameMetadataSource]types.Galgame)
	for result := range resultsChan {
		for key, value := range result {
			allResults[key] = value
		}
	}

	// 合并搜索结果
	mergedResult := MergeGalgameObjects(priority, allResults)

	return mergedResult, nil
}
