package sources

import (
	"YoshinoGal/internal/scraper/structs"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func joinWithCommas(strs []string) string {
	var builder strings.Builder

	for i, str := range strs {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(str)
	}

	return builder.String()
}

// 限流器设置
var (
	// 每5分钟200个请求，计算出每个请求之间的最小间隔时间
	requestInterval = time.Minute * 5 / 200
	// 记录上一次请求的时间，初始化为当前时间减去间隔，确保第一个请求可以立即执行
	lastRequestTime = time.Now().Add(-requestInterval)
)

type SingleTitle struct {
	language string // 地区
	title    string //标题
	main     bool   // 是否是主标题
	official bool   // 是否被官方承认（？不确定
	latin    string // 标题的拉丁语表示
}

type SingleScreenshot struct {
	Sexual    float32 `json:"sexual"`
	Url       string  `json:"url"`
	VoteCount int     `json:"votecount"`
	Violence  float32 `json:"violence"`
}

// VNDBSingleGame 单个游戏搜索出的元数据
type VNDBSingleGame struct {
	Released      string             `json:"released"`
	Id            string             `json:"id"`
	Titles        []SingleTitle      `json:"titles"`
	DevStatus     int                `json:"devstatus"`
	Title         string             `json:"title"`
	Rating        float32            `json:"rating"`
	LengthMinutes int                `json:"length_minutes"`
	Screenshots   []SingleScreenshot `json:"screenshots"`
	Image         struct {
		Url string `json:"url"`
	} `json:"image"`
	Length      int    `json:"length"`
	Description string `json:"description"`
}

type VNDBSearchResponse struct {
	More    bool             `json:"more"` // 理论来说不需要考虑这个
	Results []VNDBSingleGame `json:"results"`
}

// ParseVNDBDescription VNDB简介中包含部分语法，解析为markdown
func ParseVNDBDescription(description string) string {
	replacements := map[string]string{
		`([cdprsu]v\d+(\.\d+)?)`:       `<a href="https://vndb.org/$1">$1</a>`,
		`(http[s]?://[^\s]+)`:          `<a href="$1">$1</a>`,
		`\[b\](.*?)\[/b\]`:             `<strong>$1</strong>`,
		`\[i\](.*?)\[/i\]`:             `<em>$1</em>`,
		`\[u\](.*?)\[/u\]`:             `<u>$1</u>`,
		`\[s\](.*?)\[/s\]`:             `<strike>$1</strike>`,
		`\[url=(.*?)\](.*?)\[/url\]`:   `<a href="$1">$2</a>`,
		`\[spoiler\](.*?)\[/spoiler\]`: `<span class="spoiler">$1</span>`,
		`\[quote\](.*?)\[/quote\]`:     `<blockquote>$1</blockquote>`,
		`\[raw\](.*?)\[/raw\]`:         `$1`,
		`\[code\](.*?)\[/code\]`:       `<pre><code>$1</code></pre>`,
	}

	for pattern, replacement := range replacements {
		re := regexp.MustCompile(pattern)
		description = re.ReplaceAllString(description, replacement)
	}

	return strings.ReplaceAll(description, "\n", "<br>")
}

func ConvertToGalgameStruct(VNDBResponse *VNDBSearchResponse) ([]structs.Galgame, error) {
	logrus.Debugf("总共搜索到了 %d 条游戏数据，开始尝试转换为内部源数据格式", len(VNDBResponse.Results))
	var galgames []structs.Galgame
	for _, g := range VNDBResponse.Results {
		var names []structs.GalgameName
		for _, t := range g.Titles {
			names = append(names, structs.GalgameName{
				Language: t.language,
				Title:    t.title,
				Main:     t.main,
				Official: t.official,
				Latin:    t.latin,
			})
		}
		var gameRating = structs.GalgameRating{
			VNDB: g.Rating,
		}
		var metaSources = structs.GalgameSource{
			VNDBUrl: "https://vndb.org/" + g.Id,
			VNDBID:  g.Id,
		}
		var screenshotsUrls []string
		for _, s := range g.Screenshots {
			if s.Sexual < 0.5 && s.Violence < 0.5 {
				screenshotsUrls = append(screenshotsUrls, s.Url)
			}
		}
		var description = ParseVNDBDescription(g.Description)
		var gal = structs.Galgame{
			Name:            g.Title,
			Names:           names,
			ReleaseDate:     g.Released,
			Score:           gameRating,
			PosterUrl:       g.Image.Url,
			Description:     description,
			MetadataSources: metaSources,
			LengthMinutes:   g.LengthMinutes,
			Length:          g.Length,
			DevStatus:       g.DevStatus,
			ScreenshotsUrls: screenshotsUrls,
		}
		//galgame := structs.GalgameSearchResult{
		//	Galgame: gal,
		//	Source:  "VNDB",
		//}
		galgames = append(galgames, gal)
		logrus.Debugf("成功转换了游戏数据：%s", g.Title)
		logrus.Debugf("转换后的数据：%+v", gal)
	}
	return galgames, nil
}

// VNDBSearch 对VNDB进行搜索并返回结果，topR为返回结果的数量
func VNDBSearch(gameName string, topR int) (map[string][]structs.Galgame, error) {
	if gameName == "" {
		return nil, errors.New("游戏名不能为空捏！")
	}
	if topR == 0 {
		topR = 1
	}
	// 确保请求遵守流控限制
	time.Sleep(time.Until(lastRequestTime.Add(requestInterval)))
	lastRequestTime = time.Now()
	var results []structs.Galgame

	logrus.Infof("正在从VNDB中搜索游戏：%s", gameName)
	// 不同的fields类型
	const (
		mainTitle     = "title"                                         // 主标题
		titles        = "titles{lang, title, latin, official, main}"    // 不同地区命名
		poster        = "image.url"                                     // 海报url
		lengthMinutes = "length_minutes"                                // 游玩预计时长
		length        = "length"                                        // 在没有具体游玩时间时获取这个 整形 1(非常短)-5(非常长)
		id            = "id"                                            //vndb id
		devStatus     = "devstatus"                                     // 整数，开发状态。 0 表示“已完成”，1 表示“正在开发”，2 表示“已取消”。
		released      = "released"                                      // 发布日期
		rating        = "rating"                                        // 评分，10-100 需自行转换为1-10
		screenshots   = "screenshots{url, sexual, violence, votecount}" // 游戏截图 将来作为游戏详情页的背景使用 选择优先级为优先选择投票数最高 且色情、暴力指数在0.5以下的
		description   = "description"                                   // TODO 简介 可能包含格式化代码，需要处理：https://vndb.org/d9#4
	)

	params := []string{mainTitle, titles, poster, lengthMinutes, length, id, devStatus, released, rating, screenshots, description}

	// 构建请求体
	requestBody := map[string]interface{}{
		"filters": []interface{}{"search", "=", gameName},
		"fields":  joinWithCommas(params),
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// 创建并发送请求
	//startTime := time.Now() // 请求开始时间
	req, _ := http.NewRequest("POST", "https://api.vndb.org/kana/vn", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "请求VNDB时发生错误喵！")
	}
	defer resp.Body.Close()

	// 我没搞懂这个限流是怎么回事，先注释掉
	//executionTime := time.Since(startTime)
	//if executionTime > time.Second {
	//	logrus.Warn("访问VNDB时超过了每分钟最大间隔！歇一会喵~")
	//}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "读取VNDB响应时发生错误")
	}

	var apiResponse VNDBSearchResponse

	err = json.Unmarshal(bodyBytes, &apiResponse)
	if err != nil {
		return nil, errors.Wrap(err, "解析VNDB响应时发生错误")
	}

	results, err = ConvertToGalgameStruct(&apiResponse)

	if err != nil {
		return nil, errors.Wrap(err, "转换VNDB响应时发生错误")
	}
	if len(results) == 0 {
		return nil, errors.New("没有搜索到相关游戏")
	}

	logrus.Debugf("VNDB搜索结果取前 %d 条", topR)
	results = results[:topR]
	return map[string][]structs.Galgame{"VNDB": results}, nil
}
