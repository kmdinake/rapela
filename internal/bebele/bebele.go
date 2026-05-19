package bebele

import (
	"encoding/json"
	"fmt"
	"io"
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

// Get a random verse from the bible
func (bs *BibleService) GetVerseOfTheDay() BibleVerse { // Notice that we are using a pointer receiver here, which allows us to modify the state of the BibleService struct
	book := strings.ToLower("John")
	chapter := 3
	verse := 16

	fmt.Println("Getting the verse...")
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/%s/books/%s/chapters/%d/verses/%d.json", bs.currentVersion.Id, book, chapter, verse)
	resp, err := httpGet(url)
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
	fmt.Println("Getting bible versions...")
	url := "https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/bibles.json"
	resp, err := httpGet(url)
	if err != nil {
		fmt.Printf("Error fetching bible versions: %v\n", err)
		return []BibleVersion{}
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("Error reading bible versions body: %v\n", readErr)
		return []BibleVersion{}
	}

	return bs.convertToBibleVersions(body)
}

// Get a bible version by its ID
func (bs *BibleService) GetBibleVersionById(bibleVersionId string) (BibleVersion, error) {
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/wldeh/bible-api/bibles/%s/%s.json", bibleVersionId, bibleVersionId)
	resp, err := httpGet(url)

	if err != nil {
		return BibleVersion{}, fmt.Errorf("Bible version not found: %s\n%v\n", bibleVersionId, err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BibleVersion{}, fmt.Errorf("Error reading bible version body: %v\n", err)
	}

	version, err := bs.convertToBibleVersion(body)
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

func (bs *BibleService) convertToBibleVersions(rawJsonData []byte) []BibleVersion {
	fmt.Println("Converting to bible version collection...")
	var versions []BibleVersion
	var data interface{}
	jsonErr := json.Unmarshal(rawJsonData, &data)
	if jsonErr != nil {
		fmt.Printf("Error parsing bible versions JSON: %v\n", jsonErr)
		return []BibleVersion{}
	}
	for _, v := range data.([]interface{}) {
		version, err := bs.convertJsonToBibleVersion(v)
		if err != nil {
			continue
		}
		versions = append(versions, version)
	}
	return versions
}

func (bs *BibleService) convertToBibleVersion(rawJsonData []byte) (BibleVersion, error) {
	fmt.Println("Converting to bible version...")
	var version BibleVersion
	var data interface{}
	jsonErr := json.Unmarshal(rawJsonData, &data)
	if jsonErr != nil {
		return BibleVersion{}, fmt.Errorf("Error parsing bible version JSON: %v\n", jsonErr)
	}
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
	Chapter int    `json:"chapter"`
	Verse   string `json:"verse"`
	Text    string `json:"text"`
}

func (v BibleVerse) String() string {
	return fmt.Sprintf("Verse of the day: %s %d:%s - %s", v.Book, v.Chapter, v.Verse, v.Text)
}
