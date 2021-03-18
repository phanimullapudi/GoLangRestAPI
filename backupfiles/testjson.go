package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type Game struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Platform    []string `json:"platform"`
}

func main() {

    myGame := Game{
        Title:       "Fifa 19",
        Description: "Football simulation game, based on UEFA players",
        Platform:    []string{"PS4"},
    }

    // For comparison, the usual way would be: j, err := json.Marshal(myGame)

    // MarshalIndent accepts:
    // 1) the data
    // 2) a prefix to place on all lines but 1st
    // 3) an indent to place before lines based on indent level

    // Print Json with indents, the pretty way:
    prettyJSON, err := json.MarshalIndent(myGame, "", "    ")
    if err != nil {
        log.Fatal("Failed to generate json", err)
    }
    fmt.Printf("%s\n", string(prettyJSON))
}