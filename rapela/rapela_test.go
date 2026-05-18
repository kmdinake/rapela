package rapela

import "testing"

func TestBibleServiceCreation(t *testing.T) {
	bs := NewBibleService("New International Version", "niv")
	if bs.CurrentVersion.Name != "New International Version" {
		t.Errorf("Expected Bible name to be 'New International Version', got '%s'", bs.CurrentVersion.Name)
	}
	if bs.CurrentVersion.Version != "niv" {
		t.Errorf("Expected Bible version to be 'niv', got '%s'", bs.CurrentVersion.Version)
	}

	bs2 := NewBibleService("New International Version", "niv")
	if bs != bs2 {
		t.Errorf("Expected singleton instances to be the same")
	}

	bs3 := NewBibleService("World English Bible", "WEB")
	if bs == bs3 {
		t.Errorf("Expected singleton instances to be different")
	}
}

func TestDefaultBibleServiceCreation(t *testing.T) {
	defaultBibleName := "King James Version of the Holy Bible"
	defaultBibleVersion := "KJV"
	bs := NewDefaultBibleService()
	if bs.CurrentVersion.Name != defaultBibleName {
		t.Errorf("Expected Bible name to be '%s', got '%s'", defaultBibleName, bs.CurrentVersion.Name)
	}
	if bs.CurrentVersion.Version != defaultBibleVersion {
		t.Errorf("Expected Bible version to be '%s', got '%s'", defaultBibleVersion, bs.CurrentVersion.Version)
	}

	bs2 := NewDefaultBibleService()
	if bs2 != bs {
		t.Errorf("Expected default singleton instances to be the same")
	}

	bs3 := NewBibleService("New International Version", "niv")
	if bs3 == bs {
		t.Errorf("Expected default singleton instances to be different")
	}
}
