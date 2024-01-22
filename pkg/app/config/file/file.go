package file

import (
	"encoding/json"
	"fmt"
	"github.com/ch-rollup/ch-rollup/pkg/app/config"
	"github.com/fsnotify/fsnotify"
	"os"
	"sync"
)

type File struct {
	lastState config.Config
	watchers  []config.WatchFunc
	lock      sync.RWMutex

	fsWatcher *fsnotify.Watcher
}

func New(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file config with path %s: %w", path, err)
	}

	var cfgJSON configJSON

	if err = json.NewDecoder(f).Decode(&cfgJSON); err != nil {
		return nil, fmt.Errorf("failed to parse cofig with path %s: %w", path, err)
	}

	cfg := bindConfigFromJSON(cfgJSON)
	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	fileConfig := File{
		lastState: cfg,
	}

	if err = fileConfig.initWatcher(path); err != nil {
		return nil, fmt.Errorf("failed to init watcher with path %s: %w", path, err)
	}

	go fileConfig.watch()

	return &fileConfig, nil
}

func (c *File) GetConfig() config.Config {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.lastState
}

func (c *File) AddWatcher(f config.WatchFunc) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.watchers = append(c.watchers, f)
}

func (c *File) initWatcher(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	c.fsWatcher = watcher

	return watcher.Add(path)
}

func (c *File) watch() {
	for {
		select {
		case <-c.fsWatcher.Events:
			// TODO: implement
		}
	}
}
