package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"os"
	"strings"
	"tools-go/download"
	"tools-go/file"
	"tools-go/film"
	"tools-go/textcoding"
	"tools-go/tg"
	"tools-go/tg/tg_iface"
)

// 大写为操作，小写为参数
var (
	D bool // 下载 JSON 格式的图集
	F bool // 扫描文件
	G bool // 发送本地文件到 TG
	H bool // 帮助
	S bool // 下载字幕
	T bool // 转换指定目录下指定格式的文件编码

	// 参数
	d string // 目标路径
	f string // 文件的格式
	k string // 关键字
	i string // 数据来源(地址、路径)
	p bool   // 是否全路径
)

// 字幕格式映射，用于快速指定格式
var formatMap = map[string]string{
	"1": ".srt.ass.ssa",
	"2": ".avi.mkv.mov.mp4.rm.rmvb.ts",
	"3": ".aac.ape.flac.mp3.wav.wma",
}

// 解析快速匹配的格式
func parseFormat(f string) string {
	// 为了通过指定数字参数快速指定多种格式，此时先做匹配
	// 又因为前面可以加"?"表示排除，需要特别处理
	fPre := "" // 为可能存在的排除符号"?"
	fReal := f // 为可能的快速指定数字"1"、"2"等
	// 当存在排除符号时，先要提取该符号和数字
	if f != "" && f[0] == '?' {
		fPre = string(f[0])
		fReal = f[1:]
	}

	// 先根据数字快速选择格式
	format, ok := formatMap[fReal]
	// 合并排除符号和格式
	format = fPre + format

	// 如果没有匹配到快速格式，说明不是快速指定的格式，需要恢复为初始输入的格式
	if !ok {
		format = f
	}

	return format
}

// color Tips 使用说明
// color.Notice 完成整个步骤的一小步骤。如 下载了某项、提示输入
// color.Success 完成整个步骤。如 整个任务下载完成
// color.Info 普通提醒。如 提示正在进行的操作
// color.Warn 轻微注意。如 没有该文件

func init() {
	// 操作
	flag.BoolVar(&H, "H", false, "帮助")
	flag.BoolVar(&D, "D", false, "下载 JSON 图集："+"如'-D -d /save/dir -i /json/path'")
	flag.BoolVar(&F, "F", false, "扫描文件：如'-F -f .mp3.mp4 -k 东风破 -p false'。"+
		"-f 文件格式，可多个用'.'分隔；"+
		"-k 搜索关键字，前加'?'符号则为排除，需返回全路径时可指定`p`的值为任意非空值")
	flag.BoolVar(&G, "G", false, "发送本地文件到 TG："+
		"如'-G -k token -f chatid -i sites -d /imgs/dir'。"+
		"-k 为 TG 的 token；"+
		"-f 为TG的频道ID；"+
		"-i 为指定处理函数，暂时可能值为'lj'；"+
		"-d 为待发送的目录")
	flag.BoolVar(&T, "T", false, "转换指定目录下的指定格式的文件编码为 UTF-8："+
		"如'-T -i /path -f .ass.srt'")
	flag.BoolVar(&S, "S", false, "下载电影的字幕："+
		"如'-S [-k 电影名] [-i /film/path]'。"+
		"如果仅指定 -k，则下载字幕到当前所在路径；"+
		"如果还指定 -i，则自动重命名字幕、再保存到相同目录；"+
		"如果仅指定 -i，则将自动从电影路径中提取关键字来搜索")

	// 参数
	flag.StringVar(&d, "d", "", "目标路径")
	flag.StringVar(&f, "f", "", "文件格式，可以用逗号分隔多个格式。"+
		"如'txt,cue,srt'。还可通过数字快速指定某类文件，1：字幕文件格式")
	flag.StringVar(&k, "k", "", "关键字")
	flag.StringVar(&i, "i", "", "数据来源(地址、路径)")
	flag.BoolVar(&p, "p", false, "是否为完整路径")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	// 解析格式
	format := parseFormat(f)

	// 根据参数指定操作
	if H {
		flag.PrintDefaults()
	} else if D {
		color.Info.Tips("开始下载 JSON 格式的图集'%s'，保存到'%s'", i, d)
		download.DLImgs(d, i)
	} else if F {
		color.Info.Tips("开始扫描目录'%s'", d)
		file.Scan(d, format, k, p)
	} else if G {
		if strings.TrimSpace(f) == "" || strings.TrimSpace(d) == "" {
			color.Warn.Tips("无法发送到TG，TG的token或chatid为空\n", i)
			return
		}

		var t tg_iface.ITask
		switch i {
		case "lj":
			t = tg_iface.LJTask{Task: tg_iface.Task{
				ChatID: f,
				Dir:    d,
			}}
		default:
			color.Warn.Tips("无法发送到TG，未知的处理标识'%s'\n", i)
			return
		}

		tg.Init(k)
		tg.SendDir(t)
	} else if S {
		color.Info.Tips("开始尝试下载'%s'的字幕", k)
		film.DLSubtitle(k, i)
	} else if T {
		color.Info.Tips("开始执行转换文本编码。格式：'%s'，路径：'%s'", format, i)
		textcoding.TransformDir(i, format)
	} else {
		usage()
	}
}

func usage() {
	_, err := fmt.Fprintf(os.Stderr, `Usage: tools-go [-DFHST] [-i addr/path] [-f format] [-k keyword]
Options:
`)
	flag.PrintDefaults()

	if err != nil {
		color.Error.Tips("显示帮助说明出错：%s", err)
		os.Exit(0)
	}
}
