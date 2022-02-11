package main

import (
	"flag"
	"fmt"
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

func init() {
	// 操作
	flag.BoolVar(&t, "t", false, "转换指定目录下指定格式的文件编码为 UTF-8，"+
		"如'-t -f .txt -p /to/the/path'")

	// 参数
	flag.StringVar(&f, "f", "", "文件格式，如'.txt'")
	flag.StringVar(&p, "p", ".", "路径")
	flag.BoolVar(&h, "h", false, "帮助")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	var err error
	if h {
		flag.PrintDefaults()
	} else if t {
		fmt.Printf("开始执行转换文本编码。格式：'%s'，路径：'%s'\n", f, p)
		err = textcoding.TransformDir(p, f)
	} else {
		usage()
	}

	// 处理可能的错误
	if err != nil {
		fmt.Printf("执行程序出错：%s\n", err)
	}
	fmt.Printf("已执行完操作\n")
}

func usage() {
	_, err := fmt.Fprintf(os.Stderr, `Usage: tools-go [-t] [-f format] [-p path]
Options:
`)
	flag.PrintDefaults()

	if err != nil {
		fmt.Printf("显示帮助说明出错：%s\n", err)
		os.Exit(0)
	}
}
