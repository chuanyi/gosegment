package dict

import (
	"segment/utils"
)

const (
	chsSingleNameFileName  = "ChsSingleName.txt"
	chsDoubleName1FileName = "ChsDoubleName1.txt"
	chsDoubleName2FileName = "ChsDoubleName2.txt"
)

var FAMILY_NAMES = []string{
	//有明显歧异的姓氏
	"王", "张", "黄", "周", "徐",
	"胡", "高", "林", "马", "于",
	"程", "傅", "曾", "叶", "余",
	"夏", "钟", "田", "任", "方",
	"石", "熊", "白", "毛", "江",
	"史", "候", "龙", "万", "段",
	"雷", "钱", "汤", "易", "常",
	"武", "赖", "文", "查",
	//没有明显歧异的姓氏
	"赵", "肖", "孙", "李",
	"吴", "郑", "冯", "陈",
	"褚", "卫", "蒋", "沈",
	"韩", "杨", "朱", "秦",
	"尤", "许", "何", "吕",
	"施", "桓", "孔", "曹",
	"严", "华", "金", "魏",
	"陶", "姜", "戚", "谢",
	"邹", "喻", "柏", "窦",
	"苏", "潘", "葛", "奚",
	"范", "彭", "鲁", "韦",
	"昌", "俞", "袁", "酆",
	"鲍", "唐", "费", "廉",
	"岑", "薛", "贺", "倪",
	"滕", "殷", "罗", "毕",
	"郝", "邬", "卞", "康",
	"卜", "顾", "孟", "穆",
	"萧", "尹", "姚", "邵",
	"湛", "汪", "祁", "禹",
	"狄", "贝", "臧", "伏",
	"戴", "宋", "茅", "庞",
	"纪", "舒", "屈", "祝",
	"董", "梁", "杜", "阮",
	"闵", "贾", "娄", "颜",
	"郭", "邱", "骆", "蔡",
	"樊", "凌", "霍", "虞",
	"柯", "昝", "卢", "柯",
	"缪", "宗", "丁", "贲",
	"邓", "郁", "杭", "洪",
	"崔", "龚", "嵇", "邢",
	"滑", "裴", "陆", "荣",
	"荀", "惠", "甄", "芮",
	"羿", "储", "靳", "汲",
	"邴", "糜", "隗", "侯",
	"宓", "蓬", "郗", "仲",
	"栾", "钭", "历", "戎",
	"刘", "詹", "幸", "韶",
	"郜", "黎", "蓟", "溥",
	"蒲", "邰", "鄂", "咸",
	"卓", "蔺", "屠", "乔",
	"郁", "胥", "苍", "莘",
	"翟", "谭", "贡", "劳",
	"冉", "郦", "雍", "璩",
	"桑", "桂", "濮", "扈",
	"冀", "浦", "庄", "晏",
	"瞿", "阎", "慕", "茹",
	"习", "宦", "艾", "容",
	"慎", "戈", "廖", "庾",
	"衡", "耿", "弘", "匡",
	"阙", "殳", "沃", "蔚",
	"夔", "隆", "巩", "聂",
	"晁", "敖", "融", "訾",
	"辛", "阚", "毋", "乜",
	"鞠", "丰", "蒯", "荆",
	"竺", "盍", "单", "欧",
	//复姓必须在单姓后面
	"司马", "上官", "欧阳",
	"夏侯", "诸葛", "闻人",
	"东方", "赫连", "皇甫",
	"尉迟", "公羊", "澹台",
	"公冶", "宗政", "濮阳",
	"淳于", "单于", "太叔",
	"申屠", "公孙", "仲孙",
	"轩辕", "令狐", "徐离",
	"宇文", "长孙", "慕容",
	"司徒", "司空", "万俟"}

type ChsName struct {
	familyNameDict  map[rune]([]rune)
	singleNameDict  map[rune]rune
	doubleName1Dict map[rune]rune
	doubleName2Dict map[rune]rune
}

func NewChsName() *ChsName {
	c := &ChsName{}
	c.familyNameDict = make(map[rune]([]rune))
	c.singleNameDict = make(map[rune]rune)
	c.doubleName1Dict = make(map[rune]rune)
	c.doubleName2Dict = make(map[rune]rune)
	for _, name := range FAMILY_NAMES {
		runes := utils.ToRunes(name)
		if len(runes) == 1 {
			if _, ok := c.familyNameDict[runes[0]]; !ok {
				c.familyNameDict[runes[0]] = nil
			}
		} else {
			if v, ok := c.familyNameDict[runes[0]]; ok {
				if v == nil {
					c.familyNameDict[runes[0]] = []rune{0}
				}
				c.familyNameDict[runes[0]] = append(c.familyNameDict[runes[0]], runes[1])
			} else {
				c.familyNameDict[runes[0]] = []rune{runes[1]}
			}
		}
	}
	return c
}

func (c *ChsName) Load(dictPath string) (err error) {
	if err = c.loadNameDict(dictPath+"/"+chsSingleNameFileName, c.singleNameDict); err == nil {
		if err = c.loadNameDict(dictPath+"/"+chsDoubleName1FileName, c.doubleName1Dict); err == nil {
			err = c.loadNameDict(dictPath+"/"+chsDoubleName2FileName, c.doubleName2Dict)
		}
	}
	return
}

func (c *ChsName) loadNameDict(filePath string, dict map[rune]rune) (err error) {
	err = utils.EachLine(filePath, func(line string) {
		if len(line) > 0 {
			runes := utils.ToRunes(line)
			dict[runes[0]] = runes[0]
		}
	})
	return
}

func (c *ChsName) Match(text []rune, start int) (result []string) {
	result = nil
	cur := start
	slen := len(text)
	if cur > slen-2 {
		return nil
	}

	f1 := text[cur]
	cur++
	f2 := text[cur]

	f2List, ok := c.familyNameDict[f1]
	if !ok {
		return nil
	}

	if f2List != nil {
		find := false
		hasZero := false
		for _, c := range f2List {
			if c == f2 {
				// 复姓
				cur++
				find = true
				break
			} else if c == 0 {
				// 单姓，首字和某个复姓的首字相同
				hasZero = true
			}
		}
		if !find && !hasZero {
			return nil
		}
	}

	if cur >= slen {
		return nil
	}

	name1 := text[cur]

	if _, ok := c.singleNameDict[name1]; ok {
		result = []string{string(text[start:(cur + 1)])}
	}

	if _, ok := c.doubleName1Dict[name1]; ok {
		cur++
		if cur >= slen {
			return result
		}

		name2 := text[cur]
		if _, ok = c.doubleName2Dict[name2]; ok {
			if result == nil {
				result = []string{}
			}
			result = append(result, string(text[start:(cur+1)]))
		}
	}
	return result
}
