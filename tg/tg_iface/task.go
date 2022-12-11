package tg_iface

import (
	"os"
)

type Task struct {
	// TG的聊天频道ID
	ChatID string

	// 目录路径
	Dir string
}

type ITask interface {
	// GetChatID 获取TG的聊天频道ID
	GetChatID() string

	// GetDir 获取目录路径
	GetDir() string

	// GetActor 获取演员列表。 如"abc-xyz-1-high"，返回 "#abc #xyz  "、["abc","xyz"]
	//
	// actorStr 后面有两个空格
	GetActor(name string) (string, []string)

	// GetNameIndex 从名字中提取序号，以便排序。如输入"abc-1-xyz"，返回 1
	GetNameIndex(name string) int

	// Sort 根据文件夹名排序。如"abc-10-xyz"、"abc-2-xyz"，后者应该排在前面
	Sort(files []os.DirEntry)
}
