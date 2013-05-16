package match

import (
	"container/list"
	"segment/dict"
	"segment/utils"
	"sort"
)

type IChsFullTextMatch interface {
	SetOptionParams(options *MatchOptions, params *MatchParameter)
	Match(posLenArr []dict.PositionLength, originalText string, count int) *list.List
}

const (
	TopRecord      int = 3
	SingleWordMask int = dict.POS_D_C | dict.POS_D_P | dict.POS_D_R | dict.POS_D_U
)

type ChsFullTextMatch struct {
	options         *MatchOptions
	params          *MatchParameter
	wordDict        *dict.WordDictionary
	root            *Node
	leafNodeList    []*Node
	posLenArr       []dict.PositionLength
	inputStringLen  int
	allCombinations []([]dict.PositionLength)
}

var freqFirst bool

func NewChsFullTextMatch(wdict *dict.WordDictionary) (m *ChsFullTextMatch) {
	m = &ChsFullTextMatch{wordDict: wdict}
	m.root = NewNode()
	m.leafNodeList = [](*Node){}
	m.allCombinations = make([]([]dict.PositionLength), 0)
	return
}

func (m *ChsFullTextMatch) SetOptionParams(options *MatchOptions, params *MatchParameter) {
	m.options = options
	m.params = params 
}

func (m *ChsFullTextMatch) Match(posLenArr []dict.PositionLength, originalText string) *list.List {
	if m.options == nil {
		m.options = NewMatchOptions()
	}
	if m.params == nil {
		m.params = NewMatchParameter()
	}
	runes := utils.ToRunes(originalText)
	masks := make([]int, len(runes))
	redundancy := m.params.Redundancy

	result := list.New()
	if len(posLenArr) == 0 {
		if m.options.UnknownWordIdentify {
			wi := dict.NewWordInfoDefault()
			wi.Word = originalText
			wi.Position = 0
			wi.WordType = dict.TNone
			wi.Rank = 1
			result.PushFront(wi)
			return result
		} else {
			position := 0
			for _, r := range runes {
				wi := dict.NewWordInfoDefault()
				wi.Word = string(r)
				wi.Position = position
				wi.WordType = dict.TNone
				wi.Rank = 1
				position++
				result.PushBack(wi)
			}
			return result
		}
	}

	leafNodeArray := m.getLeafNodeArray(posLenArr, originalText)

	// 获取前TopRecord个单词序列
	j := 0
	for _, node := range leafNodeArray {
		if leafNodeArray[j] == nil {
			break
		}
		if j >= TopRecord || j >= len(leafNodeArray) {
			break
		}
		comb := make([]dict.PositionLength, node.AboveCount)
		i := node.AboveCount - 1
		cur := node
		for i >= 0 {
			comb[i] = cur.PosLen
			cur = cur.Parent
			i--
		}
		m.allCombinations = append(m.allCombinations, comb)
		j++
	}

	// Force single word
	// 强制一元分词
	if m.options.ForceSingleWord {
		comb := make([]dict.PositionLength, len(runes))
		for i := 0; i < len(comb); i++ {
			pl := dict.NewPositionLength(i, 1, dict.NewWordAttr(string(runes[i]), dict.POS_UNK, 0.0))
			pl.Level = 3
			comb[i] = pl
		}
		m.allCombinations = append(m.allCombinations, comb)
	}

	if len(m.allCombinations) > 0 {
		positionCollection := m.mergeAllCombinations(redundancy)
		curPc := positionCollection.Front()
		for curPc != nil {
			pl := curPc.Value.(dict.PositionLength)
			wi := dict.NewWordInfoDefault()
			wi.Word = string(runes[pl.Position:(pl.Position + pl.Length)])
			wi.Pos = pl.WordAttri.Pos
			wi.Frequency = pl.WordAttri.Frequency
			wi.WordType = dict.TSimplifiedChinese
			wi.Position = pl.Position
			switch pl.Level {
			case 0:
				wi.Rank = m.params.BestRank
			case 1:
				wi.Rank = m.params.SecRank
			case 2:
				wi.Rank = m.params.ThirdRank
			case 3:
				wi.Rank = m.params.SingleRank
			default:
				wi.Rank = m.params.BestRank
			}

			result.PushBack(wi)
			if pl.Length > 1 {
				for k := pl.Position; k < pl.Position+pl.Length; k++ {
					masks[k] = 2
				}
			} else {
				masks[pl.Position] = 1
			}
			curPc = curPc.Next()
		}
	}

	// 合并未登录词
	unknownWords, needRemoveSingleWord := m.getUnknownWords(masks, runes)
	// 合并结果序列到对应位置中
	if len(unknownWords) > 0 {
		cur := result.Front()
		if needRemoveSingleWord && !m.options.ForceSingleWord {
			// remove single word need be removed
			for cur != nil {
				if utils.RuneLen(cur.Value.(*dict.WordInfo).Word) == 1 {
					if masks[cur.Value.(*dict.WordInfo).Position] == 11 {
						removeItem := cur
						cur = cur.Next()
						result.Remove(removeItem)
						continue
					}
				}
				cur = cur.Next()
			}
		}

		cur = result.Front()
		j = 0
		for cur != nil {
			if cur.Value.(*dict.WordInfo).Position >= unknownWords[j].Position {
				result.InsertBefore(unknownWords[j], cur)
				j++
				if j >= len(unknownWords) {
					break
				}
			}

			if cur.Value.(*dict.WordInfo).Position < unknownWords[j].Position {
				cur = cur.Next()
			}
		}

		for j < len(unknownWords) {
			result.PushBack(unknownWords[j])
			j++
		}
	}

	return result
}

