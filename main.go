package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Messages struct {
	Text string
	Num  int
}

func ConectRoute() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("*.html")

	return router
}

func GetCollection() *mongo.Collection {
	//mongoDBのクライアント作成＋mongoDBに接続
	mongoCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ := mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))

	//mongoDBのDBのCollectionの取得
	collection := client.Database("GOChat").Collection("messages")

	return collection
}

func GetMessage(collection *mongo.Collection) []*Messages {
	// 検索オプション設定
	findOptions := options.Find()

	// 戻り値
	var results []*Messages

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Messages
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	return results
}

func MainChat(ctx *gin.Context) {

	collection := GetCollection()

	results := GetMessage(collection)

	ctx.HTML(200, "index.html", gin.H{"messages": results})

}

func SendMessage(ctx *gin.Context) {
	text := ctx.PostForm("message")
	message := Messages{text, 0}

	collection := GetCollection()

	//Insert the data.
	collection.InsertOne(context.TODO(), message)

	ctx.Redirect(302, "/")
}

func main() {

	router := ConectRoute()

	router.GET("/", MainChat)
	router.POST("/message", SendMessage)

	router.Run()
}
