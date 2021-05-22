package main

import (
	"fmt"
	"strings"
	"tools-go/textcoding"
)

func main() {
	// 输入的参数
	var args string
	fmt.Printf("参数说明：\n1. t[ransform] [suffix] [dir]: 转换目录下指定文本的编码为 UTF-8\n" +
		"2. r[ename] [dir]: 重命名目录下的[特俗]文件\n" +
		"备注：参数 dir 可空，此时为当前目录")
	_, err := fmt.Scanf("%s\n", &args)
	if err != nil {
		fmt.Printf("输入错误：%s\n", err)
		return
	}
	// 根据输入的参数执行操作
	choice := strings.Split(args, " ")
	switch choice[0] {
	case "t", "transform":
		// 输入的参数：后缀、目录
		suffix := ""
		dir := "."
		if len(choice) == 3 {
			suffix = choice[2]
		}
		if len(choice) == 4 {
			dir = choice[3]
		}
		err = textcoding.TransformDir(dir, suffix)
	}
	// 处理可能的错误
	if err != nil {
		fmt.Printf("执行出错：%s\n", err)
	}
}
