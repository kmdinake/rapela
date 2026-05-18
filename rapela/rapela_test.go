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
	bs := NewDefaultBibleService()
	if bs.CurrentVersion.Name != "King James" {
		t.Errorf("Expected Bible name to be 'King James', got '%s'", bs.CurrentVersion.Name)
	}
	if bs.CurrentVersion.Version != "kjv" {
		t.Errorf("Expected Bible version to be 'kjv', got '%s'", bs.CurrentVersion.Version)
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
