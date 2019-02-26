# GRMLN

GRMLN is a Graph Driver Provider for Go based on the [TinkerPop Gremlin specs](http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements).

## Example Usage

```go
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
```