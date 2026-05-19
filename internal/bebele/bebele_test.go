package bebele

import (
	"errors"
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

	bs := &BibleService{}
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

	bs := &BibleService{}
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
	if got := bs.GetBibleVersion().Id; got != "en-kjv" {
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

	bs := &BibleService{}
	_, err := bs.GetBibleVersionById("en-kjv")
	if err == nil {
		t.Fatal("expected GetBibleVersionById to return error for invalid JSON")
	}
}

func TestConvertToBibleVersionInvalidType(t *testing.T) {
	bs := &BibleService{}
	_, err := bs.convertToBibleVersion([]byte(`"notobject"`))
	if err == nil {
		t.Fatal("expected convertToBibleVersion to return error for non-object JSON")
	}
}

func TestConvertJsonToBibleVersionInvalidLanguage(t *testing.T) {
	bs := &BibleService{}
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
	verse := BibleVerse{Book: "John", Chapter: "3", Verse: "16", Text: "For God so loved the world..."}
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

	bs := &BibleService{}
	versions := bs.GetBibleVersions()
	if len(versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(versions))
	}
}

func TestGetVerseOfTheDayInvalidJSON(t *testing.T) {
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

	bs := &BibleService{currentVersion: BibleVersion{Id: "en-kjv"}}
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

	bs := &BibleService{}
	versions := bs.GetBibleVersions()
	if len(versions) != 0 {
		t.Fatalf("expected no versions when httpGet fails, got %d", len(versions))
	}
}

func TestSetBibleVersionNoopWhenSameId(t *testing.T) {
	current := BibleVersion{Id: "en-kjv"}
	bs := &BibleService{currentVersion: current}
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

	bs := &BibleService{currentVersion: BibleVersion{Id: "en-kjv"}}
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

	bs := &BibleService{}
	version, err := bs.GetBibleVersionById("es-rvr1960")
	if err != nil {
		t.Fatalf("GetBibleVersionById returned error: %v", err)
	}
	if version.Id != "es-rvr1960" {
		t.Fatalf("expected version id es-rvr1960, got %q", version.Id)
	}
}
