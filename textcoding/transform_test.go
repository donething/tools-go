package textcoding

import (
	"testing"
)

func TestTransformDir(t *testing.T) {
	TransformDir(`D:\音乐\无损\华语经典\经典老歌  百万畅销 LP黑胶 10CD`, "cue")
	t.Logf("[完成]转换编码完成\n")
}
