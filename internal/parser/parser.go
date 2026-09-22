package parser

import (
	"fmt"
	"os"

	"github.com/moulindavid/cv-builder/internal/resume"
	"gopkg.in/yaml.v3"
)

func ResumeFile(path string) (resume.Resume, error) {
	var v resume.Resume
	err := decode(path, &v)
	return v, err
}
func TargetFile(path string) (resume.Target, error) {
	var v resume.Target
	err := decode(path, &v)
	return v, err
}
func decode(path string, dst any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	d := yaml.NewDecoder(f)
	d.KnownFields(true)
	if err := d.Decode(dst); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}
