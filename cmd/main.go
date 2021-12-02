package main

import (
	"flag"
	"fmt"
	"github_api"
	"github_api/aggregate"
	"github_api/api"
	"log"
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "usernames", "usernames.txt", "File containing github usernames")
	flag.Parse()

	usernames, err := github_api.Read(filepath)
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewClient()
	aggregator := aggregate.New(client)

	for _, un := range usernames {
		res, err := aggregator.AggregateUser(un)
		if err != nil {
			log.Fatal(err)
		}
		aggregator.Print(res)
		fmt.Println()
	}
}
