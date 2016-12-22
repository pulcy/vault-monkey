package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	http "net/http"
	"net/url"
	"os"

	k8s "github.com/YakLabs/k8s-client"
	"github.com/pkg/errors"
)

//go:generate ./make-type HorizontalPodAutoscaler autoscaling/v1
//go:generate ./make-type Secret v1
//go:generate ./make-type Deployment extensions/v1beta1
//go:generate ./make-type Pod v1
//go:generate ./make-type ConfigMap v1
//go:generate ./make-type ReplicaSet extensions/v1beta1
//go:generate ./make-type Service v1
//go:generate ./make-type ServiceAccount v1

const (
	tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	caFile    = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

type (
	// Client is an http client for kubernetes
	Client struct {
		username           string
		password           string
		token              string
		authHeader         string
		server             string
		certPool           *x509.CertPool
		clientCert         []byte
		clientKey          []byte
		insecureSkipVerify bool
		client             *http.Client
	}

	// OptionsFunc is a function passed to new for setting options on a new client.
	OptionsFunc func(*Client) error
)

// New creates a new client.
func New(options ...OptionsFunc) (*Client, error) {
	c := &Client{}
	for _, f := range options {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	if c.client == nil {
		// TODO: allow passing in CA to use.
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.insecureSkipVerify,
			},
		}

		if c.certPool != nil {
			tr.TLSClientConfig.RootCAs = c.certPool
		}

		if c.clientCert != nil && c.clientKey != nil {
			cert, err := tls.X509KeyPair(c.clientCert, c.clientKey)
			if err != nil {
				return nil, errors.Wrap(err, "X509KeyPair failed")
			}
			tr.TLSClientConfig.Certificates = []tls.Certificate{cert}
			tr.TLSClientConfig.BuildNameToCertificate()
		}

		c.client = &http.Client{
			Transport: tr,
		}
	}

	if c.token != "" {
		c.authHeader = "Bearer " + c.token
	} else {
		if c.username != "" {
			c.authHeader = base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
		}
	}

	return c, nil
}

// NewInCluster creates a cluster suitable for use within the kubernetes cluster.
func NewInCluster() (*Client, error) {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	if host == "" {
		return nil, errors.New("KUBERNETES_SERVICE_HOST is not set")
	}
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	if port == "" {
		return nil, errors.New("KUBERNETES_SERVICE_PORT is not set")
	}
	server := fmt.Sprintf("https://%s:%s", host, port)

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read token file: "+tokenFile)
	}

	return New(SetToken(string(token)), SetServer(server), SetCAFromFile(caFile))
}

// SetUsername sets the username to be used for authentication.
func SetUsername(username string) func(*Client) error {
	return func(c *Client) error {
		c.username = username
		return nil
	}
}

// SetPassword sets the password to be used for authentication.
func SetPassword(password string) func(*Client) error {
	return func(c *Client) error {
		c.password = password
		return nil
	}
}

// SetToken sets the token to be used for authentication. If this is set, it
// takes precedence.
func SetToken(token string) func(*Client) error {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}

// SetServer sets the target API server
func SetServer(server string) func(*Client) error {
	return func(c *Client) error {
		// make sure its valid
		if _, err := url.Parse(server); err != nil {
			return errors.Wrap(err, "failed to parse server")
		}
		c.server = server
		return nil
	}
}

// SetClient allows the caller to specify a custom http.Client to use
func SetClient(client *http.Client) func(*Client) error {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

// SetCA sets a CA to verify the servers certificate
func SetCA(cert []byte) func(*Client) error {
	return func(c *Client) error {
		c.certPool = x509.NewCertPool()
		if ok := c.certPool.AppendCertsFromPEM(cert); !ok {
			return errors.New("AppendCertsFromPEM failed")
		}
		return nil
	}
}

// SetCA sets a CA to verify the servers certificate from a file
func SetCAFromFile(path string) func(*Client) error {
	return func(c *Client) error {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		return SetCA(data)(c)
	}
}

// SetClientCert sets the certificate to be used for authentication.
func SetClientCert(cert []byte) func(*Client) error {
	return func(c *Client) error {
		c.clientCert = cert
		return nil
	}
}

// SetClientCertFromFile sets the certificate to be used for authentication from a file.
func SetClientCertFromFile(path string) func(*Client) error {
	return func(c *Client) error {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		return SetClientCert(data)(c)
	}
}

// SetClientKey sets the key to be used for authentication.
func SetClientKey(key []byte) func(*Client) error {
	return func(c *Client) error {
		c.clientKey = key
		return nil
	}
}

// SetClientKeyFromFile sets the key to be used for authentication from a file.
func SetClientKeyFromFile(path string) func(*Client) error {
	return func(c *Client) error {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		return SetClientKey(data)(c)
	}
}

// SetInsecureSkipVerify allows the caller to skip verification of the servers cert
func SetInsecureSkipVerify(skip bool) func(*Client) error {
	return func(c *Client) error {
		c.insecureSkipVerify = skip
		return nil
	}
}
func (c *Client) newRequest(method, path string, v interface{}) (*http.Request, error) {
	if v != nil {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(method, c.server+path, bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
		if c.authHeader != "" {
			req.Header.Add("Authorization", c.authHeader)
		}
		return req, nil
	}
	// weirdness with passing nil interface, so just copy/paste
	req, err := http.NewRequest(method, c.server+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	return req, nil
}

func readStatus(body io.Reader) (*k8s.Status, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response body")
	}
	var out k8s.Status
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal status")
	}
	return &out, nil
}

func (c *Client) do(method, path string, in interface{}, out interface{}, codes ...int) (int, error) {
	req, err := c.newRequest(method, path, in)
	if err != nil {
		return 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}

	// make errcheck happy
	defer func() {
		_ = resp.Body.Close()
	}()

	if len(codes) == 0 {
		codes = []int{
			200,
		}
	}

	found := false
	for _, i := range codes {
		if i == resp.StatusCode {
			found = true
			break
		}
	}

	if !found {
		status, err := readStatus(resp.Body)
		if err != nil {
			return resp.StatusCode, errors.Wrapf(err, "unable to read status: %d", resp.StatusCode)
		}
		return resp.StatusCode, status
	}

	if out != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, errors.Wrap(err, "failed to read response body")
		}
		if err := json.Unmarshal(body, out); err != nil {
			return resp.StatusCode, err
		}
	}
	return resp.StatusCode, nil
}

func listOptionsQuery(opts *k8s.ListOptions) string {
	if opts == nil {
		return ""
	}
	val := url.Values{}
	if opts.LabelSelector.MatchLabels != nil && len(opts.LabelSelector.MatchLabels) > 0 {
		labels := url.Values{}
		for k, v := range opts.LabelSelector.MatchLabels {
			labels.Set(k, v)
		}
		val.Set("labelSelector", labels.Encode())
	}
	if opts.FieldSelector != nil && len(opts.FieldSelector) > 0 {
		fields := url.Values{}
		for k, v := range opts.FieldSelector {
			fields.Set(k, v)
		}
		val.Set("fieldSelector", fields.Encode())
	}

	return val.Encode()
}
