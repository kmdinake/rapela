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
	defaultBibleName := "King James Version of the Holy Bible"
	defaultBibleVersion := "KJV"
	if bibleServiceSingleton == nil || (bibleServiceSingleton.CurrentVersion.Name != defaultBibleName && bibleServiceSingleton.CurrentVersion.Version != defaultBibleVersion) {
		bibleServiceSingleton = bebele.NewBibleService(defaultBibleName, defaultBibleVersion)
	}
	return bibleServiceSingleton
}
