package file

import (
	"github.com/amyangfei/dmux/registry/models"
	"github.com/amyangfei/dmux/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func walk(dir string) (rFiles, lnNames, lnTargets []string, err error) {
	files, _err := ioutil.ReadDir(dir)
	if _err != nil {
		err = _err
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			if file.Mode().IsRegular() {
				rFiles = append(rFiles, file.Name())
			} else {
				lnfile, _err := os.Readlink(filepath.Join([]string{dir, file.Name()}...))
				if _err != nil {
					err = _err
					return
				}
				lnstat, _err := os.Stat(lnfile)
				if _err != nil {
					err = _err
					return
				}
				if !lnstat.IsDir() {
					lnNames = append(lnNames, file.Name())
					lnTargets = append(lnTargets, lnfile)
				}
			}
		}
	}
	return
}

// LookupFiles find all candinate files according to the given entry
func LookupFiles(entry *models.Entry) (result []string, err error) {
	if len(entry.Include) > 0 {
		for _, filename := range entry.Include {
			result = append(
				result, filepath.Join([]string{entry.Path, filename}...))
		}
	} else {
		rFiles, lnNames, lnTargets, _err := walk(entry.Path)
		if _err != nil {
			err = _err
			return
		}
		for _, filename := range rFiles {
			if !utils.In(filename, entry.Exclude) {
				result = append(
					result, filepath.Join([]string{entry.Path, filename}...))
			}
		}
		for i := 0; i < len(lnNames); i++ {
			if !utils.In(lnNames[i], entry.Exclude) {
				result = append(result, lnTargets[i])
			}
		}
	}
	return
}
