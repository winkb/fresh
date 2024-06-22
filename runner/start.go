package runner

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

type Starter struct {
	startChannel chan string
	stopChannel  chan bool
}

func newStarter() *Starter {
	return &Starter{
		startChannel: make(chan string, 1000),
		stopChannel:  make(chan bool),
	}
}

var (
	mainLog    logFunc
	watcherLog logFunc
	runnerLog  logFunc
	buildLog   logFunc
	appLog     logFunc
)

func (l *Starter) flushEvents() {
	for {
		select {
		case eventName := <-l.startChannel:
			mainLog("receiving event %s", eventName)
		default:
			return
		}
	}
}

func (l *Starter) start(ctx context.Context, s *mySetting) {
	loopIndex := 0
	buildDelay := s.buildDelay()

	started := false

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				loopIndex++
				mainLog("Waiting (loop %d)...", loopIndex)
				changeInfo := <-l.startChannel

				mainLog("receiving first event %s", changeInfo)
				mainLog("sleeping for %d milliseconds", buildDelay)
				time.Sleep(buildDelay * time.Millisecond)
				mainLog("flushing events")

				l.flushEvents()

				mainLog("Started! (%d Goroutines)", runtime.NumGoroutine())
				err := removeBuildErrorsLog(s)
				if err != nil {
					mainLog(err.Error())
				}

				buildFailed := false

				var eventName string
				var filename string

				if strings.Contains(changeInfo, ":") {
					filename = strings.TrimSpace(strings.Split(changeInfo, ":")[0])
					filename = strings.TrimFunc(filename, func(r rune) bool {
						return r == '"'
					})
					eventName = strings.Split(changeInfo, ":")[1]
				}

				if shouldRebuild(s, changeInfo) {
					errorMessage, ok := build(s, map[string]string{
						"event_name": eventName,
						"filename":   filename,
					})
					if !ok {
						buildFailed = true
						mainLog("Build Failed: \n %s", errorMessage)
						if !started {
							os.Exit(1)
						}
						createBuildErrorsLog(s, errorMessage)
					}
				}

				if !buildFailed {
					if started {
						// stopChannel <- true
					}
					// don't run the app if the build failed
					// run()

					started = true
					mainLog(strings.Repeat("-", 20))

				}
			}
		}
	}()
}

func initLogFuncs(s *mySetting) {
	mainLog = newLogFunc(s, "main")
	watcherLog = newLogFunc(s, "watcher")
	runnerLog = newLogFunc(s, "runner")
	buildLog = newLogFunc(s, "build")
	appLog = newLogFunc(s, "app")
}

func setEnvVars(s *mySetting) {
	os.Setenv("DEV_RUNNER", "1")
	wd, err := os.Getwd()
	if err == nil {
		os.Setenv("RUNNER_WD", wd)
	}

	for k, v := range s.settings {
		key := strings.ToUpper(fmt.Sprintf("%s%s", envSettingsPrefix, k))
		os.Setenv(key, v)
	}
}

// Watches for file changes in the root directory.
// After each file system event it builds and (re)starts the application.
func Start(ctx context.Context, s *mySetting) {

	var defaultSetting = map[string]string{
		"config_path":       "./runner.conf",
		"build_commands":    "[\"ls\"]",
		"root":              ".",
		"tmp_path":          "./tmp",
		"build_name":        "runner-build",
		"build_log":         "runner-build-errors.log",
		"valid_ext":         ".go, .tpl, .tmpl, .html",
		"no_rebuild_ext":    ".tpl, .tmpl, .html",
		"ignored":           "assets, tmp",
		"build_delay":       "600",
		"colors":            "1",
		"log_color_main":    "cyan",
		"log_color_build":   "yellow",
		"log_color_runner":  "green",
		"log_color_watcher": "magenta",
		"log_color_app":     "",
	}

	for k, v := range defaultSetting {
		if _, ok := s.settings[k]; !ok {
			s.settings[k] = v
		}
	}

	t := newStarter()

	initLimit()
	initSettings(s)

	initLogFuncs(s)
	initFolders(s)
	setEnvVars(s)
	t.watch(ctx, s)
	t.start(ctx, s)

	<-ctx.Done()
}
