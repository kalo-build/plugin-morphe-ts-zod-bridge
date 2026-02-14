package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kalo-build/go-util/strcase"
)

// BridgeWriter handles writing generated converter files to disk.
type BridgeWriter struct {
	OutputPath string
}

// NewBridgeWriter creates a new BridgeWriter instance.
func NewBridgeWriter(outputPath string) *BridgeWriter {
	return &BridgeWriter{
		OutputPath: outputPath,
	}
}

// ensureDir creates a directory if it doesn't exist.
func (w *BridgeWriter) ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// writeFile writes content to a file, creating parent directories as needed.
func (w *BridgeWriter) writeFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := w.ensureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return os.WriteFile(path, content, 0644)
}

// WriteModelConverter writes a single model converter file.
func (w *BridgeWriter) WriteModelConverter(modelName string, content []byte) error {
	fileName := toFileName(modelName) + ".ts"
	filePath := filepath.Join(w.OutputPath, "models", fileName)
	return w.writeFile(filePath, content)
}

// WriteStructureConverter writes a single structure converter file.
func (w *BridgeWriter) WriteStructureConverter(structureName string, content []byte) error {
	fileName := toFileName(structureName) + ".ts"
	filePath := filepath.Join(w.OutputPath, "structures", fileName)
	return w.writeFile(filePath, content)
}

// WriteEntityConverter writes a single entity converter file.
func (w *BridgeWriter) WriteEntityConverter(entityName string, content []byte) error {
	fileName := toFileName(entityName) + ".ts"
	filePath := filepath.Join(w.OutputPath, "entities", fileName)
	return w.writeFile(filePath, content)
}

// toFileName converts a PascalCase type name to a kebab-case file name.
func toFileName(typeName string) string {
	return strings.ToLower(strcase.ToKebabCase(typeName))
}
