// Copyright (c) 2016 Epracom Advies.
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

package service

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Seal seals the vault
func (s *VaultService) Seal() error {
	if err := s.vaultClient.Sys().Seal(); err != nil {
		return maskAny(err)
	}
	return nil
}

// Unseal unseals the vault.
// If calls a given process several times to obtain unseal keys
func (s *VaultService) Unseal(keyCmd []string) error {
	sys := s.vaultClient.Sys()
	status, err := sys.ResetUnsealProcess()
	if err != nil {
		return maskAny(err)
	}
	if !status.Sealed {
		// Already unsealed
		s.log.Info("Already unsealed")
		return nil
	}
	keyNrs := []int{}
	for i := 1; i <= status.N; i++ {
		keyNrs = append(keyNrs, i)
	}
	shuffleInts(keyNrs)
	replace := func(s string, keyNr int) (string, error) {
		tmpl, err := template.New(fmt.Sprintf("keyTmpl")).Parse(s)
		if err != nil {
			return "", maskAny(err)
		}
		data := struct {
			Key int
		}{
			Key: keyNr,
		}
		buffer := &bytes.Buffer{}
		if err := tmpl.Execute(buffer, data); err != nil {
			return "", maskAny(err)
		}
		return buffer.String(), nil
	}
	for _, keyNr := range keyNrs {
		s.log.Debug("Fetching key %d", keyNr)
		path, err := replace(keyCmd[0], keyNr)
		if err != nil {
			return maskAny(err)
		}
		cmd := exec.Command(path)
		for _, rawArg := range keyCmd[1:] {
			arg, err := replace(rawArg, keyNr)
			if err != nil {
				return maskAny(err)
			}
			cmd.Args = append(cmd.Args, arg)
		}
		rawKey, err := cmd.Output()
		if err != nil {
			return maskAny(err)
		}
		key := strings.TrimSpace(string(rawKey))
		if key == "" {
			return maskAny(fmt.Errorf("Result from $(%s %s) is empty", path, strings.Join(cmd.Args, " ")))
		}
		if status, err := sys.Unseal(key); err != nil {
			return maskAny(err)
		} else if !status.Sealed {
			s.log.Info("Vault is now unsealed")
			return nil
		}
	}
	return nil
}

func shuffleInts(slc []int) {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
}
