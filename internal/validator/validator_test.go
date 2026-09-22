package validator

import (
	"strings"
	"testing"

	"github.com/moulindavid/cv-builder/internal/resume"
)

func valid() resume.Resume {
	return resume.Resume{
		Profile:   resume.Profile{Name: "D", Titles: resume.Localized{"fr": "T", "en": "T"}},
		Contact:   resume.Contact{Email: "d@example.com"},
		Summaries: resume.Localized{"fr": "S", "en": "S"},
		Skills:    []resume.Skill{{ID: "java", Name: "Java"}},
		Experience: []resume.Experience{{ID: "x", Company: "C", Start: "2020", Technologies: []string{"java"},
			Bullets: []resume.Bullet{{ID: "b", Text: resume.Localized{"fr": "F", "en": "E"}}}}},
	}
}
func TestValidate(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*resume.Resume)
		want   string
	}{
		{"valid", func(*resume.Resume) {}, ""},
		{"unknown skill", func(r *resume.Resume) { r.Experience[0].Technologies = []string{"rust"} }, "unknown skill"},
		{"duplicate", func(r *resume.Resume) { r.Skills = append(r.Skills, r.Skills[0]) }, "duplicate id"},
		{"unsupported locale", func(r *resume.Resume) { r.Experience[0].Bullets[0].Text["truncated fragment"] = "" }, "unsupported locale"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := valid()
			tt.mutate(&r)
			err := Validate(r)
			if tt.want == "" && err != nil {
				t.Fatal(err)
			}
			if tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error=%v", err)
			}
		})
	}
}
