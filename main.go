/*
Copyright © 2026 Keoagile Dinake kmdinake@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import "github.com/kmdinake/rapela/cmd"

func main() {
	cmd.Execute()
}

/*
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

	verseOfTheDay = bs.GetVerseOfTheDay()
	fmt.Printf("Dumela ngwana waka! Tseya sebaka se go rapela.\nVerse Of The Day: %s\n", verseOfTheDay)
}
*/
