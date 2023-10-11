package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Path    string `json:"path"`
	Workers int    `json:"workers"`
	Timeout int    `json:"timeout"`
	Outpath string `json:"outpath"`
}

func NewConfig(path string, workers, timeout int, outpath string) Config {
	return Config{
		Path:    path,
		Workers: workers,
		Timeout: timeout,
		Outpath: outpath,
	}
}

// 全局变量
var CONFIG Config
var LOG *log.Logger

func InitConf() {
	c := flag.Bool("c", false, "是否使用配置文件，默认否")
	path := flag.String("p", "", "书源url/本地路径")
	workers := flag.Int("w", 8, "并发数")
	timeout := flag.Int("t", 3, "超时时间,单位秒")
	outpath := flag.String("o", "", "输出路径")
	flag.Parse()

	if *c {
		fmt.Println("从config.json中读取配置,忽略其它参数")
		data, err := os.ReadFile("config.json")
		if err != nil {
			fmt.Println("config.json读取失败: ", err.Error())
			return
		}
		if err := json.Unmarshal(data, &CONFIG); err != nil {
			fmt.Println("config.json解析失败: ", err.Error())
			return
		}
	} else {
		if *path == "" {
			fmt.Println("请输入配置文件路径或书源路径")
			flag.PrintDefaults()
			return
		}

		CONFIG = NewConfig(*path, *workers, *timeout, *outpath)

		fmt.Println("命令行参数写入到config.json中...")
		data, _ := json.MarshalIndent(&CONFIG, "", "  ")
		err := os.WriteFile("config.json", data, 0644)
		if err != nil {
			fmt.Println("写入 config.json 失败: ", err.Error())
			return
		}
	}

	//初始化结果输出目录
	if CONFIG.Outpath != "" {
		if _, err := os.Stat(CONFIG.Outpath); os.IsNotExist(err) {
			os.MkdirAll(CONFIG.Outpath, os.ModePerm)
		}
	}
	//初始化Logger
	f, err := os.OpenFile(filepath.Join(CONFIG.Outpath, "log.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("日志文件初始化失败: ", err.Error())
	}
	LOG = log.New(f, "", log.LstdFlags)

}
