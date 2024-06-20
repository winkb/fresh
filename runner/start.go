package runner

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	startChannel chan string
	stopChannel  chan bool
	mainLog      logFunc
	watcherLog   logFunc
	runnerLog    logFunc
	buildLog     logFunc
	appLog       logFunc
)

func flushEvents() {
	for {
		select {
		case eventName := <-startChannel:
			mainLog("receiving event %s", eventName)
		default:
			return
		}
	}
}

func start(ctx context.Context) {
	loopIndex := 0
	buildDelay := buildDelay()

	started := false

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				loopIndex++
				mainLog("Waiting (loop %d)...", loopIndex)
				eventName := <-startChannel

				mainLog("receiving first event %s", eventName)
				mainLog("sleeping for %d milliseconds", buildDelay)
				time.Sleep(buildDelay * time.Millisecond)
				mainLog("flushing events")

				flushEvents()

				mainLog("Started! (%d Goroutines)", runtime.NumGoroutine())
				err := removeBuildErrorsLog()
				if err != nil {
					mainLog(err.Error())
				}

				buildFailed := false
				if shouldRebuild(eventName) {
					errorMessage, ok := build()
					if !ok {
						buildFailed = true
						mainLog("Build Failed: \n %s", errorMessage)
						if !started {
							os.Exit(1)
						}
						createBuildErrorsLog(errorMessage)
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

func init() {
	startChannel = make(chan string, 1000)
	stopChannel = make(chan bool)
}

func initLogFuncs() {
	mainLog = newLogFunc("main")
	watcherLog = newLogFunc("watcher")
	runnerLog = newLogFunc("runner")
	buildLog = newLogFunc("build")
	appLog = newLogFunc("app")
}

func setEnvVars() {
	os.Setenv("DEV_RUNNER", "1")
	wd, err := os.Getwd()
	if err == nil {
		os.Setenv("RUNNER_WD", wd)
	}

	for k, v := range Settings {
		key := strings.ToUpper(fmt.Sprintf("%s%s", envSettingsPrefix, k))
		os.Setenv(key, v)
	}
}

// Watches for file changes in the root directory.
// After each file system event it builds and (re)starts the application.
func Start(ctx context.Context, filePath string) {
	initLimit()
	initSettings()

	Settings["root"] = filePath

	initLogFuncs()
	initFolders()
	setEnvVars()
	watch(ctx)
	start(ctx)
	startChannel <- "/"

	<-ctx.Done()
}
