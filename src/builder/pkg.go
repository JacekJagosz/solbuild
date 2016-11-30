//
// Copyright © 2016 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package builder

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Package is the main item we deal with, avoiding the internals
type Package struct {
	Name    string
	Version string
	Release int
}

// YmlPackage is a parsed ypkg build file
type YmlPackage struct {
	Name    string
	Version string
	Release int
}

// XMLUpdate represents an update in the package history
type XMLUpdate struct {
	Release int `xml:"release,attr"`
	Date    string
	Version string
	Comment string
	Name    string
	Email   string
}

// XMLSource is the actual source info for each pspec.xml
type XMLSource struct {
	Homepage string
	Name     string
}

// XMLPackage contains all of the pspec.xml metadata
type XMLPackage struct {
	Name    string
	Source  XMLSource
	History []XMLUpdate `xml:"History>Update"`
}

// NewPackage will attempt to parse the given path, and return a new Package
// instance if this succeeds.
func NewPackage(path string) (*Package, error) {
	if strings.HasSuffix(path, ".xml") {
		return NewXMLPackage(path)
	}
	return NewYmlPackage(path)
}

// NewXMLPackage will attempt to parse the pspec.xml file @ path
func NewXMLPackage(path string) (*Package, error) {
	var by []byte
	var err error
	var fi *os.File

	fi, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	by, err = ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}
	xpkg := &XMLPackage{}
	if err = xml.Unmarshal(by, xpkg); err != nil {
		return nil, err
	}
	if len(xpkg.History) < 1 {
		return nil, errors.New("xml: Malformed pspec file")
	}

	upd := xpkg.History[0]
	ret := &Package{
		Name:    strings.TrimSpace(xpkg.Source.Name),
		Version: strings.TrimSpace(upd.Version),
		Release: upd.Release,
	}

	if ret.Name == "" {
		return nil, errors.New("xml: Missing name in package")
	}
	if ret.Version == "" {
		return nil, errors.New("xml: Missing version in package")
	}
	if ret.Release < 0 {
		return nil, fmt.Errorf("xml: Invalid release in package: %d", ret.Release)
	}
	return ret, nil
}

// NewYmlPackage will attempt to parse the ypkg package.yml file @ path
func NewYmlPackage(path string) (*Package, error) {
	return nil, errors.New("ypkg: Not yet implemented")
}
