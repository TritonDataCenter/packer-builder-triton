package triton

// AccessConfig is for common configuration related to Triton access
type AccessConfig struct {
	Endpoint string `mapstructure:"sdc_url"`
	Account  string `mapstructure:"sdc_account"`
	KeyID    string `mapstructure:"sdc_key_id"`
	KeyPath  string `mapstructure:"sdc_key_path"`
}
