package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/evandigby/grmln"
)

func main() {
	c, err := grmln.Dial(
		context.Background(),
		"ws://localhost:8182/gremlin",
		grmln.DefaultMimeType,
	)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer c.Close()

	op := grmln.NewOperator(c)
	sop := op.NewSession()
	defer sop.CloseDefault(context.Background())

	err = sop.EvalDefault(context.Background(), `g.V()`, func(resp *grmln.Response) {
		data, err := json.MarshalIndent(resp.Result.Data, "", "   ")
		if err != nil {
			log.Fatalf("Error marshalling response: %v", err)
		}

		fmt.Println(string(data))
	})

	if grmln.IsUnauthorized(err) {
		log.Fatal("you are not authorized!")
	}

	if err != nil {
		log.Fatal("Error: ", err)
	}
}
