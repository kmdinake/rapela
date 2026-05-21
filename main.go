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

	version, err := bs.GetBibleVersion()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current Bible Version: %s\n", version)

	verseOfTheDay := bs.GetVerseOfTheDay()
	fmt.Printf("Dumela ngwana waka! Tseya sebaka se go rapela.\nVerse Of The Day: %s\n", verseOfTheDay)

	versions, err := bs.GetBibleVersions()
	if err != nil {
		panic(err)
	}
	fmt.Println("Available Bible Versions:")
	for _, version := range versions {
		fmt.Printf("- %s\n", version)
	}

	randomVersion := versions[rand.IntN(len(versions))]
	fmt.Printf("Setting Bible version to: %s\n", randomVersion)
	bs.SetBibleVersion(randomVersion)

	version, err = bs.GetBibleVersion()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current Bible Version after setting: %s\n", version)
}
