package download

import (
	"encoding/json"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dohttp"
	"github.com/gookit/color"
	"os"
	"path/filepath"
)

var client = dohttp.New(60, false, false)

// DLImgs 下载 JSON 格式的图集
//
// JSON 格式如：{"图集名1": {"1.jpg": "https://a.com/1.jpg"}}
func DLImgs(rootDir string, jsonPath string) {
	var data map[string]map[string]string
	// 读取文件
	bs, err := dofile.Read(jsonPath)
	if err != nil {
		color.Error.Tips("读取 JSON 文件出错：%s", err)
		return
	}

	// 解析
	err = json.Unmarshal(bs, &data)
	if err != nil {
		color.Error.Tips("解析 JSON 文本出错：%s", err)
		return
	}

	// 遍历下载图集
	for title, albums := range data {
		// 创建目录
		dstPath := filepath.Join(rootDir, dofile.ValidFileName(title, "_"))
		color.Notice.Tips("开始下载'%s'", title)
		err = os.MkdirAll(dstPath, 0755)
		if err != nil {
			color.Error.Tips("创建目标路径'%s'出错，无法下载'%s'：%s", dstPath, title, err)
			continue
		}

		// 下载、保存到本地
		for name, url := range albums {
			imgPath := filepath.Join(dstPath, dofile.ValidFileName(name, "_"))
			_, errDl := client.Download(url, imgPath, true, nil)
			if errDl != nil {
				color.Error.Tips("下载图片'%s'出错：%s", imgPath, err)
				// 图集中只要有一个图片下载失败，就认为整个图集下载失败，直接开始下一个图集
				break
			}
		}
		color.Notice.Tips("图集'%s'下载完成", title)
	}
	color.Success.Tips("已完成 下载 JSON 格式的图集")
}
