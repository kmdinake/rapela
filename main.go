package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BibleService struct {
	verseOfTheDay  BibleVerse
	currentVersion BibleVersion
}

// Get a random verse from the bible
func (bs *BibleService) GetVerseOfTheDay() BibleVerse { // Notice that we are using a pointer receiver here, which allows us to modify the state of the BibleService struct
	book := strings.ToLower("John")
	chapter := 3
	verse := 16

	fmt.Println("Getting the verse...")
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/%s/books/%s/chapters/%d/verses/%d.json", bs.currentVersion.Version, book, chapter, verse)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching verse: %v\n", err)
		return BibleVerse{}
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("Error reading verse body: %v\n", readErr)
		return BibleVerse{}
	}

	verseOfTheDay := BibleVerse{
		Id:      fmt.Sprintf("%s-%d-%d", book, chapter, verse),
		Book:    strings.ToTitle(book),
		Chapter: chapter,
		Verse:   fmt.Sprintf("%d", verse),
		Text:    "",
	}

	jsonErr := json.Unmarshal(body, &verseOfTheDay)
	if jsonErr != nil {
		fmt.Printf("Error parsing verse JSON: %v\n", jsonErr)
		return BibleVerse{}
	}

	return verseOfTheDay
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
	Id      string `json:"id"`
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   string `json:"verse"`
	Text    string `json:"text"`
}

func (v BibleVerse) String() string {
	return fmt.Sprintf("Verse of the day: %s %d:%s - %s", v.Book, v.Chapter, v.Verse, v.Text)
}

func main() {
	var bs = &BibleService{
		currentVersion: BibleVersion{
			Name:     "King James Version",
			Version:  "en-kjv",
			Language: "English",
		},
	}
	var verse = bs.GetVerseOfTheDay()

	fmt.Println("Dumela ngwana waka! Tseya sebaka se go rapela.")
	fmt.Println(verse)
}
