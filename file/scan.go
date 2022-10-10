package file

import (
	"fmt"
	"github.com/gookit/color"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Not 排除的标志字符
const Not = '?'

// ScanFiles 扫描文件，返回指定的文件名列表
//
// 过滤优先级是`formats`先于`keyword`，即`formats`不满足就不包含了
//
// formats 指定格式，可空、可含多个，如".mp3.mp4.mkv"。最前面可加"?"，表示为排除包含指定格式的文件
//
// keyword 需要包含的关键字，不区分大小写，可空。最前面可加"?"，表示为排除包含指定关键字的文件
//
// fullpath 文件名是否为完整的路径，为`false`表示仅需文件名
func ScanFiles(dir string, formats string, keyword string, fullpath bool) ([]string, error) {
	// 容器
	paths := make([]string, 0, 10)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 不直接处理目录
		if info.IsDir() {
			return nil
		}

		// 小写的过滤格式、过滤关键字、文件名，以不区分大小写过滤
		lf := strings.ToLower(formats)
		lk := strings.ToLower(keyword)
		ln := strings.ToLower(info.Name())

		// 传递了需过滤的格式时，跳过不是指定的格式的文件
		if lf != "" && ((lf[0] == Not && strings.Contains(lf, filepath.Ext(ln))) ||
			(lf[0] != Not && !strings.Contains(lf, filepath.Ext(ln)))) {
			return nil
		}

		// 传递了关键字时，跳过不含关键字的文件
		if lk != "" && ((lk[0] == Not && strings.Contains(ln, lk[1:])) ||
			(lk[0] != Not && !strings.Contains(ln, lk))) {
			return nil
		}

		// 为需要返回的文件名
		name := info.Name()
		// 当需要返回完整路径时
		if fullpath {
			p, errAbs := filepath.Abs(path)
			if errAbs != nil {
				return errAbs
			}
			name = p
		}

		// 添加
		paths = append(paths, name)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

// Scan 扫描目录，参数参考 ScanFiles
func Scan(dir string, formats string, keyword string, fullpath bool) {
	color.Info.Tips("收到扫描目录的参数，目录'%s'，格式'%s'，关键字'%s'，是否全路径'%t'",
		dir, formats, keyword, fullpath)

	list, err := ScanFiles(dir, formats, keyword, fullpath)
	if err != nil {
		color.Error.Tips("扫描目录'%s'出错：%s", dir, err)
		return
	}

	// 合并为字符串
	str := strings.Join(list, "\n\n")
	color.Info.Tips("已扫描目录'%s'，共 %d 个文件", dir, len(list))

	// 小文本直接打印到终端
	if len(str) <= 128 {
		color.Info.Tips("文件列表：")
		fmt.Println(str)
		return
	}

	// 大文本保存到当前目录
	err = os.WriteFile(fmt.Sprintf("文件列表_%d.txt", time.Now().UnixMilli()), []byte(str), 0644)
	if err != nil {
		color.Error.Tips("保存文件列表到本地文件时出错：%s", err)
		return
	}
	color.Info.Tips("已保存文件列表到当前目录")
}
