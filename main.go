package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017"

type Item struct {
	Category string `json:"category"`
	Model    string `json:"model"`
	Price    int    `json:"price"`
	Producer string `json:"producer"`
}

func FindAll(client *mongo.Client) {
	find, err := client.Database("test").Collection("shop").Find(context.TODO(), bson.D{})
	if err != nil {
		return
	}

	var Result []Item
	for find.Next(context.TODO()) {
		var res Item
		err = find.Decode(&res)
		if err != nil {
			return
		}
		Result = append(Result, res)
	}
	final, _ := json.Marshal(Result)
	fmt.Println(string(final))
}

func CountItemsInCategory(client *mongo.Client) {

	filter := bson.M{"category": "Phone"}

	find, err := client.Database("test").Collection("shop").CountDocuments(context.TODO(), filter)
	if err != nil {
		return
	}
	fmt.Printf("Found %d items in %s category \n", find, filter["category"])
}

func CountCategories(client *mongo.Client) {
	filter := "category"

	find, err := client.Database("test").Collection("shop").Distinct(context.TODO(), filter, bson.D{})
	if err != nil {
		return
	}
	counter := len(find)
	fmt.Printf("There are %d categories \n", counter)
}

func CountAllProviders(client *mongo.Client) {
	filter := "producer"

	find, err := client.Database("test").Collection("shop").Distinct(context.TODO(), filter, bson.D{})
	if err != nil {
		return
	}
	fmt.Println("producers:")
	for _, v := range find {
		fmt.Println(v)
	}
}

func And(client *mongo.Client) {
	filter := bson.D{{"category", "Phone"}, {"price", 1000}}
	find, err := client.Database("test").Collection("shop").Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	var results []Item
	if err = find.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Printf("Filter %s is %v and %s is %v and result %v \n", filter[0].Key, filter[0].Value, filter[1].Key, filter[1].Value, result)
	}

}

func Or(client *mongo.Client) {
	filter := bson.D{
		{"$or", []interface{}{
			bson.M{"model": "LG2"},
			bson.M{"model": "Samsung 21"},
		}},
	}
	find, err := client.Database("test").Collection("shop").Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return
	}
	var results []Item
	if err = find.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	fmt.Println("$or function")
	for _, result := range results {
		fmt.Println(result)
	}
}

func In(client *mongo.Client) {
	filter := bson.D{
		{"producer", bson.D{{"$in", []interface{}{
			"LG",
			"Huawei",
		}},
		}},
	}
	find, err := client.Database("test").Collection("shop").Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return
	}
	var results []Item
	if err = find.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	fmt.Println("$in function")
	for _, result := range results {
		fmt.Println(result)
	}
}

func Update(client *mongo.Client) {
	filter1 := bson.D{{"price", bson.D{{"$lt", 500}}}}
	update1 := bson.D{{"$set", bson.D{{"RAM", 2}}}}

	_, err := client.Database("test").Collection("shop").UpdateMany(context.TODO(), filter1, update1)
	if err != nil {
		fmt.Println(err)
		return
	}

	filter2 := bson.D{{"price", bson.D{{"$gte", 500}}}}
	update2 := bson.D{{"$set", bson.D{{"RAM", 4}}}}

	_, err = client.Database("test").Collection("shop").UpdateMany(context.TODO(), filter2, update2)

}
func FindByRam(client *mongo.Client) {
	filter := bson.D{{"RAM", 4}}
	find, err := client.Database("test").Collection("shop").Find(context.TODO(), filter)

	var Result []Item
	for find.Next(context.TODO()) {
		var res Item
		err = find.Decode(&res)
		if err != nil {
			return
		}
		Result = append(Result, res)
	}
	final, _ := json.Marshal(Result)
	fmt.Println(string(final))
}

func IncreasePrice(client *mongo.Client) {

	filter := bson.D{{"RAM", 4}}
	find, err := client.Database("test").Collection("shop").Find(context.TODO(), filter)
	if err != nil {
		return
	}
	for find.Next(context.TODO()) {
		_, err = client.Database("test").Collection("shop").UpdateMany(context.TODO(), find.Current, bson.D{{"$mul", bson.D{{"price", 1.15}}}})

	}

}
func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	FindAll(client)
	CountItemsInCategory(client)
	CountCategories(client)
	CountAllProviders(client)
	And(client)
	Or(client)
	In(client)
	Update(client)
	FindByRam(client)
	IncreasePrice(client)
}
