package main

import (
	"fmt"

	"github.com/kmdinake/rapela/rapela"
)

func main() {
	bs := rapela.NewDefaultBibleService()
	versions := bs.GetBibleVersions()
	fmt.Println("Available Bible Versions:")
	for _, version := range versions {
		fmt.Printf("- %s: %s (%s)\n", version.Id, version.Name, version.Version)
	}
	// var verse = bs.GetVerseOfTheDay()

	// fmt.Println("Dumela ngwana waka! Tseya sebaka se go rapela.")
	// fmt.Println(verse)
}
