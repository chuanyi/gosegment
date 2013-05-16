package framework

import (
	"segment/utils"
)

const INC = 50

type Stemmer struct {
    b []rune
    i, i_end, j, k int
}

func NewStemmer() (s *Stemmer) {
    s = &Stemmer{}
    s.b = make([]rune, INC)
    s.i = 0
    s.i_end = 0
    return 
}

func (s *Stemmer) Add(ch rune) {
    if s.i == len(s.b) {
        new_b := make([]rune, (s.i+INC))
        for c := 0; c < s.i; c++ {
            new_b[c] = s.b[c]
        }
        s.b = new_b
    }
    s.b[s.i] = ch
    s.i++
}

func (s *Stemmer) addArr(w []rune, wlen int) {
    if s.i + wlen >= len(s.b) {
    	new_b := make([]rune, (s.i + wlen + INC))
    	for c := 0; c < s.i; c++ {
    	    new_b[c] = s.b[c]
    	}
    	s.b = new_b
    }
    for c := 0; c < wlen; c++ {
        s.b[s.i] = w[c]
        s.i++
    }
}

func (s *Stemmer) ToString() string {
    return string(s.b[0:s.i_end])
}

func (s *Stemmer) getResultLength() int {
    return s.i_end
}

func (s *Stemmer) getResultBuffer() []rune {
    return s.b
}

func (s *Stemmer) cons(i int) bool {
    switch s.b[i] {
    	case 'a','e','i','o','u':
    	    return false
    	case 'y':
    	    if i == 0 {
    	        return true
    	    } else {
    	        return !s.cons(i-1)
    	    }
    }
    return true
}

func (s *Stemmer) m() int {
    n := 0
    i := 0
    for {
        if i > s.j {
            return n
        }
        if !s.cons(i) {
            break
        }
        i++
    }
    i++
    for {
        for {
            if i > s.j {
                return n
            }
            if s.cons(i) {
                break
            }
            i++
        }
        i++
        n++
        for {
            if i > s.j {
                return n
            }
            if !s.cons(i) {
                break
            }
            i++
        }
        i++
    }
    return n
}

func (s *Stemmer) vowelinstem() bool {
    var i int
    for i = 0; i < s.j; i++ {
        if !s.cons(i) {
            return true
        }
    }
    return false
}

func (s *Stemmer) doublec(j int) bool {
    if j < 1 {
        return false
    }
    if s.b[j] != s.b[j-1] {
        return false
    }
    return s.cons(j)
}

func (s *Stemmer) cvc(i int) bool {
    if i < 2 || !s.cons(i) || s.cons(i-1) || !s.cons(i-2) {
        return false
    }
    ch := s.b[i]
    if ch == 'w' || ch == 'x' || ch == 'y' {
        return false
    }
    return true
}

func (s *Stemmer) ends(str string) bool {
    l := utils.RuneLen(str)
    o := s.k - l + 1
    if o < 0 {
        return false
    }
    sc := utils.ToRunes(str)
    for i := 0; i < l; i++ {
        if s.b[o+i] != sc[i] {
            return false
        }
    }
    s.j = s.k - l
    return true
}

func (s *Stemmer) setto(str string) {
    l := utils.RuneLen(str)
    o := s.j + 1
    sc := utils.ToRunes(str)
    for i := 0; i < l; i++ {
        s.b[o+i] = sc[i]
    }
    s.k = s.j + l
}

func (s *Stemmer) r(str string) {
    if s.m() > 0 {
        s.setto(str)
    }
}

func (s *Stemmer) step1() {
    if s.b[s.k] == 's' {
        if s.ends("sses") {
            s.k -= 2
        } else if s.ends("ies") {
            s.setto("i")
        } else if s.b[s.k-1] != 's' {
            s.k--
        }
    }
    if s.ends("eed") {
        if s.m() > 0 {
            s.k--
        }
    } else if (s.ends("ed") || s.ends("ing")) && s.vowelinstem() {
        s.k = s.j
        if s.ends("at") {
            s.setto("ate")
        } else if s.ends("bl") {
            s.setto("ble")
        } else if s.ends("iz") {
            s.setto("ize")
        } else if s.doublec(s.k) {
            s.k--
            ch := s.b[s.k]
            if ch == 'l' || ch == 's' || ch == 'z' {
                s.k++
            }
        } else if s.m() == 1 && s.cvc(s.k) {
            s.setto("e")
        }
    }
}

