package rapela

import "github.com/kmdinake/rapela/internal/bebele"

type BibleService = bebele.BibleService

type BibleVersion = bebele.BibleVersion

type BibleVerse = bebele.BibleVerse

var bibleServiceSingleton *BibleService

func NewBibleService(bibleName string, bibleVersion string) *BibleService {
	if bibleServiceSingleton == nil || (bibleServiceSingleton.CurrentVersion.Name != bibleName && bibleServiceSingleton.CurrentVersion.Version != bibleVersion) {
		bibleServiceSingleton = bebele.NewBibleService(bibleName, bibleVersion)
	}
	return bibleServiceSingleton
}

func NewDefaultBibleService() *BibleService {
	if bibleServiceSingleton == nil || (bibleServiceSingleton.CurrentVersion.Name != "King James" && bibleServiceSingleton.CurrentVersion.Version != "kjv") {
		bibleServiceSingleton = bebele.NewBibleService("King James", "kjv")
	}
	return bibleServiceSingleton
}
