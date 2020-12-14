package spider

import (
	"encoding/json"
	"fmt"
)

type Book struct {
	Name         string
	Author       string
	Introduction string
	Url          string
}

type Chapter struct {
	Url     string
	Title   string
	Content string
	ID      int64
	BookID  string
}

func (c *Chapter) toJson() string {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}
