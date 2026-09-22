package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/moulindavid/cv-builder/internal/ats"
	"github.com/moulindavid/cv-builder/internal/latex"
	"github.com/moulindavid/cv-builder/internal/parser"
	"github.com/moulindavid/cv-builder/internal/renderer"
	"github.com/moulindavid/cv-builder/internal/resume"
	"github.com/moulindavid/cv-builder/internal/tailor"
	"github.com/moulindavid/cv-builder/internal/validator"
)

const dataPath = "data/resume.yaml"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}
	switch args[0] {
	case "build":
		return build(args[1:])
	case "validate":
		return validate()
	case "extract-text":
		return extract(args[1:])
	case "tailor":
		return tailorJob(args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
func usage() error {
	fmt.Fprintln(os.Stderr, "Usage: cv <build|validate|extract-text|tailor> [options]")
	return flag.ErrHelp
}
func load() (resume.Resume, error) {
	r, err := parser.ResumeFile(dataPath)
	if err != nil {
		return r, err
	}
	if err := validator.Validate(r); err != nil {
		return r, err
	}
	return r, nil
}
func validate() error {
	_, err := load()
	if err == nil {
		fmt.Println("resume is valid")
	}
	return err
}
func build(args []string) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	lang := fs.String("lang", "fr", "language: fr or en")
	target := fs.String("target", "", "target name or YAML path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	r, err := load()
	if err != nil {
		return err
	}
	slug := "senior-backend"
	if *target != "" {
		path := *target
		if filepath.Ext(path) == "" {
			path = filepath.Join("targets", path+".yaml")
		}
		t, err := parser.TargetFile(path)
		if err != nil {
			return err
		}
		r = tailor.ApplyTarget(r, t, *lang)
		if t.Name != "" {
			slug = t.Name
		}
	}
	if err := os.MkdirAll("output", 0755); err != nil {
		return err
	}
	base := fmt.Sprintf("david-moulin-%s-%s", slug, *lang)
	tex := filepath.Join("output", base+".tex")
	f, err := os.Create(tex)
	if err != nil {
		return err
	}
	renderErr := renderer.Render(f, "templates/ats.tex", r, *lang)
	closeErr := f.Close()
	if renderErr != nil {
		return renderErr
	}
	if closeErr != nil {
		return closeErr
	}
	pdf := filepath.Join("output", base+".pdf")
	if err := latex.Compile(tex, pdf); err != nil {
		return err
	}
	fmt.Println(pdf)
	return nil
}
func extract(args []string) error {
	fs := flag.NewFlagSet("extract-text", flag.ContinueOnError)
	lang := fs.String("lang", "en", "language")
	pdf := fs.String("pdf", "", "PDF path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *pdf == "" {
		*pdf = filepath.Join("output", "david-moulin-senior-backend-"+*lang+".pdf")
	}
	text, err := ats.Extract(*pdf)
	if err != nil {
		return err
	}
	fmt.Print(text)
	return nil
}
func tailorJob(args []string) error {
	fs := flag.NewFlagSet("tailor", flag.ContinueOnError)
	job := fs.String("job", "", "job description file")
	lang := fs.String("lang", "en", "language")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*job) == "" {
		return fmt.Errorf("--job is required")
	}
	raw, err := os.ReadFile(*job)
	if err != nil {
		return fmt.Errorf("read job: %w", err)
	}
	r, err := load()
	if err != nil {
		return err
	}
	fmt.Print(tailor.Format(tailor.Analyze(r, string(raw), *lang)))
	return nil
}
