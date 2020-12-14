package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gobwas/glob/util/runes"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
)

type OnSaveNovelInfo func(book *Book) error
type OnSaveChapter func(chapter *Chapter) error

func Crawl(chaptersUrl string, onSaveNovelInfo OnSaveNovelInfo, onSaveChapter OnSaveChapter) {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.biquge.com"), // 允许爬取的域名
		colly.Async(true),                      // 异步模式
	)

	var contentCollector *colly.Collector

	// 获取页面内容后的回调
	c.OnResponse(func(res *colly.Response) {
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(res.Body)))
		author := []rune(HandleText(dom.Find("#info > p:nth-child(2)").First().Text()))
		author = author[runes.LastIndex(author, []rune("："))+1:]
		book := Book{
			Name:         dom.Find("#info > h1").First().Text(),
			Author:       string(author),
			Introduction: HandleText(dom.Find("#intro").Text()),
			Url:          res.Request.URL.String(),
		}
		if onSaveNovelInfo != nil {
			if err := onSaveNovelInfo(&book); err != nil {
				fmt.Println(err)
				return
			}
		}

		contentCollector = getContentCollector(onSaveChapter)

		fmt.Printf("Link found: %v\n", book)

		dom.Find("#list > dl > dd > a").Each(func(i int, selection *goquery.Selection) {
			link, exist := selection.Attr("href")
			if exist {
				contentCollector.Visit(res.Request.AbsoluteURL(link))
			}

		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit(chaptersUrl)
	c.Wait()
	contentCollector.Wait()
}

func getContentCollector(onSaveChapter OnSaveChapter) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("www.biquge.com"),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5, //并发数
		//RandomDelay: 5 * time.Second,
		//Delay:       3 * time.Second,
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		link := e.Request.URL.String()

		idString := link[strings.LastIndex(link, "/")+1 : strings.LastIndex(link, ".")]
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		contentDom := e.DOM.Find("#content").First()
		// 删除正文里的子标签
		contentDom.Children().Each(func(i int, selection *goquery.Selection) {
			selection.Remove()
		})

		chapter := Chapter{
			Url:     link,
			Title:   e.DOM.Find("#wrapper > div.content_read > div.box_con > div.bookname > h1").First().Text(),
			Content: contentDom.Text(),
			ID:      id,
		}
		// Print link
		chapter.Content = HandleText(chapter.Content)

		if onSaveChapter != nil {
			if err := onSaveChapter(&chapter); err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Printf("【%s】(%s)下载完成\n", chapter.Title, chapter.Url)
	})

	return c
}

/**
 * 清理正文
 */
func HandleText(text string) string {
	if re, err := regexp.Compile("[\t\u3000\n ]+"); err != nil {
		return ""
	} else {
		return strings.Trim(re.ReplaceAllString(text, "\n"), "\n\t ")
	}
}
