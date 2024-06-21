package runner

import (
	"os"
	"path/filepath"
	"strings"
)

func initFolders(s *mySetting) {
	runnerLog("InitFolders")
	path := s.tmpPath()
	runnerLog("mkdir %s", path)
	err := os.Mkdir(path, 0755)
	if err != nil {
		runnerLog(err.Error())
	}
}

func isTmpDir(s *mySetting, path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(s.tmpPath())

	return absolutePath == absoluteTmpPath
}

func isIgnoredFolder(s *mySetting, path string) bool {
	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	for _, e := range strings.Split(s.settings["ignored"], ",") {
		if strings.TrimSpace(e) == paths[0] {
			return true
		}
	}
	return false
}

func isWatchedFile(s *mySetting, path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(s.tmpPath())

	if strings.HasPrefix(absolutePath, absoluteTmpPath) {
		return false
	}

	ext := filepath.Ext(path)

	for _, e := range strings.Split(s.settings["valid_ext"], ",") {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}

func shouldRebuild(s *mySetting, eventName string) bool {
	for _, e := range strings.Split(s.settings["no_rebuild_ext"], ",") {
		e = strings.TrimSpace(e)
		fileName := strings.Replace(strings.Split(eventName, ":")[0], `"`, "", -1)
		if strings.HasSuffix(fileName, e) {
			return false
		}
	}

	return true
}

func createBuildErrorsLog(s *mySetting, message string) bool {
	file, err := os.Create(s.buildErrorsFilePath())
	if err != nil {
		return false
	}

	_, err = file.WriteString(message)
	if err != nil {
		return false
	}

	return true
}

func removeBuildErrorsLog(s *mySetting) error {
	err := os.Remove(s.buildErrorsFilePath())

	return err
}
