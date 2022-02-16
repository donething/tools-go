package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"os"
	"tools-go/textcoding"
)

var (
	// 操作
	t bool // 转换指定目录下指定格式的文件编码

	// 参数
	f string // 格式
	p string // 路径
	h bool   // 帮助
)
var (
	// 字幕格式映射，用于快速指定格式
	fortmatMap = map[string]string{"1": "srt,ass,ssa"}
)

func init() {
	// 操作
	flag.BoolVar(&t, "t", false, "转换指定目录下指定格式的文件编码为 UTF-8，"+
		"如'-t -f .txt -p /to/the/path'")

	// 参数
	flag.StringVar(&f, "f", "",
		"文件格式，可以用逗号分隔多个格式，如'txt,cue,srt'。可通过数字快速指定文件类型，1：字幕格式")
	flag.StringVar(&p, "p", ".", "路径")
	flag.BoolVar(&h, "h", false, "帮助")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	// 由于可以指定数字参数快速指定多种格式，此处需要根据数字参数获取快速指定的格式
	format, ok := fortmatMap[f]
	if !ok {
		format = f
	}

	// 执行结果
	var result string
	var err error
	// 根据参数指定操作
	if h {
		flag.PrintDefaults()
	} else if t {
		color.Notice.Tips("开始执行转换文本编码。格式：'%s'，路径：'%s'\n", format, p)
		result, err = textcoding.TransformDir(p, format)
	} else {
		usage()
	}

	// 处理可能的错误
	if err != nil {
		color.Error.Tips("执行程序出错：%s\n", err)
	}

	color.Notice.Tips("%s\n", result)
}

func usage() {
	_, err := fmt.Fprintf(os.Stderr, `Usage: tools-go [-t] [-f format] [-p path]
Options:
`)
	flag.PrintDefaults()

	if err != nil {
		color.Error.Tips("显示帮助说明出错：%s\n", err)
		os.Exit(0)
	}
}
