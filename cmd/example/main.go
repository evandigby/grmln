package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/evandigby/grmln"
)

func main() {
	c, err := grmln.Dial("ws://localhost:8182/gremlin")
	if err != nil {
		log.Fatal(err)
	}

	err = c.EvalDefault(`g.V()`, func(resp *grmln.Response) {
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
