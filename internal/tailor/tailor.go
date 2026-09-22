package tailor

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/moulindavid/cv-builder/internal/resume"
)

type Match struct {
	Keyword string
	Present bool
}
type Report struct {
	Score            int
	Matches          []Match
	Missing          []string
	RankedSkills     []resume.Skill
	SuggestedSummary string
}

func ApplyTarget(r resume.Resume, t resume.Target, lang string) resume.Resume {
	weights := weights(t.Priorities)
	sort.SliceStable(r.Skills, func(i, j int) bool { return scoreSkill(r.Skills[i], weights) > scoreSkill(r.Skills[j], weights) })
	for i := range r.Experience {
		sort.SliceStable(r.Experience[i].Bullets, func(a, b int) bool {
			return scoreTags(r.Experience[i].Bullets[a].Tags, weights) > scoreTags(r.Experience[i].Bullets[b].Tags, weights)
		})
		if t.MaxBullets > 0 && len(r.Experience[i].Bullets) > t.MaxBullets {
			r.Experience[i].Bullets = r.Experience[i].Bullets[:t.MaxBullets]
		}
	}
	if title := t.Titles.Get(lang); title != "" {
		r.Profile.Titles[lang] = title
	}
	return r
}
func Analyze(r resume.Resume, job, lang string) Report {
	known := map[string]resume.Skill{}
	for _, skill := range r.Skills {
		known[normalize(skill.Name)] = skill
		known[normalize(skill.ID)] = skill
	}
	jobText := normalize(job)
	seen := map[string]bool{}
	var matches []Match
	var found []resume.Skill
	for key, skill := range known {
		if !seen[skill.ID] && containsKeyword(jobText, key) {
			seen[skill.ID] = true
			matches = append(matches, Match{Keyword: skill.Name, Present: true})
			found = append(found, skill)
		}
	}
	var missing []string
	for _, keyword := range externalKeywords {
		key := normalize(keyword)
		if containsKeyword(jobText, key) {
			if _, exists := known[key]; !exists {
				missing = append(missing, keyword)
			}
		}
	}
	sort.Strings(missing)
	score := 0
	if total := len(found) + len(missing); total > 0 {
		score = 100 * len(found) / total
	}
	return Report{Score: score, Matches: matches, Missing: missing, RankedSkills: found, SuggestedSummary: r.Summaries.Get(lang)}
}

var externalKeywords = []string{
	"AWS", "Azure", "C#", "Kafka", "Kubernetes", "Node.js", "Python", "Terraform",
}

func containsKeyword(text, keyword string) bool {
	if keyword == "" {
		return false
	}
	for _, token := range tokenize(text) {
		if normalize(token) == keyword {
			return true
		}
	}
	return false
}

func Format(rep Report) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Match score: %d%%\n", rep.Score)
	for _, m := range rep.Matches {
		fmt.Fprintf(&b, "Matched skill: %s\n", m.Keyword)
	}
	for _, keyword := range rep.Missing {
		fmt.Fprintf(&b, "Missing keyword: %s\n", keyword)
	}
	fmt.Fprintf(&b, "Suggested summary (source-safe): %s\n", rep.SuggestedSummary)
	return b.String()
}
func weights(p []string) map[string]int {
	m := map[string]int{}
	for i, v := range p {
		m[normalize(v)] = len(p) - i
	}
	return m
}
func scoreSkill(s resume.Skill, w map[string]int) int {
	return w[normalize(s.ID)] + w[normalize(s.Name)] + scoreTags(s.Tags, w)
}
func scoreTags(tags []string, w map[string]int) int {
	n := 0
	for _, t := range tags {
		n += w[normalize(t)]
	}
	return n
}
func normalize(s string) string {
	return strings.ToLower(strings.Trim(strings.TrimSpace(s), ".,;:!?()[]{}"))
}
func tokenize(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '.' })
}
