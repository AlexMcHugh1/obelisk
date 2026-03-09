package main

import (
    "log"
    "obelisk/internal/database"
)

func main() {
    db, err := database.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
}