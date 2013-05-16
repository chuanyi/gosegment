/**
 * func:  stop word match
 * status: ok
 * 
 */

package dict

import (
	"segment/utils"
	"strings"
)

type StopWord struct {
	stopWordTbl map[string]bool
}

func NewStopWord() (s *StopWord) {
	s = &StopWord{}
	s.stopWordTbl = make(map[string]bool)
	return
}

func (s *StopWord) Load(file string) (err error) {
	err = utils.EachLine(file, func(line string) {
		if len(line) > 0 {
			if utils.FirstRune(line) < 128 {
				s.stopWordTbl[strings.ToLower(line)] = true
			} else {
				s.stopWordTbl[line] = true
			}
		}
	})
	return
}

func (s *StopWord) IsStopWord(word string, filterEnglish bool, filterEnglishLength int, filterNumeric bool, filterNumbericLength int) bool {
	if len(word) == 0 {
		return false
	}

	r := utils.FirstRune(word)
	if r < 128 {
		slen := utils.RuneLen(word)
		if filterEnglish {
			if slen > filterEnglishLength && (r < '0' || r > '9') {
				return true
			}
		}
		if filterNumeric {
			if slen > filterNumbericLength && (r >= '0' && r <= '9') {
				return true
			}
		}
		return s.stopWordTbl[strings.ToLower(word)]
	}

	return s.stopWordTbl[word]
}
