package gt

import (
	"fmt"
	"html/template"
	"os"
)

func GenerateHTMLFile(outputFile string, data GTData) error {

	tmpl, err := template.New("greentext").Parse(string(data.Template))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
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
