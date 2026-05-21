package bebele

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
)

// region Core Type definitions
type BibleVersion struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"`
}

func (v BibleVersion) String() string {
	return fmt.Sprintf("%s: %s (%s)", v.Id, v.Name, v.Version)
}

type BibleBook struct {
	Id string `json:"id"`
}

type BibleChapter struct {
	Id string `json:"id"`
}

type BibleVerse struct {
	Id      string       `json:"id"`
	Book    BibleBook    `json:"book"`
	Chapter BibleChapter `json:"chapter"`
	Text    string       `json:"text"`
}

func (v BibleVerse) String() string {
	return fmt.Sprintf("%s %s:%s - %s", v.Book.Id, v.Chapter.Id, v.Id, v.Text)
}

type BibleService interface {
	SetBibleVersion(version BibleVersion) error
	GetBibleVersion() (BibleVersion, error)
	GetBibleVersions() ([]BibleVersion, error)
	GetBibleVersionById(versionId string) (BibleVersion, error)
	GetBooksBy(versionId string) ([]BibleBook, error)
	GetChaptersBy(versionId string, bookId string) ([]BibleChapter, error)
	GetVersesBy(versionId string, bookId string, chapterId string) ([]BibleVerse, error)
	GetVerseBy(versionId string, bookId string, chapterId string, verseId string) (BibleVerse, error)
	GetVerseOfTheDay() BibleVerse
}

// endregion

// region Wldeh Bible Service Implementation
// API origin https://github.com/wldeh/bible-api
type WldehBibleService struct {
	verseOfTheDay  BibleVerse
	currentVersion BibleVersion
}

var httpGet = http.Get

func (bs WldehBibleService) SetBibleVersion(version BibleVersion) error {
	if bs.currentVersion.Id != version.Id {
		bs.currentVersion = version
		return nil
	}
	return fmt.Errorf("Bible version %s already set", version.Id)
}

func (bs WldehBibleService) GetBibleVersion() (BibleVersion, error) {
	if bs.currentVersion.Id == "" {
		return BibleVersion{}, fmt.Errorf("No bible version set")
	}
	return bs.currentVersion, nil
}

func (bs WldehBibleService) GetBibleVersions() ([]BibleVersion, error) {
	url := "https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/bibles.json"
	var data any
	if err := bs.fetchJsonDataFrom(url, &data); err != nil {
		return []BibleVersion{}, fmt.Errorf("Error fetching bible versions: %v\n", err)
	}

	return bs.convertToBibleVersions(data), nil
}

func (bs WldehBibleService) GetBibleVersionById(versionId string) (BibleVersion, error) {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/%s/%s.json", versionId, versionId)
	var data any
	if err := bs.fetchJsonDataFrom(url, &data); err != nil {
		return BibleVersion{}, fmt.Errorf("Error fetching bible version %s: %v\n", versionId, err)
	}

	version, err := bs.convertToBibleVersion(data)
	if err != nil {
		return BibleVersion{}, fmt.Errorf("Error converting bible version %s JSON: %v\n", versionId, err)
	}

	if version.Id != versionId {
		return BibleVersion{}, fmt.Errorf("Bible version not found: %s\n", versionId)
	}

	return version, nil
}

func (bs WldehBibleService) GetBooksBy(versionId string) ([]BibleBook, error) {
	fmt.Printf("Getting books for bible version %s\n", versionId)
	url := fmt.Sprintf("https://api.github.com/repos/wldeh/bible-api/contents/bibles/%s/books", versionId)
	var data any
	if err := bs.fetchJsonDataFrom(url, &data); err != nil {
		return []BibleBook{}, fmt.Errorf("Error fetching %s books: %v\n", versionId, err)
	}

	var books []BibleBook
	for _, b := range data.([]interface{}) {
		if item, ok := b.(map[string]interface{}); ok {
			book := BibleBook{}
			book.Id = item["name"].(string)
			books = append(books, book)
		}
	}

	return books, nil
}

func (bs WldehBibleService) GetChaptersBy(versionId string, bookId string) ([]BibleChapter, error) {
	fmt.Printf("Getting chapters for bible version %s book %s\n", versionId, bookId)
	url := fmt.Sprintf("https://api.github.com/repos/wldeh/bible-api/contents/bibles/%s/books/%s/chapters", versionId, bookId)
	var data any
	if err := bs.fetchJsonDataFrom(url, &data); err != nil {
		return []BibleChapter{}, fmt.Errorf("Error fetching %s book %s chapters: %v\n", versionId, bookId, err)
	}

	var chapters []BibleChapter
	for _, c := range data.([]interface{}) {
		if item, ok := c.(map[string]interface{}); ok {
			chapterName := item["name"].(string)
			if !strings.Contains(chapterName, ".json") {
				chapter := BibleChapter{}
				chapter.Id = item["name"].(string)
				chapters = append(chapters, chapter)
			}
		}
	}

	return chapters, nil
}

func (bs WldehBibleService) GetVersesBy(versionId string, bookId string, chapterId string) ([]BibleVerse, error) {
	fmt.Printf("Getting verses for bible version %s book %s chapter %s\n", versionId, bookId, chapterId)
	url := fmt.Sprintf("https://raw.githubusercontent.com/wldeh/bible-api/main/bibles/%s/books/%s/chapters/%s.json", versionId, bookId, chapterId)
	var data any
	if err := bs.fetchJsonDataFrom(url, &data); err != nil {
		return []BibleVerse{}, fmt.Errorf("Error fetching bible version %s book %s chapter %s verses: %v\n", versionId, bookId, chapterId, err)
	}

	var verses []BibleVerse
	if dataMap, ok := data.(map[string]interface{}); ok {
		for _, v := range dataMap["data"].([]interface{}) {
			if item, ok := v.(map[string]interface{}); ok {
				verse := BibleVerse{
					Id:   item["verse"].(string),
					Text: item["text"].(string),
					Book: BibleBook{
						Id: item["book"].(string),
					},
					Chapter: BibleChapter{
						Id: item["chapter"].(string),
					},
				}
				verses = append(verses, verse)
			}
		}
	}

	return verses, nil
}