func (m *ChsFullTextMatch) getLeafNodeArray(posLenArr []dict.PositionLength, originalText string) [](*Node) {
	// Split by isolated point
	result := make([]*Node, TopRecord)
	lastRightBoundary := posLenArr[0].Position + posLenArr[0].Length
	lastIndex := 0
	freqFirst = m.options.FrequencyFirst

	for i := 1; i < len(posLenArr); i++ {
		if posLenArr[i].Position >= lastRightBoundary {
			// last is isolated point
			c := i - lastIndex
			arr := make([]dict.PositionLength, c)
			for j := 0; j < c; j++ {
				arr[j] = posLenArr[lastIndex+j]
			}
			leafNodeArray := m.getLeafNodeArrayCore(arr, lastRightBoundary-posLenArr[lastIndex].Position)
			sort.Sort(Nodes(leafNodeArray))
			m.combineNodeAttr(result, leafNodeArray)
			lastIndex = i
		}

		newRightBoundary := posLenArr[i].Position + posLenArr[i].Length
		if newRightBoundary > lastRightBoundary {
			lastRightBoundary = newRightBoundary
		}
	}

	if lastIndex < len(posLenArr) {
		// last is isolated point
		c := len(posLenArr) - lastIndex
		arr := make([]dict.PositionLength, c)
		for j := 0; j < c; j++ {
			arr[j] = posLenArr[lastIndex+j]
		}
		leafNodeArray := m.getLeafNodeArrayCore(arr, lastRightBoundary-posLenArr[lastIndex].Position)
		sort.Sort(Nodes(leafNodeArray))
		m.combineNodeAttr(result, leafNodeArray)
	}

	return result
}

func (m *ChsFullTextMatch) getLeafNodeArrayCore(posLenArr []dict.PositionLength, orginalTextLength int) []*Node {
	m.leafNodeList = [](*Node){}
	m.posLenArr = posLenArr
	m.inputStringLen = orginalTextLength
	m.buildTree(m.root, 0)	
	return m.leafNodeList
}

func (m *ChsFullTextMatch) combineNodeAttr(result []*Node, arr []*Node) {
    c := len(result) - len(arr)
    if c > 0 {
       for i := 0; i < c; i++ {
          arr = append(arr, nil)
       }
    }

	// 复制 arr 链表
	for i := 0; i < len(arr); i++ {
		if i == 0 {
			if arr[i] == nil {
				return
			}
			continue
		}

		if i >= len(result) {
			break
		}

		if arr[i] == nil {
			arr[i] = arr[i-1]
		}

		fst := NewNodeClone(arr[i])
		node := fst
		n := arr[i]

		n = n.Parent
		for j := 1; j < arr[i].AboveCount; j++ {
			node.Parent = NewNodeClone(n)
			node = node.Parent
			n = n.Parent
		}

		arr[i] = fst
	}

	// 如果result 的有效数量少于arr,将result的有效值填充到和arr相等
	// 如果result 没有一个有效值，则不做处理
	for i := 0; i < len(result); i++ {
		if i >= len(arr) {
			break
		}

		if result[i] == nil && arr[i] != nil {
			if i > 0 {
				result[i] = result[i-1]
			}
		}
	}

	for i := 0; i < len(result); i++ {
		j := i
		if len(arr) <= i {
			j = len(arr) - 1
		}
		if arr[j] == nil {
			if result[i] == nil {
				return
			} else {
				for arr[j] == nil {
					j--
				}
			}
		}

		if result[i] == nil {
			// 只有在result 没有一个有效值时才会到这个分支
			result[i] = arr[j]
		} else {
			n := arr[j]
			for k := 0; k < arr[j].AboveCount-1; k++ {
				n = n.Parent
			}

			n.Parent = result[i]
			aboveCount := arr[j].AboveCount + result[i].AboveCount
			result[i] = arr[j]
			result[i].AboveCount = aboveCount
		}
	}
}

