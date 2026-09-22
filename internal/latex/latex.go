package latex

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Compile(texPath, pdfPath string) error {
	engine, err := exec.LookPath("lualatex")
	if err != nil {
		return fmt.Errorf("lualatex not found: install a LuaLaTeX distribution: %w", err)
	}
	outDir := filepath.Dir(texPath)
	cmd := exec.Command(engine, "-interaction=nonstopmode", "-halt-on-error", "-output-directory", outDir, texPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lualatex: %w\n%s", err, out)
	}
	generated := filepath.Join(outDir, filepath.Base(texPath[:len(texPath)-len(filepath.Ext(texPath))])+".pdf")
	if err := os.Rename(generated, pdfPath); err != nil {
		return fmt.Errorf("move PDF: %w", err)
	}
	return nil
}
