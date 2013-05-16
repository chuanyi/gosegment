package dict

import (
	"container/list"
	"segment/utils"
	"strconv"
	"strings"
)

type PositionLength struct {
	Level     int
	Position  int
	Length    int
	WordAttri *WordAttr
}

func NewPositionLength(pos int, len int, word *WordAttr) PositionLength {
	return PositionLength{0, pos, len, word}
}

type WordDictionary struct {
	wordDict       map[string](*WordAttr)
	firstCharDict  map[rune](*WordAttr)
	doubleCharDict map[int32](*WordAttr)
	tripleCharDict map[int64](*[]byte)
	ChineseName    *ChsName
}

func NewWordDictionary() *WordDictionary {
	return &WordDictionary{}
}

func (d *WordDictionary) Load(fileName string) (err error) {
	d.wordDict = make(map[string](*WordAttr))
	d.firstCharDict = make(map[rune](*WordAttr))
	d.doubleCharDict = make(map[int32](*WordAttr))
	d.tripleCharDict = make(map[int64](*[]byte))

	waList, err := d.loadFromTextFile(fileName)
	if err != nil {
		return err
	}

	for e := waList.Front(); e != nil; e = e.Next() {
		key := strings.ToLower(e.Value.(*WordAttr).Word)
		runes := utils.ToRunes(key)

		if len(runes) == 1 {
			d.firstCharDict[runes[0]] = e.Value.(*WordAttr)
			continue
		}

		if len(runes) == 2 {
			doubleChar := runes[0]*65536 + runes[1]
			d.doubleCharDict[doubleChar] = e.Value.(*WordAttr)
			continue
		}

		d.wordDict[key] = e.Value.(*WordAttr)
		tripleChar := int64(int32(runes[0]))*int64(0x100000000) + int64(int32(runes[1]))*int64(65536) + int64(int32(runes[2]))
		var wordLenArray []byte
		v, ok := d.tripleCharDict[tripleChar]
		if !ok {
			wordLenArray = make([]byte, 4)
			wordLenArray[0] = byte(len(runes))
			d.tripleCharDict[tripleChar] = &wordLenArray
		} else {
			find := false
			i := 0
			for i = 0; i < len(*v); i++ {
				if (*v)[i] == byte(len(runes)) {
					find = true
					break
				}
				if (*v)[i] == byte(0) {
					(*v)[i] = byte(len(runes))
					find = true
					break
				}
			}
			if !find {
				var temp []byte = make([]byte, len(*v)*2)
				copy(temp, (*v))
				temp[i] = byte(len(runes))
				d.tripleCharDict[tripleChar] = &temp
			}
		}
	}
	return nil
}

func (d *WordDictionary) loadFromTextFile(fileName string) (dicts *list.List, err error) {
	dicts = list.New()
	err = utils.EachLine(fileName, func(line string) {
		words := strings.Split(string(line), "|")
		if len(words) == 3 {
			word := strings.TrimSpace(words[0])
			pos, _ := strconv.ParseInt(words[1], 0, 0)
			frequency, _ := strconv.ParseFloat(words[2], 64)
			dicts.PushBack(NewWordAttr(word, int(pos), frequency))
		}
	})
	return
}

func (d *WordDictionary) GetWordAttr(word []rune) *WordAttr {
	if len(word) == 1 {
		if wa, ok := d.firstCharDict[word[0]]; ok {
			return wa
		}
	} else if len(word) == 2 {
		doubleChar := word[0]*65536 + word[1]
		if wa, ok := d.doubleCharDict[doubleChar]; ok {
			return wa
		}
	} else {
		if wa, ok := d.wordDict[strings.ToLower(string(word))]; ok {
			return wa
		}
	}
	return nil
}

func (d *WordDictionary) GetAllMatchs(text string, chineseNameIdentify bool) (result []PositionLength) {
	result = []PositionLength{}
	if len(text) == 0 {
		return
	}

	rtext := utils.ToRunes(text)

	keyText := rtext
	if rtext[0] < 128 {
		keyText = utils.ToRunes(strings.ToLower(text))
	}

	for i := 0; i < len(rtext); i++ {
		fst := keyText[i]

		var chsNames []string = nil
		if chineseNameIdentify {
			chsNames = d.ChineseName.Match(rtext, i)
			for _, name := range chsNames {
				wa := NewWordAttr(name, POS_A_NR, 0)
				result = append(result, PositionLength{0, i, utils.RuneLen(name), wa})
			}
		}

		if fwa, ok := d.firstCharDict[fst]; ok {
			result = append(result, PositionLength{0, i, 1, fwa})
		}

		if i < len(keyText)-1 {
			doubleChar := keyText[i]*65536 + keyText[i+1]
			if fwa, ok := d.doubleCharDict[doubleChar]; ok {
				result = append(result, PositionLength{0, i, 2, fwa})
			}
		}

		if i >= len(keyText)-2 {
			continue
		}

		tripleChar := int64(int32(keyText[i]))*0x100000000 + int64(int32(keyText[i+1]))*65536 + int64(int32(keyText[i+2]))
		if lenList, ok := d.tripleCharDict[tripleChar]; ok {
			for _, ilen := range *lenList {
				if ilen == 0 {
					break
				}
				if (i + int(ilen)) > len(keyText) {
					continue
				}
				key := string(keyText[i:(i + int(ilen))])
				if wa, ok := d.wordDict[key]; ok {
					if chsNames != nil {
						find := false
						for _, name := range chsNames {
							if wa.Word == name {
								find = true
								break
							}
						}
						if find {
							continue
						}
					}
					result = append(result, PositionLength{0, i, int(ilen), wa})
				}
			}
		}
	}
	
	return
}
