package spider

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sort"
)

type chapterList []*Chapter

func SaveToTxt(url string) {
	var chapters chapterList
	var bookName string
	Crawl(url, func(book *Book) error {
		bookName = book.Name
		return nil
	}, func(chapter *Chapter) error {
		chapters = append(chapters, chapter)
		return nil
	})
	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].ID < chapters[j].ID
	})

	if file, err := os.OpenFile(bookName+".txt", os.O_CREATE|os.O_WRONLY, 0755); err != nil {
		fmt.Println(err)
	} else {
		for _, chapter := range chapters {
			if _, err := file.WriteString(chapter.Title + "\n" + chapter.Content + "\n"); err != nil {
				fmt.Println(err)
				return
			}
		}
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}

}

func SaveToMongo(url string) {
	var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接mongo数据库
	var client, _ = mongo.Connect(context.TODO(), clientOptions)
	var novelCol = client.Database("noveldb").Collection("novel")
	var chapterCol = client.Database("noveldb").Collection("chapters")

	var bookID string
	Crawl(url, func(book *Book) error {
		//-------------------获取booID------------------------
		var result bson.M
		if err := novelCol.FindOne(context.TODO(), bson.M{"name": book.Name}).Decode(&result); err != nil {
			fmt.Println(err)
			if res, err := novelCol.InsertOne(context.TODO(), book); err != nil {
				return err
			} else {
				bookID = res.InsertedID.(primitive.ObjectID).String()
				return nil
			}
		} else {
			bookID = result["_id"].(primitive.ObjectID).String()
		}
		return nil
		//-------------------获取booID------------------------
	}, func(chapter *Chapter) error {
		chapter.BookID = bookID
		if _, err := chapterCol.InsertOne(context.TODO(), chapter); err != nil {
			fmt.Println(err)
			return err
		} else {
			return nil
		}
	})
	if err := client.Disconnect(context.TODO()); err != nil {
		fmt.Println(err)
	}
}
