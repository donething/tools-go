package tg_iface

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// LJTask legsjapan
type LJTask struct {
	Task
}

func (t LJTask) GetChatID() string {
	return t.ChatID
}
func (t LJTask) GetDir() string {
	return t.Dir
}
func (t LJTask) GetNameIndex(name string) int {
	if strings.Index(name, "high") == -1 {
		return -1
	}
	reg := regexp.MustCompile("(\\d)+")
	iStr := reg.FindString(name)

	i, err := strconv.Atoi(iStr)
	if err != nil {
		i = -1
	}

	return i
}

func (t LJTask) GetActor(name string) (string, []string) {
	reg := regexp.MustCompile("[^-]([a-zA-Z]+)-")
	actor := reg.FindAllString(name, -1)

	actorStr := ""
	for _, str := range actor {
		// 末尾会多个"-"，需要去除
		actorStr += fmt.Sprintf("#%s  ", strings.Replace(str, "-", "", -1))
	}
	return strings.TrimSpace(actorStr) + "  ", actor
}

func (t LJTask) Sort(files []os.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		// 演员人数少的排前面
		_, actori := t.GetActor(files[i].Name())
		_, actorj := t.GetActor(files[j].Name())
		if len(actori) != len(actorj) {
			return len(actori) < len(actorj)
		}

		// 再按序号排序
		return t.GetNameIndex(files[i].Name()) < t.GetNameIndex(files[j].Name())
	})
}
