package match

type MatchOptions struct {
	ChineseNameIdentify        bool // 中文人名识别
	FrequencyFirst             bool // 词频优先
	MultiDimensionality        bool // 多元分词
	EnglishMultiDimensionality bool // 英文多元分词，这个开关，会将英文中的字母和数字分开
	FilterStopWords            bool // 过滤停用词
	IgnoreSpace                bool // 忽略空格、回车、Tab
	ForceSingleWord            bool // 强制一元分词
	// TraditionalChineseEnabled   bool // 繁体中文开关
	// OutputSimplifiedTraditional bool // 同时输出简体和繁体
	UnknownWordIdentify bool // 未登录词识别
	FilterEnglish       bool // 过滤英文，这个选项只有在过滤停用词选项生效时才有效
	FilterNumeric       bool // 过滤数字，这个选项只有在过滤停用词选项生效时才有效
	IgnoreCapital       bool // 忽略英文大小写
	EnglishSegment      bool // 英文分词
	SynonymOutput       bool // 同义词输出功能一般用于对搜索字符串的分词，不建议在索引时使用
	WildcardOutput      bool // 通配符匹配输出
	WildcardSegment     bool // 对通配符匹配的结果分词
}

func NewMatchOptions() *MatchOptions {
	return &MatchOptions{MultiDimensionality: true, FilterStopWords: true, IgnoreSpace: true, UnknownWordIdentify: true}
}
