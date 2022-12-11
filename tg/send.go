package tg

import (
	"fmt"
	"github.com/donething/utils-go/dofile"
	"github.com/donething/utils-go/dotgpush"
	"github.com/gookit/color"
	"math"
	"os"
	"path/filepath"
	"strings"
	"tools-go/tg/tg_iface"
)

// 每个图集可以有的图片数量
const mediaMax = 9

var tg *dotgpush.TGBot

func Init(token string) {
	if tg == nil {
		tg = dotgpush.NewTGBot(token)
		err := tg.SetProxy("socks5://127.0.0.1:1080")
		if err != nil {
			panic(err)
		}
	}
}

// SendDir 发送目录
func SendDir(task tg_iface.ITask) {
	// 读取文件
	files, err := os.ReadDir(task.GetDir())
	if err != nil {
		color.Error.Tips("发送目录出错，无法读取目录'%s'：%s\n", task.GetDir(), err)
		os.Exit(0)
	}

	// 排序
	task.Sort(files)

	for _, file := range files {
		if file.IsDir() {
			dst := filepath.Join(task.GetDir(), file.Name())
			color.Primary.Tips("发送目录'%s'\n", dst)

			// 演员名列表。如"#abc  #xyz"
			actorStr, _ := task.GetActor(file.Name())

			err := sendImgDir(task.GetChatID(), dst, actorStr)
			if err != nil {
				color.Error.Tips("发送目录出错：%s\n", err)
			}
		}
	}
}

// 发送指定目录内的图片
func sendImgDir(chatID string, dir string, actorStr string) error {
	// 图集内的图片等于指定数量后，需要发送
	var count = 0
	// 总part数
	var total = 0
	// 已发送第几个part
	var part = 1
	// 容器
	photos := make([]dotgpush.Media, 0, mediaMax)
	// 目录名
	dirName := filepath.Base(dir)

	// 读取文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	total = int(math.Ceil(float64(len(files)) / float64(mediaMax)))

	for _, file := range files {
		bs, err := dofile.Read(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		count++
		photo := dotgpush.Media{
			Type:  dotgpush.Photo,
			Media: bs,
		}
		// 设置图集的标题（仅设置第一个元素的`Caption`时有效）
		if count%mediaMax == 1 {
			// 提取名字
			name := file.Name()
			name = name[:strings.LastIndex(name, ".")]

			i := strings.Index(name, "-")
			if i != -1 {
				name = name[:i]
			}

			caption := fmt.Sprintf("%s[%s  Part%d]", actorStr, dirName, part)
			// color.Primary.Tips("设置图集'%s'的标题：'%s'。此时总图片数为：%d\n", dirName, caption, count)
			photo.Caption = caption
		}
		photos = append(photos, photo)

		// 图集满photosMax个后需要发送
		if count%mediaMax == 0 {
			color.Notice.Tips("开始发送图集'%s'的第 %d/%d 部分\n", dirName, part, total)
			err := execSendMedia(chatID, photos)
			if err != nil {
				return err
			}

			part++
			photos = photos[:0]
		}
	}

	color.Notice.Tips("开始发送图集'%s'的第 %d/%d 部分\n", dirName, part, total)
	err = execSendMedia(chatID, photos)
	if err != nil {
		return err
	}

	color.Success.Tips("已发送图集'%s'，共发送 %d 张图片\n", dirName, count)
	return nil
}

// 执行发送图片
func execSendMedia(chatID string, photos []dotgpush.Media) error {
	msg, err := tg.SendMediaGroup(chatID, photos)
	if err != nil {
		color.Error.Tips("网络出错，将重试：%s\n", err)
		return execSendMedia(chatID, photos)
	}

	if !msg.Ok {
		return fmt.Errorf("发送图片出错：ErrCode: %d, Msg: %s", msg.ErrorCode, msg.Description)
	}
	return nil
}
