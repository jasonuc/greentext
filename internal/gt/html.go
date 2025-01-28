package gt

import (
	"fmt"
	"html/template"
	"os"
)

func GenerateHTMLFile(outputFile string, data GTData) error {
	tmpl := template.New("template")
	if isFilePath(data.Template) {
		if _, err := tmpl.ParseFiles(string(data.Template)); err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}
	} else {
		if _, err := tmpl.Parse(string(data.Template)); err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}
	}

	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func isFilePath(template []byte) bool {
	info, err := os.Stat(string(template))
	if err != nil {
		return false // returns false because template is not file path
	}
	return !info.IsDir() // returns true because template is file path
}
