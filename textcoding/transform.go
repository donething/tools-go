// Package textcoding 转换文本、文件的编码
package textcoding

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

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
	case "GB-18030":
		decoder = simplifiedchinese.GB18030.NewDecoder()
	case "GBK":
		decoder = simplifiedchinese.GBK.NewDecoder()
	case "GB-2312":
		decoder = simplifiedchinese.HZGB2312.NewDecoder()
	default:
		return nil, fmt.Errorf("还未适配的编码：%s", result.Charset)
	}
	// 转换
	reader := transform.NewReader(bytes.NewReader(bs), decoder)
	// 读取结果
	return ioutil.ReadAll(reader)
}

// TransformFile 转换文本文件的编码为 UTF-8
// 只转换能匹配到的编码类型，参看 TransformText()
func TransformFile(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	data, err := TransformText(bs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}
