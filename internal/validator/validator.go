package validator

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"

	"github.com/moulindavid/cv-builder/internal/resume"
)

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
		if id == "" {
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
		if s.Name == "" {
			errs = append(errs, "skill "+s.ID+" has no name")
		}
	}
	for _, e := range r.Experience {
		add(e.ID, "experience")
		if e.Company == "" || e.Start == "" {
			errs = append(errs, "experience "+e.ID+" requires company and start")
		}
		for _, id := range e.Technologies {
			if !skillIDs[id] {
				errs = append(errs, fmt.Sprintf("experience %s references unknown skill %s", e.ID, id))
			}
		}
		for _, b := range e.Bullets {
			add(e.ID+"/"+b.ID, "bullet")
			requireLocales(&errs, "bullet "+e.ID+"/"+b.ID, b.Text)
		}
	}
	for _, e := range r.Education {
		add(e.ID, "education")
	}
	for _, l := range r.Links {
		if u, err := url.ParseRequestURI(l.URL); err != nil || u.Scheme == "" {
			errs = append(errs, "invalid link: "+l.URL)
		}
	}
	requireLocales(&errs, "profile.titles", r.Profile.Titles)
	requireLocales(&errs, "summaries", r.Summaries)
	if len(errs) > 0 {
		return fmt.Errorf("resume validation failed:\n- %s", strings.Join(errs, "\n- "))
	}
	return nil
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
