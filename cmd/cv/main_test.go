package main

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"David Moulin", "david-moulin"},
		{" Senior Backend / Java ", "senior-backend-java"},
		{"Élodie Durand", "élodie-durand"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := slugify(tt.input); got != tt.want {
				t.Fatalf("slugify(%q)=%q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestOutputPath(t *testing.T) {
	got := outputPath("David Moulin", "senior-backend", "en", ".pdf")
	want := "output/david-moulin-senior-backend-en.pdf"
	if got != want {
		t.Fatalf("outputPath()=%q, want %q", got, want)
	}
}

func TestValidateLanguage(t *testing.T) {
	if err := validateLanguage("en"); err != nil {
		t.Fatal(err)
	}
	if err := validateLanguage("de"); err == nil {
		t.Fatal("expected unsupported language error")
	}
}
