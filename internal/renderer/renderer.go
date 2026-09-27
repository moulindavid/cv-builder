package renderer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/moulindavid/cv-builder/internal/resume"
)

type View struct {
	Resume resume.Resume
	Lang   string
	Labels map[string]string
	Skills []SkillGroup
}

type SkillGroup struct {
	Name   string
	Skills []resume.Skill
}

func Render(w io.Writer, templatePath string, r resume.Resume, lang string) error {
	if lang != "fr" && lang != "en" {
		return fmt.Errorf("unsupported language %q", lang)
	}
	names := make(map[string]string, len(r.Skills))
	for _, skill := range r.Skills {
		names[skill.ID] = skill.Name
	}
	funcs := template.FuncMap{
		"text": func(v resume.Localized) string { return v.Get(lang) },
		"esc":  Escape,
		"date": func(v string) string {
			if v == "present" {
				if lang == "fr" {
					return "Aujourd'hui"
				}
				return "Present"
			}
			return v
		},
		"period": func(start, end string) string {
			if start == end {
				return start
			}
			if end == "present" {
				if lang == "fr" {
					end = "Aujourd'hui"
				} else {
					end = "Present"
				}
			}
			return start + "–" + end
		},
		"displayURL": func(raw string) string {
			value := strings.TrimPrefix(raw, "https://")
			value = strings.TrimPrefix(value, "http://")
			value = strings.TrimPrefix(value, "www.")
			return strings.TrimSuffix(value, "/")
		},
		"tech": func(ids []string) string {
			out := make([]string, 0, len(ids))
			for _, id := range ids {
				name := names[id]
				if name == "" {
					name = id
				}
				out = append(out, name)
			}
			return strings.Join(out, ", ")
		},
	}
	raw, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("read template: %w", err)
	}
	t, err := template.New("ats.tex").Funcs(funcs).Parse(string(raw))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	labels := map[string]string{"summary": "SUMMARY", "skills": "TECHNICAL SKILLS", "experience": "EXPERIENCE", "education": "EDUCATION", "projects": "PERSONAL PROJECTS", "languages": "LANGUAGES"}
	if lang == "fr" {
		labels = map[string]string{"summary": "PROFIL", "skills": "COMPÉTENCES TECHNIQUES", "experience": "EXPÉRIENCE", "education": "FORMATION", "projects": "PROJETS PERSONNELS", "languages": "LANGUES"}
	}
	categoryNames := map[string]resume.Localized{
		"backend":      {"fr": "Backend", "en": "Backend"},
		"data":         {"fr": "Données", "en": "Data"},
		"distributed":  {"fr": "Systèmes distribués", "en": "Distributed Systems"},
		"platform":     {"fr": "Cloud & DevOps", "en": "Cloud & DevOps"},
		"production":   {"fr": "Observabilité & Analytics", "en": "Observability & Analytics"},
		"frontend":     {"fr": "Frontend", "en": "Frontend"},
		"tooling":      {"fr": "Outillage", "en": "Tooling"},
		"integrations": {"fr": "Intégrations", "en": "Integrations"},
		"ai-assisted":  {"fr": "Développement assisté par agents", "en": "AI-Assisted Engineering - Agentic Workflows"},
		"historical":   {"fr": "Autres technologies", "en": "Other Technologies"},
	}
	orderedCategories := []string{"backend", "data", "distributed", "platform", "production", "frontend", "ai-assisted", "historical", "tooling"}
	byCategory := make(map[string][]resume.Skill)
	for _, skill := range r.Skills {
		if skill.Category == "experience-only" || skill.Category == "historical" || skill.Category == "ai-assisted" || skill.Category == "project-only" {
			continue
		}
		byCategory[skill.Category] = append(byCategory[skill.Category], skill)
	}
	groups := make([]SkillGroup, 0, len(byCategory))
	for _, category := range orderedCategories {
		if skills := byCategory[category]; len(skills) > 0 {
			groups = append(groups, SkillGroup{Name: categoryNames[category].Get(lang), Skills: skills})
			delete(byCategory, category)
		}
	}
	for category, skills := range byCategory {
		groups = append(groups, SkillGroup{Name: category, Skills: skills})
	}
	if err := t.ExecuteTemplate(w, "ats.tex", View{Resume: r, Lang: lang, Labels: labels, Skills: groups}); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	return nil
}

var replacer = strings.NewReplacer(
	`\`, `\textbackslash{}`, "&", `\&`, "%", `\%`, "$", `\$`, "#", `\#`,
	"_", `\_`, "{", `\{`, "}", `\}`, "~", `\textasciitilde{}`, "^", `\textasciicircum{}`,
)

func Escape(s string) string { return replacer.Replace(s) }
