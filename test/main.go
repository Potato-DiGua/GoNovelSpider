package main

import (
	"container/list"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"goSpiderProject/spider"
	"strings"
)

func main() {
	list.New()
	c := colly.NewCollector(
		colly.AllowedDomains("www.biquge.com"), // 允许爬取的域名
	)
	// 获取页面内容后的回调
	c.OnResponse(func(res *colly.Response) {
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(res.Body)))
		contentDom := dom.Find("#content").First()
		fmt.Printf("%q\n", contentDom.Text())
		// 删除正文里的子标签
		contentDom.Children().Each(func(i int, selection *goquery.Selection) {
			selection.Remove()
		})
		fmt.Println(spider.HandleText(contentDom.Text()))
	})
	c.Visit("https://www.biquge.com/135_135747/7738977.html")
}
