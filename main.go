package main

import "log"

func main() {

	store, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := store.Init(); err != nil {
		log.Fatal(err.Error())
	}

	apiServer := NewServer(":8080", store)
	apiServer.Run()
}
