package oracle

// Config contains the configuration parameters of the oracle.
type Config struct {
	// KeyFile has to be a path to a PEM format encoded secp256k1 key
	KeyFile     string `configkey:"oracle.key_file" validate:"required"`
	KeyPassFile string `configkey:"oracle.key_pass.file"`
	KeyPass     string `configkey:"oracle.key_pass"`
}
