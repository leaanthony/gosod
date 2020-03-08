package templatedir

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
)

// TemplateDir defines a directory containing directories and files, including template files
type TemplateDir struct {
	path           string
	templateFilter string
	dirs           []string
	standardFiles  []string
	templateFiles  []string
	ignoredFiles   map[string]struct{}
}

// New attempts to create a new TemplateDir from the given (absolute) path
func New(path string) (*TemplateDir, error) {

	// If the path does not exist then return an error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("Path " + path + " does not exist")
	}

	return &TemplateDir{
		path:           path,
		templateFilter: ".tmpl",
		ignoredFiles:   make(map[string]struct{}),
	}, nil

}

// IgnoreFilename will add the given filename to the list of files to ignore
// during extraction
func (t *TemplateDir) IgnoreFilename(filename string) {
	t.ignoredFiles[filename] = struct{}{}
}

// Extract the templates to the given directory, using data as input
func (t *TemplateDir) Extract(targetDirectory string, data interface{}) error {

	// Get the absolute path
	targetDirectory, err := filepath.Abs(targetDirectory)
	if err != nil {
		return err
	}

	// If the targetDirectory doesn't exist, then create it
	if _, err := os.Stat(targetDirectory); os.IsNotExist(err) == true {
		// Create the targetDirectory
		err = os.MkdirAll(targetDirectory, 0755)
		if err != nil {
			return err
		}
	}

	// Process the template files
	err = t.processTemplateDirFiles(targetDirectory, data)
	if err != nil {
		return err
	}

	return nil
}

func (t *TemplateDir) processTemplateDirFiles(targetDirectory string, data interface{}) error {
	// Categorise all files
	err := t.categoriseFiles()
	if err != nil {
		return err
	}

	// Create all directories
	err = t.createDirectories(targetDirectory)
	if err != nil {
		return err
	}

	// Process TemplateDirs
	err = t.processTemplateDirs(targetDirectory, data)
	if err != nil {
		return err
	}

	// Copy files
	err = t.copyFiles(targetDirectory)
	if err != nil {
		return err
	}

	return nil
}

func (t *TemplateDir) categoriseFiles() error {
	return filepath.Walk(t.path, t.categoriseFile)
}

func (t *TemplateDir) categoriseFile(path string, info os.FileInfo, err error) error {

	// Process error
	if err != nil {
		return err
	}

	// Is it a directory?
	if info.IsDir() {
		// Ignore base dir
		if path != t.path {
			t.dirs = append(t.dirs, path)
		}
		return nil
	}

	// Get the filename
	filename := filepath.Base(path)

	// Is it a file we are ignoring?
	_, ignored := t.ignoredFiles[filename]
	if ignored {
		return nil
	}

	// Is it a template?
	if strings.Index(filename, t.templateFilter) > -1 {
		t.templateFiles = append(t.templateFiles, path)
		return nil
	}

	// Treat as standard file
	t.standardFiles = append(t.standardFiles, path)
	return nil
}

func (t *TemplateDir) convertPathTarget(path string, targetDirectory string) string {
	relativePath := strings.TrimPrefix(path, t.path)
	result := filepath.Join(targetDirectory, relativePath)
	return result
}

func (t *TemplateDir) createDirectories(targetDirectory string) error {

	// Iterate all directories and attempt to create them
	for _, dirPath := range t.dirs {

		targetDir := t.convertPathTarget(dirPath, targetDirectory)

		// Create the directory
		err := os.MkdirAll(targetDir, 0755)

		// Ignore directory exists errors
		if err != nil && err != syscall.EEXIST {
			return err
		}
	}

	return nil
}

func (t *TemplateDir) processTemplateDirs(targetDirectory string, data interface{}) error {

	// Iterate template files
	for _, templateFile := range t.templateFiles {

		// Parse template
		tmpl, err := template.ParseFiles(templateFile)
		if err != nil {
			return err
		}

		// Convert path to target path
		targetFile := t.convertPathTarget(templateFile, targetDirectory)

		// update filename
		baseDir := filepath.Dir(targetFile)
		filename := filepath.Base(targetFile)
		filename = strings.ReplaceAll(filename, t.templateFilter, "")
		targetFile = filepath.Join(baseDir, filename)

		// Create target file
		writer, err := os.Create(targetFile)
		if err != nil {
			return err
		}

		err = tmpl.Execute(writer, data)
		if err != nil {
			writer.Close()
			return err
		}

		writer.Close()

	}

	return nil
}

func (t *TemplateDir) copyFiles(targetDirectory string) error {

	// Iterate over files
	for _, filename := range t.standardFiles {

		targetFilename := t.convertPathTarget(filename, targetDirectory)
		err := t.copyFile(filename, targetFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TemplateDir) copyFile(source, target string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(target)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
