package main

import (
	"encoding/json"
	"fmt"

	"github.com/newtoallofthis123/sqlite_http/api"
	"github.com/newtoallofthis123/sqlite_http/db"
)

func main() {
	listenAddr := ":8080"
	api := api.NewApi(listenAddr)

	d, err := db.NewDb("test.db")
	if err != nil {
		panic(err)
	}
	err = d.Discover()
	if err != nil {
		panic(err)
	}

	rows, err := d.RunQuery("SELECT * from hello")
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(rows)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	api.Run()
}
