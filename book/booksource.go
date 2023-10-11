package book

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"learn/verifybooksource/conf"
	"net/http"
	"os"
	"strings"
	"time"
)

type BookSource map[string]interface{}

func (bs BookSource) Check(timeout int) bool {
	url, ok := bs["bookSourceUrl"]
	if !ok {
		conf.LOG.Println("读取bookSourceUrl失败")
		return false
	}
	conf.LOG.Println("检查: ", url)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url.(string))
	if err != nil {
		conf.LOG.Println("请求失败:", url, err.Error())
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
func NewBSList(config conf.Config) (*[]BookSource, error) {
	var originalBooksourcelist []BookSource

	if strings.HasPrefix(config.Path, "http") {
		//url
		fmt.Println("读取书源url...")
		// Ignore SSL certificate verification
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.Get(config.Path)
		if err != nil {
			fmt.Println("书源下载失败,请确认书源地址是否正确: ", config.Path)
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取书源失败: ", err.Error())
			return nil, err

		}
		if err := json.Unmarshal(body, &originalBooksourcelist); err != nil {
			fmt.Println("书源解析失败: ", err.Error())
			return nil, err

		}
		fmt.Println("书源读取完成")
		fmt.Println("书源内容写入文件 originalBooksourcelist.json ...")
		err = os.WriteFile("originalBooksourcelist.json", body, 0644)
		if err != nil {
			fmt.Println("书源内容写入文件失败:", err.Error())
			return nil, err

		}
	} else {
		//file
		fmt.Println("读取书源文件...")
		data, err := os.ReadFile(config.Path)
		if err != nil {
			fmt.Println("书源文件打开失败: ", config.Path)
			return nil, err

		}
		if err := json.Unmarshal(data, &originalBooksourcelist); err != nil {
			fmt.Println("书源解析失败: ", err.Error())
			return nil, err

		}
		fmt.Println("书源读取完成")

	}
	return &originalBooksourcelist, nil
}

func SaveToJson(bsl *[]BookSource, path string) {
	fmt.Println("保存结果到 ", path)
	if len(*bsl) == 0 {
		fmt.Println("内容为空，无需保存")
		return
	}
	data, err := json.MarshalIndent(bsl, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data: ", err)
		return
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}
	fmt.Println("保存结果完成")

}
