package utils

import (
	"bufio"
	"container/list"
	"errors"
	"io"
	"os"
	"unicode/utf8"
)

// first rune in string
func FirstRune(s string) (r rune) {
	r, _ = utf8.DecodeRuneInString(s)
	return
}

// rune count in string
func RuneLen(s string) int {
	return utf8.RuneCountInString(s)
}

// string to []rune
func ToRunes(s string) (a []rune) {
	a = make([]rune, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		a[i] = r
		i++
	}
	return
}

// read text file line by line
func EachLine(file string, handle func(string)) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	bf := bufio.NewReader(f)
	for {
		line, isPrefix, err := bf.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if isPrefix {
			return errors.New("Error: unexcepted long line.")
		}
		handle(string(line))
	}

	return nil
}

func InsertAfterList(rl *list.List, ol *list.List, mark *list.Element) *list.Element {
	rcur := mark
	for cur := ol.Front(); cur != nil; cur = cur.Next() {
		rcur = rl.InsertAfter(cur.Value, rcur)
	}
	return rcur
}

func IntMin(a int, b int) int {
	if a > b {
		return b
	}
	return a
}
