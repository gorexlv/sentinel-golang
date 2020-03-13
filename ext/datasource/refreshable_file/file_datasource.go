package refreshable_file

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

var (
	logger = logging.GetDefaultLogger()
)

type RefreshableFileDataSource struct {
	datasource.Base
	sourceFilePath string
	once           sync.Once
}

func FileDataSourceStarter(sourceFilePath string) *RefreshableFileDataSource {
	ds := &RefreshableFileDataSource{
		sourceFilePath: sourceFilePath,
	}
	return ds
}

func (s *RefreshableFileDataSource) ReadSource() ([]byte, error) {
	f, err := os.Open(s.sourceFilePath)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		return nil, errors.Errorf("The rules file is not existed, err:%+v.", errors.WithStack(err))
	}
	return ioutil.ReadAll(f)
}

func (s *RefreshableFileDataSource) Initialize() {
	s.doUpdate()
	// start watcher
	s.once.Do(
		func() {
			go func() {
				watcher, err := fsnotify.NewWatcher()
				defer watcher.Close()

				if err != nil {
					panic(fmt.Sprintf("Fail to new a watcher of fsnotify, err:%+v", err))
				}
				err = watcher.Add(s.sourceFilePath)
				if err != nil {
					panic(fmt.Sprintf("Fail add a watcher on file(%s), err:%+v", s.sourceFilePath, err))
				}

				for {
					select {
					case ev := <-watcher.Events:
						if ev.Op&fsnotify.Write == fsnotify.Write {
							s.doUpdate()
						}

						if ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Rename == fsnotify.Rename {
							logger.Errorf("The file source(%s) was removed or renamed.", s.sourceFilePath)
							// todo remove
							return
						}
					case err := <-watcher.Errors:
						logger.Errorf("Watch err on file(%s), err:", s.sourceFilePath, err)
						time.Sleep(time.Second * 3)
					}
				}
			}()
		})
}

func (s *RefreshableFileDataSource) doUpdate() {
	src, err := s.ReadSource()
	if err != nil {
		logger.Errorf("read source: %+v", err)
		return
	}

	if err := s.Handle(src); err != nil {
		logger.Errorf("handle source: %+v", err)
	}
}

func (s *RefreshableFileDataSource) Close() error {
	return nil
}
