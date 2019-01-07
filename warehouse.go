// Package fswarehouse provides a sekrat.Warehouse implemntation that saves data
// to the local filesystem.
package fswarehouse

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sekrat/sekrat"
	"github.com/spf13/afero"
)

// Warehouse is a sekrat.Warehouse that saves data to the local filesystem.
type Warehouse struct {
	BaseDir string
	fs      fs
}

// IDs returns the array of all secret IDs known to the Warehouse.
func (warehouse *Warehouse) IDs() []string {
	warehouse.setup()

	keys := make([]string, 0)

	Walk(
		warehouse.BaseDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				keys = append(keys, strings.Replace(path, warehouse.BaseDir+"/", "", 1))
			}

			return nil
		},
	)

	return keys
}

// Store takes a secret ID and a chunk of data, saves the data indexed by the
// ID, and returns an error. If there are problems along the way, the error is
// populated. Otherwise, the error is nil.
//
// Note: As the data saved by this Warehouse is encoded via base64, it is
// decoded before being returned.
func (warehouse *Warehouse) Store(id string, data []byte) error {
	warehouse.setup()

	abs := filepath.Join(warehouse.BaseDir, id)

	err := CreateDir(filepath.Dir(abs), 0755)
	if err != nil {
		return errors.New("could not write")
	}

	encoded := []byte(base64.StdEncoding.EncodeToString(data))

	err = ioutil.WriteFile(abs, encoded, 0644)
	if err != nil {
		return errors.New("could not write")
	}

	return nil
}

// Retrieve takes a secret id and returns the data for that secret and an error.
// If there are problems along the way, the data is nil and the error is
// populated. Otherwise, the data is populated and the error is nil.
//
// Note: To avoid encoding issues upon load, the data is encoded via base64
// before it is saved.
func (warehouse *Warehouse) Retrieve(id string) ([]byte, error) {
	warehouse.setup()

	abs := filepath.Join(warehouse.BaseDir, id)
	if !FileExists(abs) {
		return nil, errors.New("not found")
	}

	data, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, errors.New("could not read file")
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, errors.New("could not decode data")
	}

	return decoded, nil
}

func (warehouse *Warehouse) setup() {
	CreateDir(warehouse.BaseDir, 0755)
}

// New takes a base path and returns a Warehouse (as a sekrat.Warehouse).
func New(baseDir string) sekrat.Warehouse {
	abs, err := filepath.Abs(baseDir)
	if err == nil {
		baseDir = abs
	}

	return &Warehouse{BaseDir: baseDir, fs: &real{fs: afero.NewOsFs()}}
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
