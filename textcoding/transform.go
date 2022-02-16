// Package textcoding 转换中文编码的文本、文件为 UTF-8
package textcoding

import (
	"errors"
	"fmt"
	"github.com/donething/utils-go/dotext"
	"github.com/gookit/color"
	"io/fs"
	"io/ioutil"
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

// TransformFile 转换中文编码的文件为 UTF-8 编码
//
// 只转换能匹配到的编码类型，参看 TransformText()
//
// 返回是否改变了源文件、原编码、可能的错误
func TransformFile(path string) (bool, string, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return false, "", err
	}
	data, encoding, err := dotext.TransformText(bs)
	if err != nil {
		return false, encoding, err
	}

	if data == nil {
		return false, encoding, nil
	}
	return true, encoding, ioutil.WriteFile(path, data, 0644)
}

// TransformDir 转换目录下中文编码的文件为 UTF-8 编码
//
// dirPath: 目录路径
//
// format: 需转编码的文件格式如，可多个("txt,cue,srt")
func TransformDir(dirPath string, format string) (string, error) {
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
		has, encoding, err := TransformFile(path)

		if errors.Is(err, errUnknownCoding) {
			skip++
			color.Info.Tips("还未适配的编码'%s'：文件：'%s'\n", encoding, path)
			return nil
		}

		if err != nil {
			fail++
			return err
		}

		if has {
			done++
			color.Success.Tips("已转换编码'%s'，文件：'%s'\n", encoding, path)
		} else {
			skip++
			color.Debug.Tips("无需转换该编码'%s'，文件：'%s'\n", encoding, path)
		}

		return nil
	})
	return fmt.Sprintf("已完成转换编码：转换 %d 个，跳过 %d 个，失败 %d 个", done, skip, fail), err
}
