package bebele

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConvertToBibleVersions(t *testing.T) {
	rawJSON := `[
        {
            "id": "en-kjv",
            "localVersionName": "King James Version of the Holy Bible",
            "localVersionAbbreviation": "KJV",
            "language": {"name": "English"}
        }
    ]`

	bs := WldehBibleService{}
	var data any
	json.Unmarshal([]byte(rawJSON), &data)
	versions := bs.convertToBibleVersions(data)

	if len(versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(versions))
	}

	got := versions[0]
	if want := "en-kjv"; got.Id != want {
		t.Errorf("Id = %q, want %q", got.Id, want)
	}
	if want := "King James Version of the Holy Bible"; got.Name != want {
		t.Errorf("Name = %q, want %q", got.Name, want)
	}
	if want := "KJV"; got.Version != want {
		t.Errorf("Version = %q, want %q", got.Version, want)
	}
	if want := "English"; got.Language != want {
		t.Errorf("Language = %q, want %q", got.Language, want)
	}
}

func TestGetBibleVersionById(t *testing.T) {
	rawJSON := `{
		"id": "en-kjv",
		"localVersionName": "King James Version of the Holy Bible",
		"localVersionAbbreviation": "KJV",
		"language": {"name": "English"}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(rawJSON))
	}))
	defer ts.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL)
	}
	defer func() { httpGet = originalGet }()

	bs := WldehBibleService{}
	version, err := bs.GetBibleVersionById("en-kjv")
	if err != nil {
		t.Fatalf("GetBibleVersionById returned error: %v", err)
	}
	if want := "en-kjv"; version.Id != want {
		t.Errorf("Id = %q, want %q", version.Id, want)
	}
}

func TestGetBibleVersionByIdNotFound(t *testing.T) {
	rawJSON := `{
		"id": "en-kjv",
		"localVersionName": "King James Version of the Holy Bible",
		"localVersionAbbreviation": "KJV",
		"language": {"name": "English"}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(rawJSON))
	}))
	defer ts.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL)
	}
	defer func() { httpGet = originalGet }()

	bs := WldehBibleService{}
	_, err := bs.GetBibleVersionById("unknown")
	if err == nil {
		t.Fatal("expected error for unknown bible version id")
	}
}

func TestNewBibleServiceSuccess(t *testing.T) {
	rawJSON := `{
        "id": "en-kjv",
        "localVersionName": "King James Version of the Holy Bible",
        "localVersionAbbreviation": "KJV",
        "language": {"name": "English"}
    }`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(rawJSON))
	}))
	defer ts.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL)
	}
	defer func() { httpGet = originalGet }()

	bs, err := NewBibleService("en-kjv")
	if err != nil {
		t.Fatalf("NewBibleService returned error: %v", err)
	}
	if got, _ := bs.GetBibleVersion(); got.Id != "en-kjv" {
		t.Fatalf("GetBibleVersion().Id = %q, want %q", got, "en-kjv")
	}
}

func TestNewBibleServiceHttpError(t *testing.T) {
	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return nil, errors.New("network failure")
	}
	defer func() { httpGet = originalGet }()

	_, err := NewBibleService("en-kjv")
	if err == nil {
		t.Fatal("expected NewBibleService to return error when httpGet fails")
	}
}

func TestGetBibleVersionByIdInvalidJSON(t *testing.T) {
	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bad": "verse"}`))
	}))
	defer rs.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(rs.URL)
	}
	defer func() { httpGet = originalGet }()

	bs := WldehBibleService{}
	_, err := bs.GetBibleVersionById("en-kjv")
	if err == nil {
		t.Fatal("expected GetBibleVersionById to return error for invalid JSON")
	}
}

func TestConvertToBibleVersionInvalidType(t *testing.T) {
	bs := WldehBibleService{}
	_, err := bs.convertToBibleVersion([]byte(`"notobject"`))
	if err == nil {
		t.Fatal("expected convertToBibleVersion to return error for non-object JSON")
	}
}

func TestConvertJsonToBibleVersionInvalidLanguage(t *testing.T) {
	bs := WldehBibleService{}
	_, err := bs.convertJsonToBibleVersion(map[string]interface{}{
		"id":                       "en-kjv",
		"localVersionName":         "King James Version of the Holy Bible",
		"localVersionAbbreviation": "KJV",
		"language":                 "English",
	})
	if err == nil {
		t.Fatal("expected convertJsonToBibleVersion to return error when language is not a map")
	}
}

func TestBibleVersionString(t *testing.T) {
	version := BibleVersion{Id: "en-kjv", Name: "King James Version of the Holy Bible", Version: "KJV"}
	if got := version.String(); !strings.Contains(got, "en-kjv") || !strings.Contains(got, "KJV") {
		t.Fatalf("unexpected BibleVersion.String output: %q", got)
	}
}

func TestBibleVerseString(t *testing.T) {
	verse := BibleVerse{Book: BibleBook{Id: "John"}, Chapter: BibleChapter{Id: "3"}, Id: "16", Text: "For God so loved the world..."}
	if got := verse.String(); !strings.Contains(got, "John 3:16") || !strings.Contains(got, "For God so loved the world...") {
		t.Fatalf("unexpected BibleVerse.String output: %q", got)
	}
}

