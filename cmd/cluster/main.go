package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/evandigby/grmln"
)

func main() {
	c := grmln.NewCluster(
		grmln.ClusterConfig{
			OnConnectError: func(addr string, err error, attempts int) {
				log.Printf("Error connecting to %s (total attempts %d): %v", addr, attempts, err)
			},
		},
		"ws://localhost:8182/gremlin",
	)
	defer c.Close()

	op := grmln.NewOperator(c)

	err := op.EvalDefault(context.Background(), `g.V()`, func(resp *grmln.Response) {
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
