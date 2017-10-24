package main

import (
	"fmt"

	"github.com/tett23/mangrove/lib/mangrove_db"
	"github.com/tett23/mangrove/lib/name_resolver"
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
	// fmt.Printf("%+v\n", videos)

	for _, item := range videos {
		p := name_resolver.CreateInstance(item.OriginalName)
		_, err := p.AddProgramData(item.Program)
		if err != nil {
			// panic(err)
		}

		n, err := p.GetName()
		if err != nil {
			panic(err)
		}
		if p.Title == "" {

			fmt.Println("original\t", item.OriginalName)
			fmt.Println("output\t\t", n)
			fmt.Println("---------")
		}
	}

	// storages, err := storage_balancer.LoadStorages()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("%+v\n", storages)
}
