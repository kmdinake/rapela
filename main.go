package main

import "fmt"

type BibleService struct {
	verseOfTheDay  BibleVerse
	currentVersion BibleVersion
}

// Get a random verse from the bible
func (bs *BibleService) GetVerseOfTheDay() BibleVerse { // Notice that we are using a pointer receiver here, which allows us to modify the state of the BibleService struct
	// In a real implementation, this would likely involve some logic to select a random verse from a database or API
	fmt.Println("Getting the verse...")
	return BibleVerse{
		Id:      1,
		Book:    "John",
		Chapter: 3,
		Verse:   16,
		Text:    "For God so loved the world that he gave his one and only Son, that whoever believes in him shall not perish but have eternal life.",
	}
}

// Get a list of available Bible versions
func (bs *BibleService) GetBibleVersions() []BibleVersion {
	panic("Method not implemented")
}

// Set the current Bible version to be used
func (bs *BibleService) SetBibleVersion(version BibleVersion) {
	panic("Method not implemented")
}

// Get the current Bible version being used
func (bs *BibleService) GetBibleVersion() BibleVersion {
	panic("Method not implemented")
}

// Get a specific verse by its reference (book, chapter, verse)
func (bs *BibleService) GetVerseByReference(book string, chapter int, verse int) BibleVerse {
	panic("Method not implemented")
}

type BibleVersion struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"`
}

type BibleVerse struct {
	Id      int    `json:"id"`
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

func (v BibleVerse) String() string {
	return fmt.Sprintf("Verse: %d, %s %d:%d - %s", v.Id, v.Book, v.Chapter, v.Verse, v.Text)
}

func main() {
	var bs = &BibleService{}
	var verse = bs.GetVerseOfTheDay()

	fmt.Println("Dumela ngwana waka! Tseya sebaka se go rapela.")
	fmt.Println(verse)
}