func TestGetBibleVersionsSuccess(t *testing.T) {
	rawJSON := `[
        {
            "id": "en-kjv",
            "localVersionName": "King James Version of the Holy Bible",
            "localVersionAbbreviation": "KJV",
            "language": {"name": "English"}
        }
    ]`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(rawJSON))
	}))
	defer ts.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL)
	}
	defer func() { httpGet = originalGet }()

	bs := WldehBibleService{}
	versions, err := bs.GetBibleVersions()
	if err != nil {
		t.Fatalf("GetBibleVersions() returned error: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(versions))
	}
}

func TestGetVerseOfTheDayInvalidJSON(t *testing.T) {
	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if parts[len(parts)-1] == "chapters" {
			_, _ = w.Write([]byte(`[{
				"name": "1",
				"path": "bibles/en-kjv/books/1chronicles/chapters/1",
				"sha": "97f98e3f49951eaf7690dc21af0d021c9dfbc67a",
				"size": 0,
				"url": "https://api.github.com/repos/wldeh/bible-api/contents/bibles/en-kjv/books/1chronicles/chapters/1?ref=main",
				"html_url": "https://github.com/wldeh/bible-api/tree/main/bibles/en-kjv/books/1chronicles/chapters/1",
				"git_url": "https://api.github.com/repos/wldeh/bible-api/git/trees/97f98e3f49951eaf7690dc21af0d021c9dfbc67a",
				"download_url": null,
				"type": "dir",
				"_links": {
				"self": "https://api.github.com/repos/wldeh/bible-api/contents/bibles/en-kjv/books/1chronicles/chapters/1?ref=main",
				"git": "https://api.github.com/repos/wldeh/bible-api/git/trees/97f98e3f49951eaf7690dc21af0d021c9dfbc67a",
				"html": "https://github.com/wldeh/bible-api/tree/main/bibles/en-kjv/books/1chronicles/chapters/1"
				}
			}]`))
		} else if parts[len(parts)-1] == "books" {
			_, _ = w.Write([]byte(`[{
				"name": "1chronicles",
				"path": "bibles/en-kjv/books/1chronicles",
				"sha": "30e45bb08b963605deff20db6efba588d58421b0",
				"size": 0,
				"url": "https://api.github.com/repos/wldeh/bible-api/contents/bibles/en-kjv/books/1chronicles?ref=main",
				"html_url": "https://github.com/wldeh/bible-api/tree/main/bibles/en-kjv/books/1chronicles",
				"git_url": "https://api.github.com/repos/wldeh/bible-api/git/trees/30e45bb08b963605deff20db6efba588d58421b0",
				"download_url": null,
				"type": "dir",
				"_links": {
				"self": "https://api.github.com/repos/wldeh/bible-api/contents/bibles/en-kjv/books/1chronicles?ref=main",
				"git": "https://api.github.com/repos/wldeh/bible-api/git/trees/30e45bb08b963605deff20db6efba588d58421b0",
				"html": "https://github.com/wldeh/bible-api/tree/main/bibles/en-kjv/books/1chronicles"
				}
			}]`))
		} else {
			_, _ = w.Write([]byte(`[{"bad": "verse"}]`))
		}
	}))
	defer rs.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		parts := strings.Split(url, "/")
		lastPart := parts[len(parts)-1]
		newUrl := fmt.Sprintf("%s/%s", rs.URL, lastPart)
		return http.Get(newUrl)
	}
	defer func() { httpGet = originalGet }()

	bs := &WldehBibleService{currentVersion: BibleVersion{Id: "en-kjv"}}
	verse := bs.GetVerseOfTheDay()
	if verse.Text != "" {
		t.Fatalf("expected empty verse on invalid JSON, got %q", verse.Text)
	}
}

func TestGetBibleVersionsWithHttpError(t *testing.T) {
	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return nil, errors.New("timeout")
	}
	defer func() { httpGet = originalGet }()

	bs := WldehBibleService{}
	versions, err := bs.GetBibleVersions()
	if err != nil {
		t.Fatalf("GetBibleVersions() returned error: %v", err)
	}
	if len(versions) != 0 {
		t.Fatalf("expected no versions when httpGet fails, got %d", len(versions))
	}
}

func TestSetBibleVersionNoopWhenSameId(t *testing.T) {
	current := BibleVersion{Id: "en-kjv"}
	bs := &WldehBibleService{currentVersion: current}
	bs.SetBibleVersion(current)
	if bs.currentVersion.Id != current.Id {
		t.Fatalf("expected unchanged version id %q, got %q", current.Id, bs.currentVersion.Id)
	}
}

func TestGetVerseOfTheDayHttpError(t *testing.T) {
	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return nil, errors.New("network down")
	}
	defer func() { httpGet = originalGet }()

	bs := &WldehBibleService{currentVersion: BibleVersion{Id: "en-kjv"}}
	verse := bs.GetVerseOfTheDay()
	if verse.Id != "" {
		t.Fatalf("expected empty verse on http error, got %q", verse.Id)
	}
}

func TestGetBibleVersionByIdFromVersions(t *testing.T) {
	rawJSON := `{
		"id": "es-rvr1960",
		"localVersionName": "Reina-Valera 1960",
		"localVersionAbbreviation": "RVR1960",
		"language": {"name": "Spanish"}
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(rawJSON))
	}))
	defer ts.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL)
	}
	defer func() { httpGet = originalGet }()

	bs := WldehBibleService{}
	version, err := bs.GetBibleVersionById("es-rvr1960")
	if err != nil {
		t.Fatalf("GetBibleVersionById returned error: %v", err)
	}
	if version.Id != "es-rvr1960" {
		t.Fatalf("expected version id es-rvr1960, got %q", version.Id)
	}
}
