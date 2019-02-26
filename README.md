# GRMLN

GRMLN is a Graph Driver Provider for Go based on the [TinkerPop Gremlin specs](http://tinkerpop.apache.org/docs/current/dev/provider/#_graph_driver_provider_requirements).

## Usage

### 1. Create a Connection

The first step is to create a connection:

#### Single Server

```go
c, err := grmln.Dial(
    context.Background(),
    "ws://localhost:8182/gremlin",
    grmln.DefaultMimeType,
)
if err != nil {
    log.Fatalf("Error connecting: %v", err)
}
defer c.Close()
```

#### Clustered

```go
c := grmln.NewCluster(
    grmln.ClusterConfig{
        OnConnectError: func(addr string, err error, attempts int) {
            log.Printf("Error connecting to %s (total attempts %d): %v", addr, attempts, err)
        },
    },
    "ws://localhost:8182/gremlin",
)
defer c.Close()
```

### 2. Create and Utilize and Operator

The second step is to wrap that connection with an operator that helps you with queries. 

```go
op := grmln.NewOperator(c)

err = op.EvalDefault(context.Background(), `g.V()`, func(resp *grmln.Response) {
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