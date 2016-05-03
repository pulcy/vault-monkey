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
	"crypto/x509"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/hashicorp/vault/api"
	"github.com/juju/errgo"
	"github.com/op/go-logging"
)

type VaultServiceConfig struct {
	VaultAddr   string // URL of the vault
	VaultCACert string // Path to a PEM-encoded CA cert file to use to verify the Vault server SSL certificate
	VaultCAPath string // Path to a directory of PEM-encoded CA cert files to verify the Vault server SSL certificate
	TokenPath   string // Path of a file containing the login token
}

type VaultService struct {
	log    *logging.Logger
	config api.Config
	token  string
}

type VaultClient struct {
	Address string
	Client  *api.Client
}

// NewVaultService creates a new VaultService and loads its configuration from the given settings.
func NewVaultService(log *logging.Logger, srvCfg VaultServiceConfig) (*VaultService, error) {
	// Create a vault client
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, maskAny(err)
	}
	var serverName string
	if srvCfg.VaultAddr != "" {
		log.Debugf("Setting vault address to %s", srvCfg.VaultAddr)
		config.Address = srvCfg.VaultAddr
		url, err := url.Parse(config.Address)
		if err != nil {
			return nil, maskAny(err)
		}
		host, _, err := net.SplitHostPort(url.Host)
		if err != nil {
			return nil, maskAny(err)
		}
		serverName = host
	}
	if srvCfg.VaultCACert != "" || srvCfg.VaultCAPath != "" {
		var newCertPool *x509.CertPool
		var err error
		if srvCfg.VaultCACert != "" {
			log.Debugf("Loading CA cert: %s", srvCfg.VaultCACert)
			newCertPool, err = api.LoadCACert(srvCfg.VaultCACert)
		} else {
			log.Debugf("Loading CA certs from: %s", srvCfg.VaultCAPath)
			newCertPool, err = api.LoadCAPath(srvCfg.VaultCAPath)
		}
		if err != nil {
			return nil, maskAny(err)
		}
		clientTLSConfig := config.HttpClient.Transport.(*http.Transport).TLSClientConfig
		clientTLSConfig.RootCAs = newCertPool
		clientTLSConfig.ServerName = serverName
	}
	var token string
	if srvCfg.TokenPath != "" {
		log.Debugf("Loading token from %s", srvCfg.TokenPath)
		var err error
		token, err = readID(srvCfg.TokenPath)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	return &VaultService{
		log:    log,
		config: *config,
		token:  token,
	}, nil
}

// newUnsealedClient creates the first single vault client that resolves to an unsealed vault instance.
func (s *VaultService) newUnsealedClient() (*api.Client, error) {
	clients, err := s.newClients()
	if err != nil {
		return nil, maskAny(err)
	}
	for _, client := range clients {
		status, err := client.Client.Sys().SealStatus()
		if err != nil {
			s.log.Warningf("Vault at %s cannot be reached: %#v", client.Address, err)
		} else if status.Sealed {
			s.log.Warningf("Vault at %s is sealed", client.Address)
		} else {
			return client.Client, nil
		}
	}
	return nil, maskAny(errgo.WithCausef(nil, VaultError, "no unsealed vault instance found"))
}

// newClients resolves the configured vault address into IP addresses and creates a one vault client
// for each IP address.
func (s *VaultService) newClients() ([]VaultClient, error) {
	url, err := url.Parse(s.config.Address)
	if err != nil {
		return nil, maskAny(err)
	}
	host, port, err := net.SplitHostPort(url.Host)
	if err != nil {
		return nil, maskAny(err)
	}
	// Is the host address already an IP address?
	ip := net.ParseIP(host)
	if ip != nil {
		// Yes, host address is an IP
		client, err := newClientFromConfig(s.config, s.token)
		if err != nil {
			return nil, maskAny(err)
		}
		return []VaultClient{VaultClient{Client: client, Address: s.config.Address}}, nil
	}

	// Get IP's for host address
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, maskAny(err)
	}

	// Create a client for each IP
	list := []VaultClient{}
	for _, ip := range ips {
		ipURL := *url
		ipURL.Host = net.JoinHostPort(ip.String(), port)
		config := s.config
		config.Address = ipURL.String()
		client, err := newClientFromConfig(config, s.token)
		if err != nil {
			return nil, maskAny(err)
		}
		list = append(list, VaultClient{Client: client, Address: config.Address})
	}

	return list, nil
}

func newClientFromConfig(config api.Config, token string) (*api.Client, error) {
	client, err := api.NewClient(&config)
	if err != nil {
		return nil, maskAny(err)
	}
	if token != "" {
		client.SetToken(token)
	}
	return client, nil
}

// asyncForEachClient creates a new vault client for each IP address and calls the given function
// for each client (asynchronous).
func (s *VaultService) asyncForEachClient(f func(client VaultClient) error) error {
	clients, err := s.newClients()
	if err != nil {
		return maskAny(err)
	}
	wg := sync.WaitGroup{}
	errors := make(chan error, len(clients))
	for _, client := range clients {
		wg.Add(1)
		go func(client VaultClient) {
			defer wg.Done()
			if err := f(client); err != nil {
				errors <- maskAny(err)
			}
		}(client)
	}
	wg.Wait()
	close(errors)

	// Gather errors
	if err := collectErrorsFromChannel(errors); err != nil {
		return maskAny(err)
	}

	return nil
}
