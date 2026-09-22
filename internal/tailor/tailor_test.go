package tailor

import (
	"testing"

	"github.com/moulindavid/cv-builder/internal/resume"
)

func TestApplyTargetRanksOnlyExistingSkills(t *testing.T) {
	r := resume.Resume{
		Profile: resume.Profile{Titles: resume.Localized{"en": "Old"}},
		Skills:  []resume.Skill{{ID: "go", Name: "Go"}, {ID: "java", Name: "Java"}},
		Experience: []resume.Experience{{Bullets: []resume.Bullet{
			{ID: "go", Tags: []string{"go"}}, {ID: "java", Tags: []string{"java"}},
		}}},
	}
	got := ApplyTarget(r, resume.Target{Priorities: []string{"java", "kubernetes"}, MaxBullets: 1}, "en")
	if got.Skills[0].ID != "java" {
		t.Fatalf("first=%s", got.Skills[0].ID)
	}
	if len(got.Skills) != 2 {
		t.Fatal("target introduced a skill")
	}
	if got.Experience[0].Bullets[0].ID != "java" {
		t.Fatal("bullet was not ranked")
	}
}
func TestAnalyzeDoesNotInvent(t *testing.T) {
	r := resume.Resume{Skills: []resume.Skill{{ID: "java", Name: "Java"}}, Summaries: resume.Localized{"en": "Backend engineer"}}
	rep := Analyze(r, "Java Kubernetes.", "en")
	if len(rep.RankedSkills) != 1 || rep.RankedSkills[0].Name != "Java" || len(rep.Missing) != 1 || rep.Missing[0] != "Kubernetes" {
		t.Fatalf("skills=%v", rep.RankedSkills)
	}
}
