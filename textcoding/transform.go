// Package textcoding 转换中文编码的文本、文件为 UTF-8
package textcoding

import (
	"errors"
	"fmt"
	"github.com/donething/utils-go/dotext"
	"github.com/gookit/color"
	"io/fs"
	"path/filepath"
	"strings"
)

var errUnknownCoding = fmt.Errorf("还未适配的编码")

// 转换结果计数
var (
	done int
	skip int
	fail int
)

// TransformDir 转换目录下中文编码的文件为 UTF-8 编码
//
// dirPath: 目录路径
//
// format: 需转编码的文件格式如，可多个("txt,cue,srt")
func TransformDir(dirPath string, format string) {
	// 遍历目录
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		// 读取文件出错
		if err != nil {
			return err
		}
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		// 跳过非指定格式的文件
		// filepath.Ext()返回格式包含点号'.'
		if strings.Index(format, strings.TrimLeft(filepath.Ext(path), ".")) == -1 {
			return nil
		}

		// 转换文件
		has, encoding, err := dotext.File2UTF8(path)

		if errors.Is(err, errUnknownCoding) {
			skip++
			color.Warn.Tips("还未适配的编码'%s'：文件：'%s'\n", encoding, path)
			return nil
		}

		if err != nil {
			fail++
			return err
		}

		if has {
			done++
			color.Notice.Tips("已转换编码'%s'，文件：'%s'\n", encoding, path)
		} else {
			skip++
			color.Debug.Tips("无需转换该编码'%s'，文件：'%s'\n", encoding, path)
		}

		return nil
	})

	if err != nil {
		color.Error.Tips("遍历路径'%s'出错：%s\n", dirPath, err)
		return
	}

	color.Success.Tips("已完成 转换编码：转换 %d 个，跳过 %d 个，失败 %d 个\n", done, skip, fail)
}