func (m *ChsFullTextMatch) mergeAllCombinations(redundancy int) *list.List {
	result := list.New()

	if (redundancy == 0 || !m.options.MultiDimensionality) && !m.options.ForceSingleWord {
		for _, v := range m.allCombinations[0] {
			result.PushBack(v)
		}
		return result
	}

	i := 0
	var cur *list.Element
	forceOnce := false

loop:
	for i <= redundancy && i < len(m.allCombinations) {
		cur = result.Front()
		for j := 0; j < len(m.allCombinations[i]); j++ {
			m.allCombinations[i][j].Level = i
			if cur != nil {
				for cur.Value.(dict.PositionLength).Position < m.allCombinations[i][j].Position {
					cur = cur.Next()
					if cur == nil {
						break
					}
				}
				if cur != nil {
					if cur.Value.(dict.PositionLength).Position != m.allCombinations[i][j].Position || cur.Value.(dict.PositionLength).Length != m.allCombinations[i][j].Length {
						result.InsertBefore(m.allCombinations[i][j], cur)
					}
				} else {
					result.PushBack(m.allCombinations[i][j])
				}
			} else {
				result.PushBack(m.allCombinations[i][j])
			}
		}

		i++
	}

	if m.options.ForceSingleWord && !forceOnce {
		i = len(m.allCombinations) - 1
		redundancy = i
		forceOnce = true
		goto loop
	}

	return result
}

func (m *ChsFullTextMatch) getUnknownWords(masks []int, orginalText []rune) (unknownWords []*dict.WordInfo, needRemoveSingleWord bool) {
	unknownWords = [](*dict.WordInfo){}

	// 找到所有未登录词
	needRemoveSingleWord = false

	j := 0
	begin := false
	beginPosition := 0
	for j < len(masks) {
		if m.options.UnknownWordIdentify {
			if !begin {
				if m.isKnownSingleWord(masks, j, orginalText) {
					begin = true
					beginPosition = j
				}
			} else {
				mergeUnknownWord := true
				if !m.isKnownSingleWord(masks, j, orginalText) {
					if j-beginPosition <= 2 {
						for k := beginPosition; k < j; k++ {
							mergeUnknownWord = false
							if masks[k] != 1 {
								word := string(orginalText[k : k+1])
								wi := dict.NewWordInfoDefault()
								wi.Word = word
								wi.Position = k
								wi.WordType = dict.TNone
								wi.Rank = m.params.UnknowRank
								unknownWords = append(unknownWords, wi)
							}
						}
					} else {
						for k := beginPosition; k < j; k++ {
							if masks[k] == 1 {
								masks[k] = 11
								needRemoveSingleWord = true
							}
						}
					}

					begin = false

					if mergeUnknownWord {
						word := string(orginalText[beginPosition:j])
						wi := dict.NewWordInfoDefault()
						wi.Word = word
						wi.Position = beginPosition
						wi.WordType = dict.TNone
						wi.Rank = m.params.UnknowRank
						unknownWords = append(unknownWords, wi)
					}
				}
			}
		} else {
			if m.isKnownSingleWord(masks, j, orginalText) {
				wi := dict.NewWordInfoDefault()
				wi.Word = string(orginalText[j])
				wi.Position = j
				wi.WordType = dict.TNone
				wi.Rank = m.params.UnknowRank
				unknownWords = append(unknownWords, wi)
			}
		}

		j++
	}

	if begin && m.options.UnknownWordIdentify {
		mergeUnknownWord := true
		if j-beginPosition <= 2 {
			for k := beginPosition; k < j; k++ {
				mergeUnknownWord = false
				if masks[k] != 1 {
					word := string(orginalText[k:(k + 1)])
					wi := dict.NewWordInfoDefault()
					wi.Word = word
					wi.Position = k
					wi.WordType = dict.TNone
					wi.Rank = m.params.UnknowRank
					unknownWords = append(unknownWords, wi)
				}
			}
		} else {
			for k := beginPosition; k < j; k++ {
				if masks[k] == 1 {
					masks[k] = 11
					needRemoveSingleWord = true
				}
			}
		}

		begin = false

		if mergeUnknownWord {
			word := string(orginalText[beginPosition:j])
			wi := dict.NewWordInfoDefault()
			wi.Word = word
			wi.Position = beginPosition
			wi.WordType = dict.TNone
			wi.Rank = m.params.UnknowRank
			unknownWords = append(unknownWords, wi)
		}
	}
	return
}

