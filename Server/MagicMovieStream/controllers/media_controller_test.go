package controllers

import "testing"

func TestIsVideoExt(t *testing.T) {
	cases := map[string]bool{
		".mp4":  true,
		".MP4":  true,
		".webm": true,
		".mov":  true,
		".mkv":  true,
		".txt":  false,
		".pdf":  false,
		"":      false,
		".mp3":  false,
	}

	for ext, want := range cases {
		if got := isVideoExt(ext); got != want {
			t.Errorf("isVideoExt(%q) = %v, want %v", ext, got, want)
		}
	}
}
