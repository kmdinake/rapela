package bebele

import (
	"io"
	"net/http"
	"net/http/httptest"
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

	bs := &BibleService{}
	versions := bs.convertToBibleVersions([]byte(rawJSON))

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
	rawJSON := `
		{
			"id": "en-kjv",
			"localVersionName": "King James Version of the Holy Bible",
			"localVersionAbbreviation": "KJV",
			"language": {"name": "English"}
		}
	`

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

	bs := &BibleService{}
	version, err := bs.GetBibleVersionById("en-kjv")
	if err != nil {
		t.Fatalf("GetBibleVersionById returned error: %v", err)
	}
	if want := "en-kjv"; version.Id != want {
		t.Errorf("Id = %q, want %q", version.Id, want)
	}
}

func TestGetVerseOfTheDay(t *testing.T) {
	verseJSON := `{
		"id": "john-3-16",
		"book": "John",
		"chapter": 3,
		"verse": "16",
		"text": "For God so loved the world..."
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(verseJSON))
	}))
	defer ts.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL)
	}
	defer func() { httpGet = originalGet }()

	bs := &BibleService{currentVersion: BibleVersion{Id: "en-kjv"}}
	verse := bs.GetVerseOfTheDay()

	if want := "John"; verse.Book != want {
		t.Errorf("Book = %q, want %q", verse.Book, want)
	}
	if want := 3; verse.Chapter != want {
		t.Errorf("Chapter = %d, want %d", verse.Chapter, want)
	}
	if want := "16"; verse.Verse != want {
		t.Errorf("Verse = %q, want %q", verse.Verse, want)
	}
	if want := "For God so loved the world..."; verse.Text != want {
		t.Errorf("Text = %q, want %q", verse.Text, want)
	}
}

func TestSetBibleVersion(t *testing.T) {
	bs := &BibleService{currentVersion: BibleVersion{Id: "en-kjv"}}
	newVersion := BibleVersion{Id: "es-rvr1960", Name: "Reina-Valera 1960", Version: "RVR1960"}

	bs.SetBibleVersion(newVersion)
	if got := bs.currentVersion.Id; got != newVersion.Id {
		t.Fatalf("currentVersion.Id = %q, want %q", got, newVersion.Id)
	}

	bs.SetBibleVersion(newVersion)
	if got := bs.currentVersion.Id; got != newVersion.Id {
		t.Fatalf("currentVersion.Id = %q after second set, want %q", got, newVersion.Id)
	}
}

func TestGetBibleVersionsReturnsEmptyOnError(t *testing.T) {
	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return nil, io.ErrUnexpectedEOF
	}
	defer func() { httpGet = originalGet }()

	bs := &BibleService{}
	versions := bs.GetBibleVersions()
	if len(versions) != 0 {
		t.Fatalf("expected no versions on http error, got %d", len(versions))
	}
}
