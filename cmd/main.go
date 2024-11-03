package main

import (
	"github.com/newtoallofthis123/sqlite_http/api"
)

func main() {
	listenAddr := ":8080"
	api := api.NewApi(listenAddr)
	api.Run()
}
