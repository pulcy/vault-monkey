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
	log          *logging.Logger
	address      string
	serverName   string
	initialToken string
	certPool     *x509.CertPool
}

type VaultClient struct {
	Address string
	Client  *api.Client
}

// NewVaultService creates a new VaultService and loads its configuration from the given settings.
func NewVaultService(log *logging.Logger, srvCfg VaultServiceConfig) (*VaultService, error) {
	// Create a vault client
	var serverName, address string
	if srvCfg.VaultAddr != "" {
		address = srvCfg.VaultAddr
		log.Debugf("Setting vault address to %s", address)
		url, err := url.Parse(address)
		if err != nil {
			return nil, maskAny(err)
		}
		host, _, err := net.SplitHostPort(url.Host)
		if err != nil {
			return nil, maskAny(err)
		}
		serverName = host
	}
	var newCertPool *x509.CertPool
	if srvCfg.VaultCACert != "" || srvCfg.VaultCAPath != "" {
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
		log:          log,
		address:      address,
		serverName:   serverName,
		initialToken: token,
		certPool:     newCertPool,
	}, nil
}

func (s *VaultService) newConfig() (*api.Config, error) {
	// Create a vault client
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, maskAny(err)
	}
	if s.address != "" {
		config.Address = s.address
	}
	if s.certPool != nil {
		clientTLSConfig := config.HttpClient.Transport.(*http.Transport).TLSClientConfig
		clientTLSConfig.RootCAs = s.certPool
		clientTLSConfig.ServerName = s.serverName
	}
	return config, nil
}

// newUnsealedClient creates the first single vault client that resolves to an unsealed vault instance.
func (s *VaultService) newUnsealedClient() (*api.Client, string, error) {
	clients, err := s.newClients()
	if err != nil {
		return nil, "", maskAny(err)
	}
	for _, client := range clients {
		// Check seal status
		status, err := client.Client.Sys().SealStatus()
		if err != nil {
			s.log.Debugf("vault at %s cannot be reached: %s", client.Address, Describe(err))
			continue
		} else if status.Sealed {
			s.log.Warningf("Vault at %s is sealed", client.Address)
			continue
		}

		// Check leader status
		resp, err := client.Client.Sys().Leader()
		if err != nil {
			s.log.Debugf("vault at %s cannot be reached: %s", client.Address, Describe(err))
			continue
		} else if resp.HAEnabled && !resp.IsSelf {
			s.log.Debugf("vault at %s is not the leader", client.Address)
			continue
		}

		s.log.Debugf("found unsealed vault client at %s", client.Address)
		return client.Client, client.Address, nil
	}
	return nil, "", maskAny(errgo.WithCausef(nil, VaultError, "no unsealed vault instance found"))
}

// newClients resolves the configured vault address into IP addresses and creates a one vault client
// for each IP address.
func (s *VaultService) newClients() ([]VaultClient, error) {
	url, err := url.Parse(s.address)
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
		config, err := s.newConfig()
		if err != nil {
			return nil, maskAny(err)
		}
		s.log.Debugf("fixed vault client at %s", config.Address)
		client, err := newClientFromConfig(config, s.initialToken)
		if err != nil {
			return nil, maskAny(err)
		}
		return []VaultClient{VaultClient{Client: client, Address: config.Address}}, nil
	}

	// Get IP's for host address
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, maskAny(err)
	}

	// Create a client for each IP
	list := []VaultClient{}
	for j := 0; j < 2; j++ {
		preferIPv6 := j == 0
		for _, ip := range ips {
			isIPv6 := ip.To4() == nil
			if preferIPv6 == isIPv6 {
				ipURL := *url
				ipURL.Host = net.JoinHostPort(ip.String(), port)
				config, err := s.newConfig()
				if err != nil {
					return nil, maskAny(err)
				}
				config.Address = ipURL.String()
				client, err := newClientFromConfig(config, s.initialToken)
				if err != nil {
					return nil, maskAny(err)
				}
				list = append(list, VaultClient{Client: client, Address: config.Address})
				s.log.Debugf("possible vault client at %s", config.Address)
			}
		}
	}

	return list, nil
}

func newClientFromConfig(config *api.Config, token string) (*api.Client, error) {
	client, err := api.NewClient(config)
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

func (s *VaultService) newAuthenticatedClient(vaultClient *api.Client) *AuthenticatedVaultClient {
	return &AuthenticatedVaultClient{log: s.log, vaultClient: vaultClient}
}
