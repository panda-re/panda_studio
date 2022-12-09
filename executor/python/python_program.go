package python

import (
	"embed"
	"fmt"
	"log"
	"os"
)

//go:embed *.py *.pyi
var pythonScript embed.FS

func ExtractPythonScript(destDir string) error {
	files, err := pythonScript.ReadDir(".")
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Printf("Writing file %s\n", file.Name())
		fileContent, err := pythonScript.ReadFile(file.Name())
		if err != nil {
			return err
		}

		destFile := fmt.Sprintf("%s/%s", destDir, file.Name())
		if err := os.WriteFile(destFile, fileContent, 0666); err != nil {
			log.Printf("error os.WriteFile error: %v", err)
			return nil
		}
	}

	return nil
}