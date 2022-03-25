package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"os"
	"tools-go/downloads"
	"tools-go/film"
	"tools-go/textcoding"
)

// 大写为操作，小写为参数
var (
	D bool // 下载 JSON 格式的图集
	H bool // 帮助
	S bool // 下载字幕
	T bool // 转换指定目录下指定格式的文件编码

	// 参数
	a string // 地址/路径
	d string // 目标路径
	f string // 文件的格式
	k string // 关键字
)
var (
	// 字幕格式映射，用于快速指定格式
	fortmatMap = map[string]string{"1": "srt,ass,ssa"}
)

// color Tips 使用说明
// color.Notice 完成整个步骤的一小步骤。如 下载了某项、提示输入
// color.Success 完成整个步骤。如 整个任务下载完成
// color.Info 普通提醒。如 提示正在进行的操作
// color.Warn 轻微注意。如 没有该文件

func init() {
	// 操作
	flag.BoolVar(&H, "H", false, "帮助")
	flag.BoolVar(&D, "D", false, "下载 JSON 图集，"+
		"如'-D -d /save/dir -a /json/path'")
	flag.BoolVar(&T, "T", false, "转换指定目录下指定格式的文件编码为 UTF-8，"+
		"如'-T -f txt -a /to/the/path'")
	flag.BoolVar(&S, "S", false, "下载电影的字幕字幕，"+
		"如'-S -k 电影名 [-a /path/the/film]'。如果指定电影路径，则自动重命名字幕文件并保存到相同路径")

	// 参数
	flag.StringVar(&a, "a", "", "源地址/路径")
	flag.StringVar(&d, "d", "", "目标路径")
	flag.StringVar(&f, "f", "", "文件格式，可以用逗号分隔多个格式。"+
		"如'txt,cue,srt'。还可通过数字快速指定某类文件，1：字幕文件格式")
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
	if H {
		flag.PrintDefaults()
	} else if D {
		color.Info.Tips("开始下载 JSON 格式的图集'%s'，保存到'%s'\n", a, d)
		downloads.DLImgs(d, a)
	} else if T {
		color.Info.Tips("开始执行转换文本编码。格式：'%s'，路径：'%s'\n", format, a)
		textcoding.TransformDir(a, format)
	} else if S {
		color.Info.Tips("开始尝试下载'%s'的字幕\n", k)
		film.DLSubtitle(k, a)
	} else {
		usage()
	}
}

func usage() {
	_, err := fmt.Fprintf(os.Stderr, `Usage: tools-go [-T] [-f format] [-a addr/path]
Options:
`)
	flag.PrintDefaults()

	if err != nil {
		color.Error.Tips("显示帮助说明出错：%s\n", err)
		os.Exit(0)
	}
}