func (bs WldehBibleService) GetVerseBy(versionId string, bookId string, chapterId string, verseId string) (BibleVerse, error) {
	fmt.Printf("Getting verse for bible version %s book %s chapter %s verse %s\n", versionId, bookId, chapterId, verseId)
	url := fmt.Sprintf("https://raw.githubusercontent.com/wldeh/bible-api/main/bibles/%s/books/%s/chapters/%s/verses/%s.json", versionId, bookId, chapterId, verseId)
	var data any
	if err := bs.fetchJsonDataFrom(url, &data); err != nil {
		return BibleVerse{}, fmt.Errorf("Error fetching bible version %s book %s chapter %s verse %s: %v\n", versionId, bookId, chapterId, verseId, err)
	}

	verse := BibleVerse{}
	if item, ok := data.(map[string]interface{}); ok {
		verse.Id = item["verse"].(string)
		verse.Text = item["text"].(string)
		verse.Book = BibleBook{
			Id: item["book"].(string),
		}
		verse.Chapter = BibleChapter{
			Id: item["chapter"].(string),
		}
	}

	return verse, nil
}

func (bs WldehBibleService) GetVerseOfTheDay() BibleVerse {
	if bs.verseOfTheDay.Id != "" {
		return bs.verseOfTheDay
	}
	book := bs.getRandomBookBy(bs.currentVersion.Id)
	chapter := bs.getRandomChapterBy(bs.currentVersion.Id, book.Id)
	bs.verseOfTheDay = bs.getRandomVerseBy(bs.currentVersion.Id, book.Id, chapter.Id)
	return bs.verseOfTheDay
}

// endregion

// region Wldeh Bible Service Utility Methods
func (bs WldehBibleService) getRandomBookBy(versionId string) BibleBook {
	books, err := bs.GetBooksBy(versionId)
	if err != nil {
		fmt.Println(err)
		return BibleBook{}
	}
	if len(books) == 0 {
		fmt.Println("No books found")
		return BibleBook{}
	}

	return books[rand.IntN(len(books))]
}

func (bs WldehBibleService) getRandomChapterBy(versionId string, bookId string) BibleChapter {
	chapters, err := bs.GetChaptersBy(versionId, bookId)
	if err != nil {
		fmt.Println(err)
		return BibleChapter{}
	}
	if len(chapters) == 0 {
		fmt.Println("No chapters found")
		return BibleChapter{}
	}

	return chapters[rand.IntN(len(chapters))]
}

func (bs WldehBibleService) getRandomVerseBy(versionId string, bookId string, chapterId string) BibleVerse {
	verses, err := bs.GetVersesBy(versionId, bookId, chapterId)
	if err != nil {
		fmt.Println(err)
		return BibleVerse{}
	}
	if len(verses) == 0 {
		fmt.Println("No verses found")
		return BibleVerse{}
	}

	return verses[rand.IntN(len(verses))]
}

func (bs WldehBibleService) fetchJsonDataFrom(url string, outputParam *any) error {
	resp, err := httpGet(url)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, outputParam)
	if err != nil {
		return err
	}

	return nil
}

func (bs WldehBibleService) convertToBibleVersions(data interface{}) []BibleVersion {
	// fmt.Println("Converting to bible version collection...")
	var versions []BibleVersion
	for _, v := range data.([]interface{}) {
		version, err := bs.convertJsonToBibleVersion(v)
		if err != nil {
			continue
		}
		versions = append(versions, version)
	}
	return versions
}

func (bs WldehBibleService) convertToBibleVersion(data interface{}) (BibleVersion, error) {
	// fmt.Println("Converting to bible version...")
	var version BibleVersion
	var err error
	version, err = bs.convertJsonToBibleVersion(data)
	if err != nil {
		return BibleVersion{}, fmt.Errorf("Error converting bible version map: %v\n", err)
	}
	return version, nil
}

func (bs WldehBibleService) convertJsonToBibleVersion(jsonData interface{}) (BibleVersion, error) {
	if bibleVersionMap, versionMapOk := jsonData.(map[string]interface{}); versionMapOk {
		bibleLanguageMap, bibleLanguageMapOk := bibleVersionMap["language"].(map[string]interface{})
		if !bibleLanguageMapOk {
			return BibleVersion{}, fmt.Errorf("Error asserting bible versionlanguage data: %T\n", bibleVersionMap["language"])
		}
		return BibleVersion{
			Id:       bibleVersionMap["id"].(string),
			Name:     bibleVersionMap["localVersionName"].(string),
			Version:  bibleVersionMap["localVersionAbbreviation"].(string),
			Language: bibleLanguageMap["name"].(string),
		}, nil
	}
	return BibleVersion{}, fmt.Errorf("Error asserting bible version data: %T\n", jsonData)
}

// endregion

func NewBibleService(bibleVersionId string) (BibleService, error) {
	bs := WldehBibleService{}

	version, err := bs.GetBibleVersionById(bibleVersionId)
	if err != nil {
		return nil, err
	}

	if version.Id == "" {
		return nil, fmt.Errorf("Bible version not found: %s", bibleVersionId)
	}

	bs.currentVersion = version
	return bs, nil
}
