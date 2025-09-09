package main

import (
    "fmt"
    "context"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
    ctx := context.Background()
    // URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
		dbUri := "neo4j://localhost"
    dbUser := "neo4j"
    dbPassword := "testpassword"
    driver, err := neo4j.NewDriverWithContext(
        dbUri,
        neo4j.BasicAuth(dbUser, dbPassword, ""))
    defer driver.Close(ctx)

    err = driver.VerifyConnectivity(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connection established.")
}
