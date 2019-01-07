package fswarehouse

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type fs struct {
	root afero.Fs
}

func (fs *fs) CreateDir(path string, mode os.FileMode) error {
	if !fs.fileExists(path) {
		err := fs.root.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *fs) fileExists(path string) bool {
	_, err := fs.root.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (fs *fs) isDir(path string) bool {
	info, err := fs.root.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func (fs *fs) walk(path string, walkFunc filepath.WalkFunc) error {
	return afero.Walk(fs.root, path, walkFunc)
}

func newfs() *fs {
	return &fs{root: afero.NewOsFs()}
}

/*
Copyright 2019 Dennis Walters

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
