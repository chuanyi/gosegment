package match

type MatchParameter struct {
	Redundancy       int // 多元分词冗余度
	UnknowRank       int // 未登录词权值
	BestRank         int // 最匹配词权值
	SecRank          int // 次匹配词权值
	ThirdRank        int // 再次匹配词权值
	SingleRank       int // 强行输出的单字的权值
	NumericRank      int // 数字的权值
	EnglishRank      int // 英文词汇权值
	EnglishLowerRank int // 英文词汇小写的权值
	EnglishStemRank  int // 英文词汇词根的权值
	SymbolRank       int // 符号的权值
	// SimplifiedTraditionalRank int // 强制同时输出简繁汉字时，非原来文本的汉字输出权值。比如原来文本是简体，这里就是输出的繁体字的权值，反之亦然。
	SynonymRank         int // 同义词权值
	WildcardRank        int // 通配符匹配结果的权值
	FilterEnglishLength int // 过滤英文选项生效时，过滤大于这个长度的英文。
	FilterNumericLength int // 过滤数字选项生效时，过滤大于这个长度的数字。
}

func NewMatchParameter() *MatchParameter {
	return &MatchParameter{Redundancy: 0, UnknowRank: 1, BestRank: 5, SecRank: 3, ThirdRank: 2, SingleRank: 1, NumericRank: 1, EnglishRank: 5, EnglishLowerRank: 3, EnglishStemRank: 2, SymbolRank: 1, SynonymRank: 1, WildcardRank: 1}
}
