package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"os"
	"tools-go/film"
	"tools-go/textcoding"
)

var (
	// 操作
	h bool // 帮助
	s bool // 下载字幕
	t bool // 转换指定目录下指定格式的文件编码

	// 参数
	a string // 地址/路径
	f string // 文件的格式
	k string // 关键字
)
var (
	// 字幕格式映射，用于快速指定格式
	fortmatMap = map[string]string{"1": "srt,ass,ssa"}
)

func init() {
	// 操作
	flag.BoolVar(&h, "h", false, "帮助")
	flag.BoolVar(&t, "t", false, "转换指定目录下指定格式的文件编码为 UTF-8，"+
		"如'-t -f .txt -a /to/the/path'")
	flag.BoolVar(&s, "s", false, "下载电影的字幕字幕，"+
		"如'-s -k 电影名 [-a /path/the/film]'。如果指定电影路径，则自动重命名字幕文件并保存到相同路径")

	// 参数
	flag.StringVar(&a, "a", "", "地址/路径")
	flag.StringVar(&f, "f", "",
		"文件格式，可以用逗号分隔多个格式，如'txt,cue,srt'。可通过数字快速指定文件类型，1：字幕格式")
	flag.StringVar(&k, "k", "", "关键字")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	// 由于可以指定数字参数快速指定多种格式，此处需要根据数字参数获取快速指定的格式
	format, ok := fortmatMap[f]
	if !ok {
		format = f
	}

	// 根据参数指定操作
	if h {
		flag.PrintDefaults()
	} else if t {
		color.Notice.Tips("开始执行转换文本编码。格式：'%s'，路径：'%s'\n", format, a)
		textcoding.TransformDir(a, format)
	} else if s {
		color.Notice.Tips("开始尝试下载'%s'的字幕\n", k)
		film.DLSubtitle(k, a)
	} else {
		usage()
	}
}

func usage() {
	_, err := fmt.Fprintf(os.Stderr, `Usage: tools-go [-t] [-f format] [-a addr/path]
Options:
`)
	flag.PrintDefaults()

	if err != nil {
		color.Error.Tips("显示帮助说明出错：%s\n", err)
		os.Exit(0)
	}
}
