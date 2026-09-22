package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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
	if err := validateLanguage(*lang); err != nil {
		return err
	}
	r, err := load()
	if err != nil {
		return err
	}
	slug := "resume"
	if *target != "" {
		t, err := loadTarget(*target, r)
		if err != nil {
			return err
		}
		r = tailor.ApplyTarget(r, t, *lang)
		slug = t.Name
	}
	pdf, err := generate(r, *lang, slug)
	if err != nil {
		return err
	}
	fmt.Println(pdf)
	return nil
}

func extract(args []string) error {
	fs := flag.NewFlagSet("extract-text", flag.ContinueOnError)
	lang := fs.String("lang", "en", "language: fr or en")
	target := fs.String("target", "", "target name or YAML path")
	pdf := fs.String("pdf", "", "PDF path (overrides --lang and --target)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *pdf == "" {
		if err := validateLanguage(*lang); err != nil {
			return err
		}
		r, err := load()
		if err != nil {
			return err
		}
		slug := "resume"
		if *target != "" {
			t, err := loadTarget(*target, r)
			if err != nil {
				return err
			}
			slug = t.Name
		}
		*pdf = outputPath(r.Profile.Name, slug, *lang, ".pdf")
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
	lang := fs.String("lang", "en", "language: fr or en")
	buildPDF := fs.Bool("build", false, "build a PDF using the job ranking")
	maxBullets := fs.Int("max-bullets", 0, "maximum bullets per experience in the built CV (0 keeps all)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*job) == "" {
		return fmt.Errorf("--job is required")
	}
	if err := validateLanguage(*lang); err != nil {
		return err
	}
	if *maxBullets < 0 {
		return fmt.Errorf("--max-bullets cannot be negative")
	}
	raw, err := os.ReadFile(*job)
	if err != nil {
		return fmt.Errorf("read job: %w", err)
	}
	r, err := load()
	if err != nil {
		return err
	}
	report := tailor.Analyze(r, string(raw), *lang)
	fmt.Print(tailor.Format(report))
	if !*buildPDF {
		return nil
	}
	r = tailor.ApplyTarget(r, resume.Target{Priorities: report.Priorities, MaxBullets: *maxBullets}, *lang)
	pdf, err := generate(r, *lang, "tailored")
	if err != nil {
		return err
	}
	fmt.Println("Built PDF:", pdf)
	return nil
}

func loadTarget(value string, r resume.Resume) (resume.Target, error) {
	path := value
	if filepath.Ext(path) == "" {
		path = filepath.Join("targets", path+".yaml")
	}
	t, err := parser.TargetFile(path)
	if err != nil {
		return t, err
	}
	if err := validator.ValidateTarget(t, r); err != nil {
		return t, fmt.Errorf("validate %s: %w", path, err)
	}
	return t, nil
}

func generate(r resume.Resume, lang, slug string) (string, error) {
	if err := os.MkdirAll("output", 0755); err != nil {
		return "", err
	}
	tex := outputPath(r.Profile.Name, slug, lang, ".tex")
	f, err := os.Create(tex)
	if err != nil {
		return "", err
	}
	renderErr := renderer.Render(f, "templates/ats.tex", r, lang)
	closeErr := f.Close()
	if renderErr != nil {
		return "", renderErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	pdf := outputPath(r.Profile.Name, slug, lang, ".pdf")
	if err := latex.Compile(tex, pdf); err != nil {
		return "", err
	}
	return pdf, nil
}

func outputPath(name, target, lang, ext string) string {
	base := fmt.Sprintf("%s-%s-%s%s", slugify(name), slugify(target), lang, ext)
	return filepath.Join("output", base)
}

func validateLanguage(lang string) error {
	if lang != "fr" && lang != "en" {
		return fmt.Errorf("unsupported language %q", lang)
	}
	return nil
}

func slugify(value string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			dash = false
		} else if b.Len() > 0 && !dash {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
