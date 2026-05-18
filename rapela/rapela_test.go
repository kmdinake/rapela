package rapela

import "testing"

func resetBibleService() {
	bibleService = nil
}

func TestNewBibleService(t *testing.T) {
	t.Cleanup(resetBibleService)

	t.Run("initializes the default bible service", func(t *testing.T) {
		bs, err := NewBibleService()
		if err != nil {
			t.Fatalf("NewBibleService() returned error: %v", err)
		}
		if bs == nil {
			t.Fatal("NewBibleService() returned nil service")
		}

		version := bs.GetBibleVersion()
		if got, want := version.Id, "en-kjv"; got != want {
			t.Errorf("GetBibleVersion().Id = %q, want %q", got, want)
		}
		if got, want := version.Name, "King James Version of the Holy Bible"; got != want {
			t.Errorf("GetBibleVersion().Name = %q, want %q", got, want)
		}
		if got, want := version.Version, "KJV"; got != want {
			t.Errorf("GetBibleVersion().Version = %q, want %q", got, want)
		}
	})

	t.Run("returns the same singleton instance", func(t *testing.T) {
		first, err := NewBibleService()
		if err != nil {
			t.Fatalf("NewBibleService() returned error: %v", err)
		}
		second, err := NewBibleService()
		if err != nil {
			t.Fatalf("second NewBibleService() returned error: %v", err)
		}
		if first != second {
			t.Fatal("NewBibleService() should return the same singleton instance")
		}
	})
}
