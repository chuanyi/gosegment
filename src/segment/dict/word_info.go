package dict

const (
	TNone              = 0
	TEnglish           = 1
	TSimplifiedChinese = 2
	//TTraditionalChinese = 3
	TNumeric = 4
	TSymbol  = 5
	TSpace   = 6
	TSynonym = 7 //同义词
)

type WordInfo struct {
	Word             string
	Pos              int
	Frequency        float64
	WordType         int
	OriginalWordType int
	Position         int
	Rank             int
}

func NewWordInfo(word string, position int, pos int, frequency float64, rank int, wordType int, originalWordType int) *WordInfo {
	return &WordInfo{Word: word, Pos: pos, Frequency: frequency, WordType: wordType, OriginalWordType: originalWordType, Position: position, Rank: rank}
}

func NewWordInfoDefault() *WordInfo {
	return &WordInfo{}
}

func NewWordInfoSome(word string, pos int, frequency float64) *WordInfo {
    return &WordInfo{Word: word, Pos: pos, Frequency: frequency}
}
