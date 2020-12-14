package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goSpiderProject/spider"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"strings"
)

//是否已gbk编码保存
var saveAsGBK bool = true

func main() {
	// Set client options
	var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var client, _ = mongo.Connect(context.TODO(), clientOptions)
	var novelCol = client.Database("noveldb").Collection("novel")
	var chapterlCol = client.Database("noveldb").Collection("chapters")
	var name string
	fmt.Printf("输入书名：")
	fmt.Scanf("%s", &name)
	name = strings.Trim(name, "\n")

	var result bson.M
	if err := novelCol.FindOne(context.TODO(), bson.M{"name": name}).Decode(&result); err != nil {
		fmt.Printf("没有找到%q\n", name)
		fmt.Printf("%v\n", result)
		fmt.Println(err)
	} else {
		bookID := result["_id"].(primitive.ObjectID).String()
		name := result["name"].(string)
		opts := options.Find().SetSort(bson.D{{"id", 1}})
		if cursor, err := chapterlCol.Find(context.TODO(), bson.M{"bookid": bookID}, opts); err != nil || cursor == nil {
			fmt.Println(err)
		} else {
			if f, err := os.OpenFile(name+"_gbk.txt", os.O_WRONLY|os.O_CREATE, 0755); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("开始整合...")
				var result spider.Chapter

				encoder := simplifiedchinese.GBK.NewEncoder()
				for cursor.Next(context.TODO()) {
					if err := cursor.Decode(&result); err != nil {
						fmt.Println(err)
					} else {
						fmt.Printf("写入：%s\n", result.Title)
						text := result.Title + "\n" + result.Content + "\n"
						if saveAsGBK {
							if gbk, err := encoder.String(text); err != nil {
								fmt.Println(err)
								return
							} else {
								text = gbk
							}
						}

						if _, err := f.WriteString(text); err != nil {
							fmt.Println(err)
						}
					}

				}
			}
		}
	}
	client.Disconnect(context.TODO())
}
