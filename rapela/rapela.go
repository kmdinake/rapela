package rapela

import "github.com/kmdinake/rapela/internal/bebele"

type BibleService = bebele.BibleService

type BibleVersion = bebele.BibleVersion

type BibleVerse = bebele.BibleVerse

type BibleChapter = bebele.BibleChapter

type BibleBook = bebele.BibleBook

var bibleService BibleService

func NewBibleService(bibleVersionOptional ...string) (BibleService, error) {
	if bibleService == nil {
		bibleVersion := "en-kjv"
		versionIndex := 0;
		if len(bibleVersionOptional) > 0 && bibleVersionOptional[versionIndex] != "" {
			bibleVersion = bibleVersionOptional[versionIndex]
		}
		var err error
		bibleService, err = bebele.NewBibleService(bibleVersion)
		if err != nil {
			return nil, err
		}
	}
	return bibleService, nil
}
