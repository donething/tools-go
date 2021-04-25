// 检测文本的编码
package textcoding

import (
	"github.com/saintfish/chardet"
	"io/ioutil"
)

// 返回指定路径的文本文件的编码，返回 编码、地区、准确度（如 GB-18030、zh、100）
func DetectFileCoding(path string) (*chardet.Result, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DetectTextCoding(bs)
}

// 检测文本的编码，返回 编码、地区、准确度（如 GB-18030、zh、100）
// [chardet: Charset detector library for golang derived from ICU](https://github.com/saintfish/chardet)
func DetectTextCoding(data []byte) (result *chardet.Result, err error) {
	detector := chardet.NewTextDetector()
	return detector.DetectBest(data)
}
