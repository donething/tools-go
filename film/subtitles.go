package film

import (
	"encoding/json"
	"fmt"
	"github.com/donething/utils-go/dohttp"
	"github.com/donething/utils-go/dotext"
	"github.com/gookit/color"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SubResp 获取字幕列表的响应
type SubResp struct {
	Code int `json:"code"` // 为 0 表示没有错误
	Data []struct {
		Gcid      string   `json:"gcid"`      // 字幕存储 ID
		Cid       string   `json:"cid"`       // 字幕存储 ID
		URL       string   `json:"url"`       // 字幕下载地址
		Ext       string   `json:"ext"`       // 字幕扩展名。如"srt"
		Name      string   `json:"name"`      // 字幕文件名。如"阿甘正传.srt"
		Duration  int      `json:"duration"`  // 时长。该值不准确
		Languages []string `json:"languages"` // 字幕语言。如["简体","繁体"]，该值不准确
	} `json:"data"`
}

const (
	// 下载字幕的地址，最后为搜索关键字
	subURL = "https://api-shoulei-ssl.xunlei.com/oracle/subtitle?gcid=&cid=&name=%s"
)

var (
	httpclient = dohttp.New(30*time.Second, false, false)
	headers    = map[string]string{
		"Host":   "api-shoulei-ssl.xunlei.com",
		"Accept": "application/json, text/plain, */*",
		"User-Agent": "thunder/11.3.7.1880 Mozilla/5.0 (Windows NT 10.0; WOW64) " +
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.122 XDASKernel/9.2.1 Safari/537.36",
		"Accept-Language": "zh-CN",
	}
)

// DLSubtitle 下载字幕
func DLSubtitle(key string, filmPath string) {
	// 如果关键字 key 为空，将提取关键字来搜索
	if key == "" {
		key = dotext.ResolveFanhao(filmPath)
		color.Info.Tips("提取到的字幕关键字：'%s'，准备搜索\n", key)
	}
	// 没有匹配到搜索关键字，返回
	if key == "" {
		color.Error.Tips("下载字幕失败：没有匹配到搜索关键字'%s'\n", filmPath)
		return
	}

	// 发送请求
	u := fmt.Sprintf(subURL, url.QueryEscape(key))
	bs, err := httpclient.Get(u, headers)
	if err != nil {
		color.Error.Tips("下载字幕失败。网络出错：%s\n", err)
		return
	}

	// 解析
	var subResp SubResp
	err = json.Unmarshal(bs, &subResp)
	if err != nil {
		color.Error.Tips("下载字幕失败。解析字幕列表出错：'%s'。URL('%s')：'%s'\n", err, u, string(bs))
		return
	}

	// 提取下载地址
	if len(subResp.Data) == 0 {
		color.Warn.Tips("下载字幕失败。没有找到电影'%s'的字幕\n", key)
		return
	}

	// 列出结果
	color.Notice.Tips("查找到的字幕文件(%d个)", len(subResp.Data))
	fmt.Printf("%-2s\t%-5s\t\t%s\n", "编号", "语言", "字幕名")
	for i, item := range subResp.Data {
		// -：表示左对齐，0.5：表示占用 0 到 5 个字符宽度
		// 注：一个字母和一个汉字都算一个字符
		fmt.Printf("%4d\t%-0.5v\t\t%-30s\n", i+1, item.Languages, item.Name)
	}

	// 暂存待下载的目标字幕
	var data = &subResp.Data[0]
	if len(subResp.Data) >= 2 {
		var choice int
		fmt.Printf("\n请输入字幕编号，来下载指定字幕：")
		_, err = fmt.Scanln(&choice)
		if err != nil {
			color.Error.Tips("输入字幕编号时出错：%s\n", err)
			return
		}
		data = &subResp.Data[choice-1]
	}

	// 如果没有指定电影路径，则将字幕保存到当前执行路径下；否则将字幕重命名为电影名，并保存到电影所在目录下
	savePath := ""
	if filmPath != "" {
		savePath = filmPath[0:strings.LastIndex(filmPath, ".")] + filepath.Ext(data.Name)
	} else {
		// 获取当前执行路径，将保存字幕文件到该路径下
		curPath, err := os.Getwd()
		if err != nil {
			color.Error.Tips("下载字幕失败。获取当前执行路径出错：%s\n", err)
			return
		}
		savePath = filepath.Join(curPath, data.Name)
	}

	color.Info.Tips("开始下载字幕文件'%s'，语言 %v\n", data.Name, data.Languages)
	color.Info.Tips("字幕的下载地址 '%s'\n", data.URL)

	// 下载字幕
	subBS, err := httpclient.Get(data.URL, headers)
	if err != nil {
		color.Error.Tips("下载字幕文件'%s'出错：%s\n", data.URL, err)
		return
	}

	// 将字幕文件编码转为 UTF-8
	subBS, encoding, err := dotext.Text2UTF8(subBS)
	if err != nil {
		color.Error.Tips("下载字幕失败。转换字幕编码'%s'到'UTF-8'出错：%s\n", encoding, err)
		return
	}

	// 判断是否确实为字幕，验证关键字：.srt 需包含"-->"、.ass 需包含"Format"，才为字幕文件
	if strings.Index(string(subBS), "-->") == -1 &&
		strings.Index(string(subBS), "Format") == -1 {
		color.Warn.Tips("下载字幕失败。经验证，该项不是字幕文件，取消下载")
		return
	}

	// 保存到本地
	err = ioutil.WriteFile(savePath, subBS, 0644)
	if err != nil {
		color.Error.Tips("下载字幕失败。保存字幕文件'%s'出错：%s\n", savePath, err)
		return
	}

	color.Success.Tips("已完成 将字幕'%s'保存到本地'%s'\n", data.Name, savePath)
}
