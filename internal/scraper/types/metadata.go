package types

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

// Galgame 定义游戏信息
type Galgame struct {
	Name            string                 `json:"name"`             // 主名称（即 GalNames中main official均为true的）
	Names           []GalgameName          `json:"names"`            // 多地区名称
	ReleaseDate     string                 `json:"release_date"`     // 发售日期
	Score           GalgameRating          `json:"score"`            // 评分
	PosterUrl       string                 `json:"poster_url"`       // 主海报地址
	Description     string                 `json:"description"`      // 简介
	MetadataSources GalgameMetadataSources `json:"metadata_sources"` // 元数据来源
	LengthMinutes   int                    `json:"length_minutes"`   // 游玩时长
	Length          int                    `json:"length"`           // vndb特供字段，在没有具体时长的情况下使用
	DevStatus       int                    `json:"dev_status"`       // vndb特供字段，开发状态
	ScreenshotsUrls []string               `json:"screenshots_urls"` // 截图地址
}

//// GalgameSearchResult 一个中间结果，用于搜索结果的排序
//type GalgameSearchResult struct {
//	Galgame
//	Source string // 搜索来源，用于后期根据优先级进行排序
//}
