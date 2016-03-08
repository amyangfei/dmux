package file

import (
	"github.com/amyangfei/dmux/registry/models"
	"io/ioutil"
	"os"
	"path/filepath"
)

func walk(dir string) ([]string, error) {
	result := make([]string, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		if !file.IsDir() {
			if file.Mode().IsRegular() {
				result = append(result, file.Name())
			} else {
				lnfile, err := os.Readlink(filepath.Join([]string{dir, file.Name()}...))
				if err != nil {
					return result, err
				}
				lnstat, err := os.Stat(lnfile)
				if err != nil {
					return result, err
				}
				if !lnstat.IsDir() {
					result = append(result, lnfile)
				}
			}
		}
	}
	return result, nil
}

// LookupFiles find all candinate files according to the given entry
func LookupFiles(entry *models.Entry) ([]string, error) {
	files, err := walk(entry.Path)
	return files, err
}
