package superMongo

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type example struct {
	Word string `json:"word"`
	Name string `json:"name"`
}

func MongoTest() {
	var data example
	cli := NewClient("localhost:27017")
	con := cli.OpenWithContext(context.Background(), "test", "test")
	con.Find(bson.D{{"name", "SuperYong"}}).One(&data)
	fmt.Printf("%+v", data)
}