func (s *Stemmer) step2() {
    if s.ends("y") && s.vowelinstem() {
        s.b[s.k] = 'i'
    }
}

func (s *Stemmer) step3() {
    if s.k == 0 {
        return
    }
    
    switch s.b[s.k-1] {
        case 'a':
            if s.ends("ational") { s.r("ate"); break }
            if s.ends("tional") { s.r("tion"); break }
        case 'c':
            if s.ends("enci") { s.r("ence"); break }
            if s.ends("anci") { s.r("ance"); break }        
        case 'e':
            if s.ends("izer") { s.r("ize"); break }
        case 'l':
            if s.ends("bli") { s.r("ble"); break }
            if s.ends("alli") { s.r("al"); break }
            if s.ends("entli") { s.r("ent"); break }
            if s.ends("eli") { s.r("e"); break }
            if s.ends("ousli") { s.r("ous"); break }
        case 'o':
            if s.ends("ization") { s.r("ize"); break }
            if s.ends("ation") { s.r("ate"); break }
            if s.ends("ator") { s.r("ate"); break }
        case 's':
            if s.ends("alism") { s.r("al"); break }
            if s.ends("iveness") { s.r("ive"); break }
            if s.ends("fulness") { s.r("ful"); break }
            if s.ends("ousness") { s.r("ous"); break }
        case 't':
            if s.ends("aliti") { s.r("al"); break }
            if s.ends("iviti") { s.r("ive"); break }
            if s.ends("biliti") { s.r("ble"); break }
        case 'g':
            if s.ends("logi") { s.r("log"); break }
    }
}

func (s *Stemmer) step4() {
	switch s.b[s.k] {
	    case 'e':
            if s.ends("icate") { s.r("ic"); break }
            if s.ends("ative") { s.r(""); break }
            if s.ends("alize") { s.r("al"); break }
	    case 'i':
	        if s.ends("iciti") { s.r("ic"); break }
	    case 'l':
	        if s.ends("ical") { s.r("ic"); break }
            if s.ends("ful") { s.r(""); break }
	    case 's':
	        if s.ends("ness") { s.r(""); break }
	}
}

func (s *Stemmer) step5() {
    if s.k == 0 {
        return
    }
    
    switch s.b[s.k - 1] {
        case 'a':
        	if s.ends("al") { break }
        	return
        case 'c':
            if s.ends("ance") { break }
            if s.ends("ence") { break }
            return
        case 'e':
        	if s.ends("er") { break }
        	return
        case 'i':
        	if s.ends("ic") { break }
        	return
        case 'l':
            if s.ends("able") { break }
            if s.ends("ible") { break }
            return
        case 'n':
            if s.ends("ant") { break }
            if s.ends("ement") { break }
            if s.ends("ment") { break }
            if s.ends("ent") { break }
            return
        case 'o':
            if s.ends("ion") && s.j >= 0 && (s.b[s.j] == 's' || s.b[s.j] == 't') { break }
            if s.ends("ou") { break }
            return
        case 's':
        	if s.ends("ism") { break }
        	return
        case 't':
            if s.ends("ate") { break }
            if s.ends("iti") { break }
            return
        case 'u':
        	if s.ends("ous") { break }
        	return
        case 'v':
        	if s.ends("ive") { break }
        	return
        case 'z':
         	if s.ends("ize") { break }
        	return
        default:
            return
    }
    if s.m() > 1 {
       s.k = s.j
    }
}

func (s *Stemmer) step6() {
    s.j = s.k
    if s.b[s.k] == 'e' {
        a := s.m()
        if a > 1 || a == 1 && !s.cvc(s.k-1) {
            s.k--
        }
    }
    if s.b[s.k] == 'l' && s.doublec(s.k) && s.m() > 1 {
        s.k--
    }
}

func (s *Stemmer) Stem() {
    s.k = s.i-1
    if s.k > 1 {
        s.step1()
        s.step2()
        s.step3()
        s.step4()
        s.step5()
        s.step6()
    }
    s.i_end = s.k + 1
    s.i = 0
}







