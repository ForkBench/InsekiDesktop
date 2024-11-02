package core

import (
	"encoding/base64"
	"github.com/ForkBench/Inseki-Core/tools"
	"os"
	"path/filepath"
)

type Analyze struct {
	Config       tools.Config
	InsekiIgnore []string
	Home         string
}

type File struct {
	Path     string
	FileName string
	IconPath string
}

func (f File) GetB64Path() string {
	return B64Encode(f.Path)
}

func FileFromUrl(url string) (File, error) {
	path, err := B64Decode(url)
	if err != nil {
		return File{}, err
	}

	// Check if the file exists
	fileinfo, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}

	iconPath := "/folder.png"

	// Check if it's a directory
	if !fileinfo.IsDir() {
		iconPath = "/file.png"
	}

	return File{Path: path, FileName: filepath.Base(path), IconPath: iconPath}, nil
}

func (a Analyze) Process(path string) (error, []tools.Response) {
	return tools.Process(path, a.Config, a.InsekiIgnore)
}

func (a Analyze) GetMainFolders() []File {
	// Find ~/Desktop, ~/Documents, ~/Downloads
	// TODO: Adapt to Windows and different languages
	paths := []File{
		{Path: a.Home + "/Desktop", FileName: "Desktop", IconPath: "/desktop.png"},
		{Path: a.Home + "/Documents", FileName: "Documents", IconPath: "/documents.png"},
		{Path: a.Home + "/Downloads", FileName: "Downloads", IconPath: "/downloads.png"},
	}

	return a.foldersStringToAbsolute(paths)
}

func (a Analyze) foldersStringToAbsolute(paths []File) []File {
	// Replace ~ with the home
	filesAbs := make([]File, 0)

	for _, file := range paths {
		absPath, _ := filepath.Abs(file.Path)
		filesAbs = append(filesAbs, File{Path: absPath, FileName: file.FileName, IconPath: file.IconPath})
	}

	return filesAbs
}

func (a Analyze) ListAllSubFiles(path string) []File {
	// List all subfiles of a folder
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	subFiles := make([]File, 0)

	for _, file := range files {
		icon := "/folder.png"
		if !file.IsDir() {
			icon = "/file.png"
		}

		subFiles = append(subFiles, File{Path: path + "/" + file.Name(), FileName: file.Name(), IconPath: icon})
	}

	return subFiles
}

func B64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func B64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (f File) ReadFile() string {
	content, err := ReadFile(f.Path)
	if err != nil {
		return ""
	}

	return content
}

func (f File) IsDirectory() bool {
	return f.IconPath == "/folder.png"
}

func (f File) FileType() string {
	if f.IsDirectory() {
		return "directory"
	}

	return "file"
}
