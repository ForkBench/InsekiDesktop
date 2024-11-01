package core

import (
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
