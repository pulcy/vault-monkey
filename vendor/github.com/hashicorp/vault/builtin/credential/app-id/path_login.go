package appId

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"app_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The unique app ID",
			},

			"user_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The unique user ID",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	appId := data.Get("app_id").(string)
	userId := data.Get("user_id").(string)

	// Ensure both appId and userId are provided
	if appId == "" || userId == "" {
		return logical.ErrorResponse("missing 'app_id' or 'user_id'"), nil
	}

	// Look up the apps that this user is allowed to access
	appsMap, err := b.MapUserId.Get(req.Storage, userId)
	if err != nil {
		return nil, err
	}
	if appsMap == nil {
		return logical.ErrorResponse("invalid user ID or app ID"), nil
	}

	// If there is a CIDR block restriction, check that
	if raw, ok := appsMap["cidr_block"]; ok {
		_, cidr, err := net.ParseCIDR(raw.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid restriction cidr: %s", err)
		}

		var addr string
		if req.Connection != nil {
			addr = req.Connection.RemoteAddr
		}
		if addr == "" || !cidr.Contains(net.ParseIP(addr)) {
			return logical.ErrorResponse("unauthorized source address"), nil
		}
	}

	appsRaw, ok := appsMap["value"]
	if !ok {
		appsRaw = ""
	}

	apps, ok := appsRaw.(string)
	if !ok {
		return nil, fmt.Errorf("internal error: mapping is not a string")
	}

	// Verify that the app is in the list
	found := false
	appIdBytes := []byte(appId)
	for _, app := range strings.Split(apps, ",") {
		match := []byte(strings.TrimSpace(app))
		// Protect against a timing attack with the app_id comparison
		if subtle.ConstantTimeCompare(match, appIdBytes) == 1 {
			found = true
		}
	}
	if !found {
		return logical.ErrorResponse("invalid user ID or app ID"), nil
	}

	// Get the raw data associated with the app
	appRaw, err := b.MapAppId.Get(req.Storage, appId)
	if err != nil {
		return nil, err
	}
	if appRaw == nil {
		return logical.ErrorResponse("invalid user ID or app ID"), nil
	}

	// Get the policies associated with the app
	policies, err := b.MapAppId.Policies(req.Storage, appId)
	if err != nil {
		return nil, err
	}

	// Check if we have a display name
	var displayName string
	if raw, ok := appRaw["display_name"]; ok {
		displayName = raw.(string)
	}

	// Store hashes of the app ID and user ID for the metadata
	appIdHash := sha1.Sum([]byte(appId))
	userIdHash := sha1.Sum([]byte(userId))
	metadata := map[string]string{
		"app-id":  "sha1:" + hex.EncodeToString(appIdHash[:]),
		"user-id": "sha1:" + hex.EncodeToString(userIdHash[:]),
	}

	return &logical.Response{
		Auth: &logical.Auth{
			DisplayName: displayName,
			Policies:    policies,
			Metadata:    metadata,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	return framework.LeaseExtend(0, 0, b.System())(req, d)
}

const pathLoginSyn = `
Log in with an App ID and User ID.
`

const pathLoginDesc = `
This endpoint authenticates using an application ID, user ID and potential the IP address of the connecting client.
`
