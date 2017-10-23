package main

import (
	"fmt"

	"github.com/tett23/mangrove/lib/mangrove_db"
	"github.com/tett23/mangrove/models"
)

func main() {
	if _, err := mangrove_db.InitDatabase("local"); err != nil {
		panic(err)
	}

	var videos models.Videos
	if err := videos.Latest(); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", videos)
}
