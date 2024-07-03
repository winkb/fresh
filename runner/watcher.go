package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/howeyc/fsnotify"
)

func (l *Starter) watchFolder(path string, ctx context.Context, s *mySetting) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	eventMp := map[string]int64{}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(s, ev.Name) {
					var eventStr = ev.String()
					var now = time.Now().Unix()

					if now-eventMp[eventStr] > 2 {
						watcherLog("sending event %s", ev)
						l.startChannel <- eventStr
						eventMp[eventStr] = now
					}

				}
			case <-ctx.Done():
				return
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func (l *Starter) watch(ctx context.Context, s *mySetting) {
	root := s.root()
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(s, path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			if isIgnoredFolder(s, path) {
				watcherLog("Ignoring %s", path)
				return filepath.SkipDir
			}

			l.watchFolder(path, ctx, s)
		}

		return err
	})
}
