package models

// GalgameMetadataSources 定义游戏来源信息
type GalgameMetadataSources struct {
	BangumiID string // BangumiID
	CnGalID   string // CnGalID
	VNDBID    string // VNDBID
}

// GalgameRating 定义游戏评分信息
type GalgameRating struct {
	VNDB    float32 `json:"VNDB"`    // VNDB
	Bangumi float32 `json:"Bangumi"` // Bangumi
	CnGal   float32 `json:"CnGal"`   // CnGal
}

// GalgameName galgame在不同地区多个名称的单例
type GalgameName struct {
	Language string `json:"language"` // 地区
	Title    string `json:"title"`    //标题
	Main     bool   `json:"main"`     // 是否是主标题
	Official bool   `json:"official"` // 是否被官方承认（？不确定
	Latin    string `json:"latin"`    // 标题的拉丁语表示
}

// Developers 定义开发商信息
type Developers struct {
	Name        string `json:"name"`
	VNDBId      string `json:"vndb_id"`
	Description string `json:"description"`
}

// GalgameMetadata
type GalgameMetadata struct {
	Name                  string                 `json:"name"`                    // 主名称（即 GalNames中main official均为true的）
	Names                 []GalgameName          `json:"names"`                   // 多地区名称
	ReleaseDate           string                 `json:"release_date"`            // 发售日期
	Score                 GalgameRating          `json:"score"`                   // 评分
	RemotePosterUrl       string                 `json:"remote_poster_url"`       // 主海报地址
	Description           string                 `json:"description"`             // 简介
	MetadataSources       GalgameMetadataSources `json:"metadata_sources"`        // 元数据来源
	LengthMinutes         int                    `json:"length_minutes"`          // 游玩时长
	Length                int                    `json:"length"`                  // vndb特供字段，在没有具体时长的情况下使用
	DevStatus             int                    `json:"dev_status"`              // vndb特供字段，开发状态
	RemoteScreenshotsUrls []string               `json:"remote_screenshots_urls"` // 截图地址
	//LocalPosterPath       string                 `json:"local_poster_path"`       // 本地海报地址
	//LocalScreenshotsPaths []string               `json:"local_screenshots_paths"` // 本地截图地址
	Developers   []Developers     `json:"developers"` // 开发商
	IsR18        bool             `json:"is_r18"`     // 是否是R18游戏
	Tags         []string         `json:"tags"`       // 标签
	GalGameLocal GalgameLocalInfo `json:"local_info"` // 本地游戏信息
}

// GalgamePlayTime 用于记录游戏的游玩时长
type GalgamePlayTime struct {
	TotalTime      int64            `json:"total_time"`       // 总游玩时长
	LatestPlayTime int64            `json:"latest_play_time"` // 最后游玩时间
	EachDayTime    map[string]int64 `json:"each_day_time"`    // 每天游玩时长
}

// GalgameLocalInfo 本地游戏信息（包含截图、海报保存位置、游玩时长等等数据）
type GalgameLocalInfo struct {
	LocalPosterPath       string          `json:"local_poster_path"`       // 本地海报地址
	LocalScreenshotsPaths []string        `json:"local_screenshots_paths"` // 本地截图地址
	GameDirPath           string          `json:"game_dir_path"`           // 游戏目录地址
	PlayTime              GalgamePlayTime `json:"play_time"`               // 游玩时长
}

//// GalgameSearchResult 一个中间结果，用于搜索结果的排序
//type GalgameSearchResult struct {
//	GalgameMetadata
//	Source string // 搜索来源，用于后期根据优先级进行排序
//}
