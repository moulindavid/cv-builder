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
	r := resume.Resume{Skills: []resume.Skill{{ID: "java", Name: "Java"}}, Tailoring: resume.Tailoring{ExternalKeywords: []string{"Kubernetes"}}, Summaries: resume.Localized{"en": "Backend engineer"}}
	rep := Analyze(r, "Java Kubernetes.", "en")
	if len(rep.RankedSkills) != 1 || rep.RankedSkills[0].Name != "Java" || len(rep.Missing) != 1 || rep.Missing[0] != "Kubernetes" {
		t.Fatalf("skills=%v", rep.RankedSkills)
	}
}

func TestAnalyzeIsDeterministicAndMatchesPhrasesAndTags(t *testing.T) {
	r := resume.Resume{Skills: []resume.Skill{
		{ID: "spring-boot", Name: "Spring Boot", Tags: []string{"microservices"}},
		{ID: "postgresql", Name: "PostgreSQL", Tags: []string{"database"}},
	}, Tailoring: resume.Tailoring{ExternalKeywords: []string{"Kubernetes"}}, Summaries: resume.Localized{"en": "Backend engineer"}}
	job := "Spring Boot microservices using PostgreSQL and Kubernetes"
	first := Format(Analyze(r, job, "en"))
	for i := 0; i < 20; i++ {
		if got := Format(Analyze(r, job, "en")); got != first {
			t.Fatalf("non-deterministic report: %q != %q", first, got)
		}
	}
	rep := Analyze(r, job, "en")
	if len(rep.RankedSkills) != 2 || rep.RankedSkills[0].ID != "spring-boot" || rep.RankedSkills[1].ID != "postgresql" {
		t.Fatalf("ranked skills=%v", rep.RankedSkills)
	}
	if len(rep.Missing) != 1 || rep.Missing[0] != "Kubernetes" {
		t.Fatalf("missing=%v", rep.Missing)
	}
}

func TestAnalyzePrioritiesRankBulletsWithoutInventing(t *testing.T) {
	r := resume.Resume{
		Skills: []resume.Skill{{ID: "java", Name: "Java", Tags: []string{"backend"}}},
		Experience: []resume.Experience{{Bullets: []resume.Bullet{
			{ID: "generic", Tags: []string{"collaboration"}},
			{ID: "backend", Tags: []string{"java", "backend"}},
		}}},
	}
	report := Analyze(r, "Senior Java backend engineer", "en")
	got := ApplyTarget(r, resume.Target{Priorities: report.Priorities}, "en")
	if got.Experience[0].Bullets[0].ID != "backend" {
		t.Fatalf("first bullet=%s, priorities=%v", got.Experience[0].Bullets[0].ID, report.Priorities)
	}
	if len(got.Skills) != 1 || got.Skills[0].ID != "java" {
		t.Fatalf("tailoring changed source skills: %v", got.Skills)
	}
}

func TestAnalyzeUsesExplicitAliases(t *testing.T) {
	r := resume.Resume{
		Skills:    []resume.Skill{{ID: "aws-s3", Name: "AWS S3", Aliases: []string{"AWS"}}},
		Tailoring: resume.Tailoring{ExternalKeywords: []string{"AWS"}},
	}
	rep := Analyze(r, "Experience with AWS", "en")
	if len(rep.RankedSkills) != 1 || rep.RankedSkills[0].ID != "aws-s3" {
		t.Fatalf("matches=%v", rep.RankedSkills)
	}
	if len(rep.Missing) != 0 {
		t.Fatalf("missing=%v", rep.Missing)
	}
}
