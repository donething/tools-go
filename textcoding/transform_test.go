package textcoding

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestTransformDir(t *testing.T) {
	dir := "D:/Temp/fan"
	// 读取目录
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("[错误]读取目录（%s）时出错：%s\n", dir, err)
		return
	}
	// 遍历文件
	for _, file := range files {
		// 跳过目录
		if file.IsDir() {
			continue
		}
		path := filepath.Join(dir, file.Name())
		err = TransformFile(path)
		if err != nil {
			log.Printf("[失败]转换文件(%s)的编码时出错：%s\n", path, err)
			continue
		}
		log.Printf("[成功]已转换文件(%s)的编码\n", path)
	}
	log.Printf("[完成]转换了指定的路径(%s)下文件的编码\n", dir)
}
