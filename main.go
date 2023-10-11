package main

import (
	"fmt"
	"learn/verifybooksource/book"
	"learn/verifybooksource/conf"
	"path/filepath"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	conf.InitConf()
	conf.LOG.Println("==============VerifyBooksource==============")
	originalBooksourcelist, err := book.NewBSList(conf.CONFIG)
	if err != nil {
		fmt.Println("书源地址解析失败: ", err)
		return
	}
	start := time.Now()
	goodBSL := []book.BookSource{}
	badBSL := []book.BookSource{}

	bslChannel := make(chan book.BookSource, len(*originalBooksourcelist))
	for _, bs := range *originalBooksourcelist {
		bslChannel <- bs
	}
	close(bslChannel)

	total := len(*originalBooksourcelist)
	var mu sync.Mutex
	processed := 0
	fmt.Println("并发验证书源...")
	var wg sync.WaitGroup
	wg.Add(conf.CONFIG.Workers)
	for i := 0; i < conf.CONFIG.Workers; i++ {
		go func() {
			for bs := range bslChannel {
				mu.Lock()
				processed += 1
				mu.Unlock()
				res := bs.Check(conf.CONFIG.Timeout)
				if res {
					goodBSL = append(goodBSL, bs)
				} else {
					badBSL = append(badBSL, bs)
				}
			}
			wg.Done()
		}()
	}
	bar := progressbar.NewOptions(total, progressbar.OptionSetWidth(10), progressbar.OptionShowCount(), progressbar.OptionSetPredictTime(false))

	for processed < total {
		bar.Set(processed)
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	cost := time.Since(start)
	fmt.Println("--------------------")
	book.SaveToJson(&goodBSL, filepath.Join(conf.CONFIG.Outpath, "good.json"))
	book.SaveToJson(&badBSL, filepath.Join(conf.CONFIG.Outpath, "bad.json"))
	fmt.Println("--------------------")
	fmt.Println("")
	fmt.Println("成果报表:")
	fmt.Println("书源总数: ", total)
	fmt.Println("有效书源数: ", len(goodBSL))
	fmt.Println("无效书源数: ", len(badBSL))
	fmt.Printf("耗时: %f s\n", cost.Seconds())
}
