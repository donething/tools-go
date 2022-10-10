# tools-go

常用工具集

## 已有功能

### 扫描目录

```shell
# 扫描目录 D:\影视
tools-go.exe -F -d D:\影视
# 只需要指定格式 .mp4.mp3
tools-go.exe -F -d D:\影视 -f .mp4.mp3
# 排除格式 .mp3、.mp4
tools-go.exe -F -d D:\影视 -f ?.mp4.mp3
# 需要文件名的完整路径
tools-go.exe -F -d D:\影视 -f ?.mp4.mp3 -p 1
# 快速选择格式，1：字幕；2：视频；3：音频
tools-go.exe -F -d D:\影视 -f 1
# 排除字幕格式
tools-go.exe -F -d D:\影视 -f ?1
```

### 下载字幕

```shell
# 指定关键字搜索字幕 黑猫警长
tools-go.exe -S -k 黑猫警长
# 将搜索黑猫警长，并重命名字幕文件为-i后面的电影名（不会包含其格式的电影名）
tools-go.exe -S -k 黑猫警长 -i ./黑猫警长.mp4
# 指定文件时，将搜索提取名 黑猫警长
tools-go.exe -S -i ./黑猫警长.mp4
```

### 转换编码

```shell
# 转换指定文本文件的编码为 UTF-8
tools-go.exe -T -i D:\影视\字幕.srt
# 转换目录下所有文本文件的编码为 UTF-8
tools-go.exe -T -i D:\影视
# 转换目录下指定格式的文本文件的编码为 UTF-8
tools-go.exe -T -i D:\影视 -f .ass.srt
# 转换除了指定格式的文本文件的编码为 UTF-8
tools-go.exe -T -i D:\影视 -f ?.ass.srt
```

## 怎样新加功能

1. 在`大写为操作`一栏中，增加表示新功能的大写字母；在`小写为参数`一栏中，增加代表其操作的小写字母（可复用已有的小写字母作为参数）。为了避免重复选择字母，最好依字母顺序排序
2. 在`init`函数中，添加操作和参数的说明。其中操作的默认值为`bool`类型，参数的默认值多为`string`类型
3. 在`main`函数中，在最后的判断语句中，添加新操作的分支，根据操作传递适当的参数给实际执行函数
