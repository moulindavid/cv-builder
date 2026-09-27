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
		Skills:    []resume.Skill{{ID: "java", Name: "Java", Category: "backend"}},
		Experience: []resume.Experience{{ID: "x", Company: "C", Start: "2020", Descriptions: resume.Localized{"fr": "D", "en": "D"}, Technologies: []string{"java"},
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
		{"unsupported locale", func(r *resume.Resume) { r.Experience[0].Bullets[0].Text["de"] = "Text" }, "unsupported locale"},
		{"missing description locale", func(r *resume.Resume) { delete(r.Experience[0].Descriptions, "fr") }, "descriptions missing fr"},
		{"incoherent dates", func(r *resume.Resume) { r.Experience[0].End = "2019" }, "starts after it ends"},
		{"unknown category", func(r *resume.Resume) { r.Skills[0].Category = "mystery" }, "unknown category"},
		{"duplicate external keyword", func(r *resume.Resume) { r.Tailoring.ExternalKeywords = []string{"Kafka", "kafka"} }, "duplicate tailoring"},
		{"empty external keyword", func(r *resume.Resume) { r.Tailoring.ExternalKeywords = []string{" "} }, "cannot be empty"},
		{"duplicate skill alias", func(r *resume.Resume) { r.Skills[0].Aliases = []string{"JVM", "jvm"} }, "duplicate alias"},
		{"unknown project skill", func(r *resume.Resume) {
			r.Projects = []resume.Project{{ID: "p", Name: "Project", Technologies: []string{"rust"}}}
		}, "project p references unknown skill rust"},
		{"project without name", func(r *resume.Resume) { r.Projects = []resume.Project{{ID: "p"}} }, "project p has no name"},
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

func TestValidateTarget(t *testing.T) {
	r := valid()
	tests := []struct {
		name   string
		target resume.Target
		want   string
	}{
		{"valid", resume.Target{Name: "backend", Priorities: []string{"java"}, MaxBullets: 2}, ""},
		{"unknown priority", resume.Target{Name: "backend", Priorities: []string{"kubernetes"}}, "unknown priority"},
		{"negative max bullets", resume.Target{Name: "backend", MaxBullets: -1}, "cannot be negative"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTarget(tt.target, r)
			if tt.want == "" && err != nil {
				t.Fatal(err)
			}
			if tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Fatalf("error=%v", err)
			}
		})
	}
}
