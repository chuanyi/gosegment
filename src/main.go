package main

import (
	"fmt"
	"segment"
	"segment/dict"
)

func main() {
	seg := segment.NewSegment()
	err := seg.Init("./dicts")
	if err != nil {
		fmt.Println("%v", err)
	}
	ret := seg.DoSegment(`盘古分词 简介: 盘古分词 是由eaglet 开发的一款基于字典的中英文分词组件
主要功能: 中英文分词，未登录词识别,多元歧义自动识别,全角字符识别能力
主要性能指标:
分词准确度:90%以上
处理速度: 300-600KBytes/s Core Duo 1.8GHz
用于测试的句子:
长春市长春节致词
长春市长春药店
IＢM的技术和服务都不错
张三在一月份工作会议上说的确实在理
于北京时间5月10日举行运动会
我的和服务必在明天做好`)
	
	for cur := ret.Front(); cur != nil; cur = cur.Next() {
	    w := cur.Value.(*dict.WordInfo)
		fmt.Print(w.Word,"(",w.Position,",",w.Rank,")/")	
	}
}
