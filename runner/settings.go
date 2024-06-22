package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pilu/config"
)

const (
	envSettingsPrefix   = "RUNNER_"
	mainSettingsSection = "Settings"
)

var colors = map[string]string{
	"reset":          "0",
	"black":          "30",
	"red":            "31",
	"green":          "32",
	"yellow":         "33",
	"blue":           "34",
	"magenta":        "35",
	"cyan":           "36",
	"white":          "37",
	"bold_black":     "30;1",
	"bold_red":       "31;1",
	"bold_green":     "32;1",
	"bold_yellow":    "33;1",
	"bold_blue":      "34;1",
	"bold_magenta":   "35;1",
	"bold_cyan":      "36;1",
	"bold_white":     "37;1",
	"bright_black":   "30;2",
	"bright_red":     "31;2",
	"bright_green":   "32;2",
	"bright_yellow":  "33;2",
	"bright_blue":    "34;2",
	"bright_magenta": "35;2",
	"bright_cyan":    "36;2",
	"bright_white":   "37;2",
}

func logColor(s *mySetting, logName string) string {
	settingsKey := fmt.Sprintf("log_color_%s", logName)
	colorName := s.settings[settingsKey]

	return colors[colorName]
}

func (l *mySetting) loadEnvSettings() {
	for key, _ := range l.settings {
		envKey := fmt.Sprintf("%s%s", envSettingsPrefix, strings.ToUpper(key))
		if value := os.Getenv(envKey); value != "" {
			l.settings[key] = value
		}
	}
}

func (l *mySetting) loadRunnerConfigSettings() {
	if _, err := os.Stat(l.configPath()); err != nil {
		return
	}

	logger.Printf("Loading settings from %s", l.configPath())
	sections, err := config.ParseFile(l.configPath(), mainSettingsSection)
	if err != nil {
		return
	}

	for key, value := range sections[mainSettingsSection] {
		l.settings[key] = value
	}
}

func initSettings(s *mySetting) {
	s.loadEnvSettings()
	s.loadRunnerConfigSettings()
}

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

type TypeGetBuildCommand func(st map[string]string, arguments map[string]string) []string

type mySetting struct {
	settings        map[string]string
	getBuildCommand TypeGetBuildCommand
}

func NewMySetting(settings map[string]string, getBuildCommand TypeGetBuildCommand) *mySetting {
	return &mySetting{
		settings:        settings,
		getBuildCommand: getBuildCommand,
	}
}

func (l *mySetting) root() string {
	return l.settings["root"]
}

func (l *mySetting) tmpPath() string {
	return l.settings["tmp_path"]
}

func (l *mySetting) buildCommand() string {
	return l.settings["build_commands"]
}

func (l *mySetting) buildName() string {
	return l.settings["build_name"]
}
func (l *mySetting) buildPath() string {
	p := filepath.Join(l.tmpPath(), l.buildName())
	if runtime.GOOS == "windows" && filepath.Ext(p) != ".exe" {
		p += ".exe"
	}
	return p
}

func (l *mySetting) buildErrorsFileName() string {
	return l.settings["build_log"]
}

func (l *mySetting) buildErrorsFilePath() string {
	return filepath.Join(l.tmpPath(), l.buildErrorsFileName())
}

func (l *mySetting) configPath() string {
	return l.settings["config_path"]
}

func (l *mySetting) buildDelay() time.Duration {
	value, _ := strconv.Atoi(l.settings["build_delay"])

	return time.Duration(value)
}
