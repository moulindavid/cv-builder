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
	positions := make(map[string]int, len(sections))
	for i, line := range strings.Split(text, "\n") {
		line = strings.ToUpper(strings.TrimSpace(line))
		if _, exists := positions[line]; !exists {
			positions[line] = i
		}
	}
	last := -1
	for _, section := range sections {
		i, found := positions[strings.ToUpper(strings.TrimSpace(section))]
		if !found {
			return fmt.Errorf("section %q not found", section)
		}
		if i <= last {
			return fmt.Errorf("section %q is out of order", section)
		}
		last = i
	}
	return nil
}
