package triton

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joyent/gocommon/client"
	"github.com/joyent/gosdc/cloudapi"
	"github.com/joyent/gosign/auth"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/template/interpolate"
)

// AccessConfig is for common configuration related to Triton access
type AccessConfig struct {
	Endpoint string `mapstructure:"triton_url"`
	Account  string `mapstructure:"triton_account"`
	KeyID    string `mapstructure:"triton_key_id"`
	KeyPath  string `mapstructure:"triton_key_path"`
}

// Prepare performs basic validation on the AccessConfig
func (c *AccessConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error

	if c.Endpoint == "" {
		// Use Joyent public cloud as the default endpoint if none is in environment
		c.Endpoint = "https://us-east-1.api.joyent.com"
	}

	if c.Account == "" {
		errs = append(errs, fmt.Errorf("triton_account is required to use the triton builder"))
	}

	if c.KeyID == "" {
		errs = append(errs, fmt.Errorf("triton_key_id is required to use the triton builder"))
	}

	if c.KeyPath == "" {
		errs = append(errs, fmt.Errorf("triton_key_path is required to use the triton builder"))
	}

	if _, err := os.Stat(c.KeyPath); err != nil {
		errs = append(errs, fmt.Errorf("Error reading private key: %s", err))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// CreateTritonClient returns an SDC client configured with the appropriate client credentials
// or an error if creating the client fails.
func (c *AccessConfig) CreateTritonClient() (*cloudapi.Client, error) {
	keyData, err := ioutil.ReadFile(c.KeyPath)
	if err != nil {
		return nil, err
	}

	userauth, err := auth.NewAuth(c.Account, string(keyData), "rsa-sha256")
	if err != nil {
		return nil, err
	}

	creds := &auth.Credentials{
		UserAuthentication: userauth,
		SdcKeyId:           c.KeyID,
		SdcEndpoint:        auth.Endpoint{URL: c.Endpoint},
	}

	return cloudapi.New(client.NewClient(
		c.Endpoint,
		cloudapi.DefaultAPIVersion,
		creds,
		&cloudapi.Logger,
	)), nil
}

func (c *AccessConfig) Comm() communicator.Config {
	return communicator.Config{}
}
