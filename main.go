package main

import "goSpiderProject/spider"

/**
 * 从[https://www.biquge.com] 爬取小说
 */
func main() {
	//spider.SaveToTxt("https://www.biquge.com/135_135747/")
	spider.SaveToMongo("https://www.biquge.com/135_135747/")
}
