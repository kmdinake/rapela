package bebele

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
)

var httpGet = http.Get

type BibleService struct {
	verseOfTheDay  BibleVerse
	currentVersion BibleVersion
}

func NewBibleService(bibleVersionId string) (*BibleService, error) {
	bs := &BibleService{}

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

func fetchJsonDataFrom(url string, outputParam *any) error {
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

func (bs *BibleService) getRandomBookBy(bibleVersionId string) BibleBook {
	fmt.Printf("Getting books for bible version %s\n", bibleVersionId)
	url := fmt.Sprintf("https://api.github.com/repos/wldeh/bible-api/contents/bibles/%s/books", bibleVersionId)
	var data any
	if err := fetchJsonDataFrom(url, &data); err != nil {
		fmt.Printf("Error fetching books: %v\n", err)
		return BibleBook{}
	}

	var books []string
	for _, b := range data.([]interface{}) {
		if item, ok := b.(map[string]interface{}); ok {
			bookName := item["name"].(string)
			books = append(books, bookName)
		}
	}
	if len(books) == 0 {
		fmt.Println("Error parsing books")
		return BibleBook{}
	}

	bookIndex := rand.IntN(len(books))
	book := BibleBook{}
	book.Name = books[bookIndex]
	return book
}

func (bs *BibleService) getRandomChapterBy(bibleVersionId string, bookName string) BibleChapter {
	fmt.Printf("Getting chapters for bible version %s book %s\n", bibleVersionId, bookName)
	url := fmt.Sprintf("https://api.github.com/repos/wldeh/bible-api/contents/bibles/%s/books/%s/chapters", bibleVersionId, bookName)
	var data any
	if err := fetchJsonDataFrom(url, &data); err != nil {
		fmt.Printf("Error fetching chapters: %v\n", err)
		return BibleChapter{}
	}

	var chapters []string
	for _, c := range data.([]interface{}) {
		if item, ok := c.(map[string]interface{}); ok {
			chapterName := item["name"].(string)
			if !strings.Contains(chapterName, ".json") {
				chapters = append(chapters, chapterName)
			}
		}
	}
	if len(chapters) == 0 {
		fmt.Println("Error parsing chapters")
		return BibleChapter{}
	}
	chapterIndex := rand.IntN(len(chapters))
	chapter := BibleChapter{}
	chapter.Name = chapters[chapterIndex]
	return chapter
}

func (bs *BibleService) getRandomVerseBy(bibleVersionId string, bookName string, chapterName string) BibleVerse {
	fmt.Printf("Getting verses for bible version %s book %s chapter %s\n", bibleVersionId, bookName, chapterName)
	url := fmt.Sprintf("https://raw.githubusercontent.com/wldeh/bible-api/main/bibles/%s/books/%s/chapters/%s.json", bibleVersionId, bookName, chapterName)
	var data any
	if err := fetchJsonDataFrom(url, &data); err != nil {
		fmt.Printf("Error fetching books: %v\n", err)
		return BibleVerse{}
	}

	var verses []BibleVerse
	if dataMap, ok := data.(map[string]interface{}); ok {
		for _, v := range dataMap["data"].([]interface{}) {
			if item, ok := v.(map[string]interface{}); ok {
				verse := BibleVerse{
					Book:    item["book"].(string),
					Chapter: item["chapter"].(string),
					Verse:   item["verse"].(string),
					Text:    item["text"].(string),
				}
				verse.Id = fmt.Sprintf("%s-%s-%s", verse.Book, verse.Chapter, verse.Verse)
				verses = append(verses, verse)
			}
		}
	}
	if len(verses) == 0 {
		fmt.Println("Error parsing verses")
		return BibleVerse{}
	}
	verseIndex := rand.IntN(len(verses))
	return verses[verseIndex]
}

// Get a random verse from the bible
func (bs *BibleService) GetVerseOfTheDay() BibleVerse { // Notice that we are using a pointer receiver here, which allows us to modify the state of the BibleService struct
	if bs.verseOfTheDay.Id != "" {
		return bs.verseOfTheDay
	}
	book := bs.getRandomBookBy(bs.currentVersion.Id)
	chapter := bs.getRandomChapterBy(bs.currentVersion.Id, book.Name)
	bs.verseOfTheDay = bs.getRandomVerseBy(bs.currentVersion.Id, book.Name, chapter.Name)
	return bs.verseOfTheDay
}

// Get a list of available Bible versions
func (bs *BibleService) GetBibleVersions() []BibleVersion {
	fmt.Println("Getting bible versions...")
	url := "https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/bibles.json"
	var data any
	if err := fetchJsonDataFrom(url, &data); err != nil {
		fmt.Printf("Error fetching bible versions: %v\n", err)
		return []BibleVersion{}
	}

	return bs.convertToBibleVersions(data)
}

// Get a bible version by its ID
func (bs *BibleService) GetBibleVersionById(bibleVersionId string) (BibleVersion, error) {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/%s/%s.json", bibleVersionId, bibleVersionId)
	var data any
	if err := fetchJsonDataFrom(url, &data); err != nil {
		return BibleVersion{}, fmt.Errorf("Error fetching bible versions: %v\n", err)
	}

	version, err := bs.convertToBibleVersion(data)
	if err != nil {
		return BibleVersion{}, fmt.Errorf("Error converting bible version JSON: %v\n", err)
	}

	if version.Id != bibleVersionId {
		return BibleVersion{}, fmt.Errorf("Bible version not found: %s\n", bibleVersionId)
	}

	return version, nil
}

// Set the current Bible version to be used
func (bs *BibleService) SetBibleVersion(version BibleVersion) {
	if bs.currentVersion.Id != version.Id {
		bs.currentVersion = version
	}
}

// Get the current Bible version being used
func (bs *BibleService) GetBibleVersion() BibleVersion {
	return bs.currentVersion
}

// Get a specific verse by its reference (book, chapter, verse)
func (bs *BibleService) GetVerseByReference(book string, chapter int, verse int) BibleVerse {
	panic("Method not implemented")
}

func (bs *BibleService) convertToBibleVersions(data interface{}) []BibleVersion {
	fmt.Println("Converting to bible version collection...")
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

func (bs *BibleService) convertToBibleVersion(data interface{}) (BibleVersion, error) {
	fmt.Println("Converting to bible version...")
	var version BibleVersion
	var err error
	version, err = bs.convertJsonToBibleVersion(data)
	if err != nil {
		return BibleVersion{}, fmt.Errorf("Error converting bible version map: %v\n", err)
	}
	return version, nil
}

func (bs *BibleService) convertJsonToBibleVersion(jsonData interface{}) (BibleVersion, error) {
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

type BibleVersion struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"`
}

func (v BibleVersion) String() string {
	return fmt.Sprintf("- %s: %s (%s)", v.Id, v.Name, v.Version)
}

type BibleVerse struct {
	Id      string `json:"id"`
	Book    string `json:"book"`
	Chapter string `json:"chapter"`
	Verse   string `json:"verse"`
	Text    string `json:"text"`
}

func (v BibleVerse) String() string {
	return fmt.Sprintf("Verse of the day: %s %s:%s - %s", v.Book, v.Chapter, v.Verse, v.Text)
}

type BibleBook struct {
	Name string `json:"name"`
}

type BibleChapter struct {
	Name string `json:"name"`
}
