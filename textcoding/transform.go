// Package textcoding 转换文本、文件的编码
package textcoding

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

var errUnknownCoding = fmt.Errorf("还未适配的编码")

// TransformText 转换文本的编码为 UTF-8
// 只转换能匹配到的编码类型：GB-18030、GBK、GB-2312
func TransformText(bs []byte) ([]byte, error) {
	// 检测文本的编码
	result, err := DetectTextCoding(bs)
	if err != nil {
		return nil, err
	}

	// 根据文本编码获取对应的编码器
	var decoder *encoding.Decoder
	switch result.Charset {
	case "UTF-8":
		return nil, nil
	case "GB-18030":
		decoder = simplifiedchinese.GB18030.NewDecoder()
	case "GBK":
		decoder = simplifiedchinese.GBK.NewDecoder()
	case "GB-2312":
		decoder = simplifiedchinese.HZGB2312.NewDecoder()
	default:
		return nil, fmt.Errorf("%w'%s'", errUnknownCoding, result.Charset)
	}
	// 转换
	reader := transform.NewReader(bytes.NewReader(bs), decoder)
	// 读取结果
	return ioutil.ReadAll(reader)
}

// TransformFile 转换文本文件的编码为 UTF-8
// 只转换能匹配到的编码类型，参看 TransformText()
func TransformFile(path string) (bool, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	data, err := TransformText(bs)
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil
	}
	return true, ioutil.WriteFile(path, data, 0644)
}

// TransformDir 转换目录下的文件的编码为 UTF-8
//
// dirPath: 目录路径
//
// format: 需转编码的文件格式如(".txt")
func TransformDir(dirPath string, format string) error {
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
		if filepath.Ext(path) != format {
			return nil
		}

		// 转换文件
		has, err := TransformFile(path)

		if errors.Is(err, errUnknownCoding) {
			fmt.Printf("暂时不需要转换编码，文件'%s'：%s\n", path, err)
			return nil
		}
		if err != nil {
			return err
		}

		if has {
			fmt.Printf("已转换文件'%s'\n", path)
		} else {
			// fmt.Printf("不需要转换编码'%s'", path)
		}
		return nil
	})
	return err
}
