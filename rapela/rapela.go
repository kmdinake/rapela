package rapela

import "github.com/kmdinake/rapela/internal/bebele"

type BibleService = bebele.BibleService

type BibleVersion = bebele.BibleVersion

type BibleVerse = bebele.BibleVerse

type BibleChapter = bebele.BibleChapter

type BibleBook = bebele.BibleBook

var bibleService BibleService

func NewBibleService() (BibleService, error) {
	if bibleService == nil {
		var err error
		bibleService, err = bebele.NewBibleService("en-kjv")
		if err != nil {
			return nil, err
		}
	}
	return bibleService, nil
}