func (m *ChsFullTextMatch) isKnownSingleWord(masks []int, index int, orginalText []rune) bool {
	state := masks[index]
	if state == 2 {
		return false
	}
	if state == 1 {
		if !m.options.UnknownWordIdentify {
			return false
		}
		// 如果单字是连词/助词/介词/代词
		wa := m.wordDict.GetWordAttr(orginalText[index:(index + 1)])
		if wa != nil {
			if (wa.Pos & SingleWordMask) != 0 {
				return false
			}
		}
	}
	return true
}

type Node struct {
	AboveCount      int
	SpaceCount      int
	FreqSum         float64
	SingleWordCount int
	PosLen          dict.PositionLength
	Parent          *Node
}

type Nodes [](*Node)

func (nodes Nodes) Len() int {
	return len(nodes)
}

func (nodes Nodes) Less(i, j int) bool {
	if nodes[i].SpaceCount < nodes[j].SpaceCount {
		return true
	} else if nodes[i].SpaceCount > nodes[j].SpaceCount {
		return false
	} else {
		if nodes[i].AboveCount < nodes[j].AboveCount {
			return true
		} else if nodes[i].AboveCount > nodes[j].AboveCount {
			return false
		} else {
			if freqFirst {
				if nodes[i].FreqSum > nodes[j].FreqSum {
					return true
				} else if nodes[i].FreqSum < nodes[j].FreqSum {
					return false
				} else {
					if nodes[i].SingleWordCount < nodes[j].SingleWordCount {
						return true
					} else if nodes[i].SingleWordCount > nodes[j].SingleWordCount {
						return false
					}
				}
			} else {
				if nodes[i].SingleWordCount < nodes[j].SingleWordCount {
					return true
				} else if nodes[i].SingleWordCount > nodes[j].SingleWordCount {
					return false
				} else {
					if nodes[i].FreqSum > nodes[j].FreqSum {
						return true
					} else if nodes[i].FreqSum < nodes[j].FreqSum {
						return false
					}
				}
			}
		}
	}
	return true
}

func (nodes Nodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func NewNode() *Node {
	return &Node{AboveCount: 0}
}

func NewNodeClone(node *Node) (newNode *Node) {
	newNode = &Node{}
	newNode.AboveCount = node.AboveCount
	newNode.SpaceCount = node.SpaceCount
	newNode.FreqSum = node.FreqSum
	newNode.SingleWordCount = node.SingleWordCount
	newNode.PosLen = node.PosLen
	newNode.Parent = nil
	return
}

func NewNodeFull(pl dict.PositionLength, parent *Node, aboveCount int, spaceCount int, singleWordCount int, freqSum float64) (node *Node) {
	node = &Node{}
	node.PosLen = pl
	node.Parent = parent
	node.AboveCount = aboveCount
	node.SpaceCount = spaceCount
	node.SingleWordCount = singleWordCount
	node.FreqSum = freqSum
	return
}

func (m *ChsFullTextMatch) buildTree(parent *Node, curIndex int) {
	// 嵌套太多的情况一般很少发生，如果发生，强行中断，
	//以免造成博弈树遍历层次过多降低系统效率
	if len(m.leafNodeList) > 8192 {
		return
	}
 
	if curIndex < len(m.posLenArr)-1 {
		if m.posLenArr[curIndex+1].Position == m.posLenArr[curIndex].Position {
			m.buildTree(parent, curIndex+1)
		}
	}

	spaceCount := parent.SpaceCount + m.posLenArr[curIndex].Position - (parent.PosLen.Position + parent.PosLen.Length)
	singleWordCount := parent.SingleWordCount
	if m.posLenArr[curIndex].Length == 1 {
		singleWordCount += 1
	}

	freqSum := 0.0
	if m.options != nil && m.options.FrequencyFirst {
		freqSum = parent.FreqSum + m.posLenArr[curIndex].WordAttri.Frequency
	}

	curNode := NewNodeFull(m.posLenArr[curIndex], parent, parent.AboveCount+1, spaceCount, singleWordCount, freqSum)
	cur := curIndex + 1
	for cur < len(m.posLenArr) {
		if m.posLenArr[cur].Position >= m.posLenArr[curIndex].Position+m.posLenArr[curIndex].Length {
			m.buildTree(curNode, cur)
			break
		}
		cur++
	}

	if cur >= len(m.posLenArr) {
		curNode.SpaceCount += m.inputStringLen - curNode.PosLen.Position - curNode.PosLen.Length
		m.leafNodeList = append(m.leafNodeList, curNode)
	}
}
