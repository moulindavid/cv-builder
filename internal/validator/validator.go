package validator

import (
	"fmt"
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	"github.com/moulindavid/cv-builder/internal/resume"
)

var supportedCategories = map[string]bool{
	"backend": true, "data": true, "distributed": true, "platform": true,
	"production": true, "frontend": true, "tooling": true, "integrations": true,
	"ai-assisted": true, "historical": true, "experience-only": true, "project-only": true,
}

func Validate(r resume.Resume) error {
	var errs []string
	if strings.TrimSpace(r.Profile.Name) == "" {
		errs = append(errs, "profile.name is required")
	}
	if _, err := mail.ParseAddress(r.Contact.Email); err != nil {
		errs = append(errs, "contact.email is invalid")
	}
	ids := map[string]string{}
	skillIDs := map[string]bool{}
	add := func(id, kind string) {
		if strings.TrimSpace(id) == "" {
			errs = append(errs, kind+" id is required")
			return
		}
		if old, ok := ids[id]; ok {
			errs = append(errs, fmt.Sprintf("duplicate id %q (%s and %s)", id, old, kind))
		}
		ids[id] = kind
	}
	for _, s := range r.Skills {
		add(s.ID, "skill")
		skillIDs[s.ID] = true
		if strings.TrimSpace(s.Name) == "" {
			errs = append(errs, "skill "+s.ID+" has no name")
		}
		if !supportedCategories[s.Category] {
			errs = append(errs, fmt.Sprintf("skill %s has unknown category %q", s.ID, s.Category))
		}
		aliases := map[string]bool{}
		for _, alias := range s.Aliases {
			key := strings.ToLower(strings.TrimSpace(alias))
			if key == "" {
				errs = append(errs, "skill "+s.ID+" has an empty alias")
			} else if aliases[key] {
				errs = append(errs, fmt.Sprintf("skill %s has duplicate alias %q", s.ID, alias))
			}
			aliases[key] = true
		}
	}
	for _, e := range r.Experience {
		add(e.ID, "experience")
		if strings.TrimSpace(e.Company) == "" {
			errs = append(errs, "experience "+e.ID+" has no company")
		}
		if strings.TrimSpace(e.Start) == "" {
			errs = append(errs, "experience "+e.ID+" has no start date")
		}
		requireLocales(&errs, "experience "+e.ID+" descriptions", e.Descriptions)
		if len(e.Roles) > 0 {
			requireLocales(&errs, "experience "+e.ID+" roles", e.Roles)
		}
		validateDates(&errs, "experience "+e.ID, e.Start, e.End)
		for _, id := range e.Technologies {
			if !skillIDs[id] {
				errs = append(errs, fmt.Sprintf("experience %s references unknown skill %s", e.ID, id))
			}
		}
		bulletIDs := map[string]bool{}
		for _, b := range e.Bullets {
			if bulletIDs[b.ID] {
				errs = append(errs, fmt.Sprintf("duplicate bullet id %q in experience %s", b.ID, e.ID))
			}
			bulletIDs[b.ID] = true
			if strings.TrimSpace(b.ID) == "" {
				errs = append(errs, "bullet id is required in experience "+e.ID)
			}
			requireLocales(&errs, "bullet "+e.ID+"/"+b.ID, b.Text)
		}
	}
	for _, project := range r.Projects {
		add(project.ID, "project")
		if strings.TrimSpace(project.Name) == "" {
			errs = append(errs, "project "+project.ID+" has no name")
		}
		for _, id := range project.Technologies {
			if !skillIDs[id] {
				errs = append(errs, fmt.Sprintf("project %s references unknown skill %s", project.ID, id))
			}
		}
	}
	for _, e := range r.Education {
		add(e.ID, "education")
		if strings.TrimSpace(e.Institution) == "" {
			errs = append(errs, "education "+e.ID+" has no institution")
		}
		requireLocales(&errs, "education "+e.ID+" degrees", e.Degrees)
		validateDates(&errs, "education "+e.ID, e.Start, e.End)
	}
	for i, l := range r.Languages {
		requireLocales(&errs, fmt.Sprintf("language %d name", i), l.Name)
		requireLocales(&errs, fmt.Sprintf("language %d level", i), l.Level)
	}
	externalKeywords := map[string]bool{}
	for _, keyword := range r.Tailoring.ExternalKeywords {
		key := strings.ToLower(strings.TrimSpace(keyword))
		if key == "" {
			errs = append(errs, "tailoring external keyword cannot be empty")
			continue
		}
		if externalKeywords[key] {
			errs = append(errs, fmt.Sprintf("duplicate tailoring external keyword %q", keyword))
		}
		externalKeywords[key] = true
	}
	for _, l := range r.Links {
		if strings.TrimSpace(l.Label) == "" {
			errs = append(errs, "link label is required")
		}
		u, err := url.ParseRequestURI(l.URL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			errs = append(errs, "invalid link: "+l.URL)
		}
	}
	requireLocales(&errs, "profile.titles", r.Profile.Titles)
	requireLocales(&errs, "summaries", r.Summaries)
	return validationError("resume", errs)
}

func ValidateTarget(t resume.Target, r resume.Resume) error {
	var errs []string
	if strings.TrimSpace(t.Name) == "" {
		errs = append(errs, "name is required")
	}
	if t.MaxBullets < 0 {
		errs = append(errs, "max_bullets cannot be negative")
	}
	if len(t.Titles) > 0 {
		requireLocales(&errs, "titles", t.Titles)
	}
	known := map[string]bool{}
	for _, skill := range r.Skills {
		known[strings.ToLower(skill.ID)] = true
		known[strings.ToLower(skill.Name)] = true
		for _, tag := range skill.Tags {
			known[strings.ToLower(tag)] = true
		}
	}
	for _, e := range r.Experience {
		for _, b := range e.Bullets {
			for _, tag := range b.Tags {
				known[strings.ToLower(tag)] = true
			}
		}
	}
	for _, priority := range t.Priorities {
		if !known[strings.ToLower(priority)] {
			errs = append(errs, fmt.Sprintf("unknown priority %q", priority))
		}
	}
	return validationError("target", errs)
}

func validateDates(errs *[]string, name, start, end string) {
	startYear, startOK := year(start)
	endYear, endOK := year(end)
	if start != "" && !startOK {
		*errs = append(*errs, name+" has invalid start date "+start)
	}
	if end != "" && end != "present" && !endOK {
		*errs = append(*errs, name+" has invalid end date "+end)
	}
	if startOK && endOK && startYear > endYear {
		*errs = append(*errs, name+" starts after it ends")
	}
}

func year(value string) (int, bool) {
	if len(value) != 4 {
		return 0, false
	}
	y, err := strconv.Atoi(value)
	return y, err == nil
}

func validationError(kind string, errs []string) error {
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s validation failed:\n- %s", kind, strings.Join(errs, "\n- "))
}

func requireLocales(errs *[]string, name string, l resume.Localized) {
	for _, lang := range []string{"fr", "en"} {
		if strings.TrimSpace(l[lang]) == "" {
			*errs = append(*errs, name+" missing "+lang)
		}
	}
	for lang := range l {
		if lang != "fr" && lang != "en" {
			*errs = append(*errs, name+" has unsupported locale "+lang)
		}
	}
}
