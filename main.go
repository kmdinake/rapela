package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/kmdinake/rapela/rapela"
)

func main() {
	bs, err := rapela.NewBibleService()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bible Service created successfully!")

	fmt.Printf("Current Bible Version: %s\n", bs.GetBibleVersion())

	fmt.Printf("Dumela ngwana waka! Tseya sebaka se go rapela.\n%s\n", bs.GetVerseOfTheDay())

	versions := bs.GetBibleVersions()
	fmt.Println("Available Bible Versions:")
	for _, version := range versions {
		fmt.Println(version)
	}

	randomVersion := versions[rand.IntN(len(versions))]
	fmt.Printf("Setting Bible version to: %s\n", randomVersion)
	bs.SetBibleVersion(randomVersion)

	fmt.Printf("Current Bible Version after setting: %s\n", bs.GetBibleVersion())
}
