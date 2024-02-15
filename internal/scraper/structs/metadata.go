package structs

// GalgameSource 定义游戏来源信息
type GalgameSource struct {
	BangumiUrl string // BangumiUrl
	VNDBUrl    string // VNDBUrl
	CnGalUrl   string // CnGalUrl
	VNDBID     string // VNDBID
}

// GalgameRating 定义游戏评分信息
type GalgameRating struct {
	VNDB    float32 // VNDB
	Bangumi float32 // Bangumi
	CnGal   float32 // CnGal
}

// GalgameName galgame在不同地区多个名称的单例
type GalgameName struct {
	Language string // 地区
	Title    string //标题
	Main     bool   // 是否是主标题
	Official bool   // 是否被官方承认（？不确定
	Latin    string // 标题的拉丁语表示
}

// Galgame 定义游戏信息
type Galgame struct {
	Name            string        // 主名称（即 GalNames中main official均为true的）
	Names           []GalgameName // 多地区名称
	ReleaseDate     string        // 发售日期
	Score           GalgameRating // 评分
	PosterUrl       string        // 主海报地址
	Description     string        // 简介
	MetadataSources GalgameSource // 元数据来源
	LengthMinutes   int           // 游玩时长
	Length          int           // vndb特供字段，在没有具体时长的情况下使用
	DevStatus       int           // vndb特供字段，开发状态
	ScreenshotsUrls []string      // 截图地址
}

//// GalgameSearchResult 一个中间结果，用于搜索结果的排序
//type GalgameSearchResult struct {
//	Galgame
//	Source string // 搜索来源，用于后期根据优先级进行排序
//}
