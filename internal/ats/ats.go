package ats

import (
	"fmt"
	"os/exec"
	"strings"
)

func Extract(pdf string) (string, error) {
	cmd := exec.Command("pdftotext", pdf, "-")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pdftotext: %w: %s", err, out)
	}
	return string(out), nil
}
func CheckOrder(text string, sections ...string) error {
	last := -1
	upper := strings.ToUpper(text)
	for _, s := range sections {
		i := strings.Index(upper, strings.ToUpper(s))
		if i < 0 {
			return fmt.Errorf("section %q not found", s)
		}
		if i <= last {
			return fmt.Errorf("section %q is out of order", s)
		}
		last = i
	}
	return nil
}
