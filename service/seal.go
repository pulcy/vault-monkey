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

package service

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/ryanuber/columnize"
)

var (
	unsealFetchMutex sync.Mutex
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// SealStatus shows the seal status of the vault
func (s *VaultService) SealStatus() error {
	m := sync.Mutex{}
	lines := []string{
		"Address | Status | Key Threshold | Key Shares | Unseal Progress",
	}
	seal := func(client VaultClient) error {
		s.log.Debugf("Fetching seal status of vault at %s", client.Address)
		status, err := client.Client.Sys().SealStatus()
		var line string
		if err != nil {
			line = fmt.Sprintf("%s | error: %s | - | - | -", client.Address, err.Error())
		} else {
			statusText := "unsealed"
			if status.Sealed {
				statusText = "sealed"
			}
			line = fmt.Sprintf("%s | %s | %d | %d | %d", client.Address, statusText, status.T, status.N, status.Progress)
		}
		m.Lock()
		defer m.Unlock()
		lines = append(lines, line)
		return nil
	}
	if err := s.asyncForEachClient(seal); err != nil {
		return maskAny(err)
	}
	fmt.Println(columnize.SimpleFormat(lines))
	return nil
}

// Seal seals the vault
func (s *VaultService) Seal() error {
	seal := func(client VaultClient) error {
		s.log.Debugf("Sealing vault at %s", client.Address)
		if err := client.Client.Sys().Seal(); err != nil {
			return maskAny(err)
		}
		s.log.Infof("Vault at %s is now sealed", client.Address)
		return nil
	}
	if err := s.asyncForEachClient(seal); err != nil {
		return maskAny(err)
	}
	return nil
}

// Unseal unseals the vault.
// If calls a given process several times to obtain unseal keys
func (s *VaultService) Unseal(keyCmd []string) error {
	unseal := func(client VaultClient) error {
		if err := s.unseal(client, keyCmd); err != nil {
			return maskAny(err)
		}
		return nil
	}
	if err := s.asyncForEachClient(unseal); err != nil {
		return maskAny(err)
	}
	return nil
}

func (s *VaultService) unseal(client VaultClient, keyCmd []string) error {
	sys := client.Client.Sys()
	s.log.Debugf("Unsealing vault at %s", client.Address)
	status, err := sys.ResetUnsealProcess()
	if err != nil {
		return maskAny(err)
	}
	if !status.Sealed {
		// Already unsealed
		s.log.Infof("Vault at %s is already unsealed", client.Address)
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
	fetchKey := func(keyNr int) (string, error) {
		// Call key command synchronized to have a smooth user experience
		unsealFetchMutex.Lock()
		defer unsealFetchMutex.Unlock()

		s.log.Debugf("Fetching key %d", keyNr)
		path, err := replace(keyCmd[0], keyNr)
		if err != nil {
			return "", maskAny(err)
		}
		cmd := exec.Command(path)
		for _, rawArg := range keyCmd[1:] {
			arg, err := replace(rawArg, keyNr)
			if err != nil {
				return "", maskAny(err)
			}
			cmd.Args = append(cmd.Args, arg)
		}
		rawKey, err := cmd.Output()
		if err != nil {
			return "", maskAny(err)
		}
		key := strings.TrimSpace(string(rawKey))
		if key == "" {
			return "", maskAny(fmt.Errorf("Result from $(%s %s) is empty", path, strings.Join(cmd.Args, " ")))
		}
		return key, nil
	}

	for _, keyNr := range keyNrs {
		key, err := fetchKey(keyNr)
		if err != nil {
			return maskAny(err)
		}
		if status, err := sys.Unseal(key); err != nil {
			return maskAny(err)
		} else if !status.Sealed {
			s.log.Infof("Vault at %s is now unsealed", client.Address)
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
