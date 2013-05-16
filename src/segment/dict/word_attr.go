package dict

type WordAttr struct {
	Word      string
	Pos       int
	Frequency float64
}

func NewWordAttr(word string, pos int, frequency float64) *WordAttr {
	return &WordAttr{word, pos, frequency}
}
