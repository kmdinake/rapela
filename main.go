package main

import "fmt"

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

func getVerse() BibleVerse {
	fmt.Println("Getting the verse...")
	return BibleVerse{
		Id:      1,
		Book:    "John",
		Chapter: 3,
		Verse:   16,
		Text:    "For God so loved the world that he gave his one and only Son, that whoever believes in him shall not perish but have eternal life.",
	}
}

func main() {
	var verse = getVerse()

	fmt.Println("Dumela ngwana waka! Tseya sebaka se go rapela.")
	fmt.Println(verse)
}
