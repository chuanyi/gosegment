package dict

import (
	"segment/utils"
	"strings"
)

const SynonymFileName = "Synonym.txt"

type Synonym struct {
	groupList     []([]string) //同义词组，文件中一行为一组，一组同义词以 “,” 分割
	wordToGroupId map[string]([]int)
}

func NewSynonym() *Synonym {
	s := &Synonym{}
	s.groupList = []([]string){}
	s.wordToGroupId = make(map[string]([]int))
	return s
}

func (s *Synonym) Load(dictPath string) (err error) {
	err = utils.EachLine(dictPath+"/"+SynonymFileName, func(line string) {
		if len(line) > 0 {
			words := strings.Split(line, ",")
			s.groupList = append(s.groupList, words)
			groupId := len(s.groupList) - 1
			for i := 0; i < len(words); i++ {
				key := strings.TrimSpace(words[i])
				if l, ok := s.wordToGroupId[key]; ok {
					if l[len(l)-1] == groupId {
						continue
					}
					s.wordToGroupId[key] = append(s.wordToGroupId[key], groupId)
				} else {
					s.wordToGroupId[key] = []int{groupId}
				}
			}
		}
	})
	return
}

func (s *Synonym) GetSynonyms(text string) []string {
	word := strings.ToLower(strings.TrimSpace(text))
	if l, ok := s.wordToGroupId[word]; ok {
		result := []string{}
		for groupId := range l {
			for _, w := range s.groupList[groupId] {
				if w == word {
					continue
				}

				found := false
				for _, wo := range result {
					if w == wo {
						found = true
						break
					}
				}
				if found {
					continue
				}

				result = append(result, w)
			}
		}
		return result
	}
	return nil
}
