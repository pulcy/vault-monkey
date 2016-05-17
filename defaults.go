// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const (
	defaultGithubTokenPathTmpl = "~/.pulcy/github-token"
)

func defaultGithubToken() string {
	path, err := homedir.Expand(defaultGithubTokenPathTmpl)
	if err != nil {
		log.Warningf("Cannot expand %s: %#v", defaultGithubTokenPathTmpl, err)
		return ""
	}
	content, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Warningf("Cannot read %s: %#v", path, err)
		return ""
	}
	return strings.TrimSpace(string(content))
}
