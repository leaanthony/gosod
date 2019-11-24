package gosod

import "github.com/leaanthony/gosod/internal/templatedir"

// TemplateDir creates a new TemplateDir structure for the given directoryPath
func TemplateDir(directoryPath string) (*templatedir.TemplateDir, error) {
	return templatedir.New(directoryPath)
}
