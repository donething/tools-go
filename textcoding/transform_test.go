package textcoding

import (
	"log"
	"testing"
)

func TestTransformDir(t *testing.T) {
	err := TransformDir(`D:\音乐\无损\华语经典\经典老歌  百万畅销 LP黑胶 10CD`, "cue")
	if err != nil {
		log.Fatalf("转换出错：%v\n", err)
	}
	log.Printf("[完成]转换编码完成\n")
}
