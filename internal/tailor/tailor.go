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
	Priorities       []string
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
	jobText := normalizeText(job)
	found := make([]resume.Skill, 0)
	matches := make([]Match, 0)
	knownTerms := make(map[string]bool)
	priorities := make([]string, 0)
	seenPriorities := make(map[string]bool)
	addPriority := func(value string) {
		key := normalizeText(value)
		if key != "" && !seenPriorities[key] {
			priorities = append(priorities, value)
			seenPriorities[key] = true
		}
	}

	for _, skill := range r.Skills {
		terms := append([]string{skill.Name, skill.ID}, skill.Aliases...)
		matched := false
		for _, term := range terms {
			normalized := normalizeText(term)
			knownTerms[normalized] = true
			if containsKeyword(jobText, normalized) {
				matched = true
			}
		}
		if matched {
			matches = append(matches, Match{Keyword: skill.Name, Present: true})
			found = append(found, skill)
			addPriority(skill.ID)
		}
		for _, tag := range skill.Tags {
			if containsKeyword(jobText, normalizeText(tag)) {
				addPriority(tag)
			}
		}
	}
	for _, experience := range r.Experience {
		for _, bullet := range experience.Bullets {
			for _, tag := range bullet.Tags {
				if containsKeyword(jobText, normalizeText(tag)) {
					addPriority(tag)
				}
			}
		}
	}

	missing := make([]string, 0)
	for _, keyword := range r.Tailoring.ExternalKeywords {
		key := normalizeText(keyword)
		if containsKeyword(jobText, key) && !knownTerms[key] {
			missing = append(missing, keyword)
		}
	}
	sort.Strings(missing)

	score := 0
	if total := len(found) + len(missing); total > 0 {
		score = 100 * len(found) / total
	}
	return Report{
		Score: score, Matches: matches, Missing: missing, RankedSkills: found,
		Priorities: priorities, SuggestedSummary: r.Summaries.Get(lang),
	}
}

func containsKeyword(text, keyword string) bool {
	if keyword == "" {
		return false
	}
	return strings.Contains(" "+text+" ", " "+keyword+" ")
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
		m[normalizeText(v)] = len(p) - i
	}
	return m
}

func scoreSkill(s resume.Skill, w map[string]int) int {
	score := max(w[normalizeText(s.ID)], w[normalizeText(s.Name)])
	for _, tag := range s.Tags {
		score = max(score, w[normalizeText(tag)])
	}
	return score
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func scoreTags(tags []string, w map[string]int) int {
	score := 0
	for _, tag := range tags {
		score = max(score, w[normalizeText(tag)])
	}
	return score
}

func normalizeText(s string) string {
	var b strings.Builder
	space := true
	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '+' || r == '#' || r == '.' {
			b.WriteRune(r)
			space = false
		} else if !space {
			b.WriteByte(' ')
			space = true
		}
	}
	fields := strings.Fields(b.String())
	for i := range fields {
		fields[i] = strings.Trim(fields[i], ".")
	}
	return strings.Join(fields, " ")
}
