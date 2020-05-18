package oracle

// Config contains the configuration parameters of the oracle.
type Config struct {
	// KeyFile has to be a path to a PEM format encoded secp256k1 key
	KeyFile     string `configkey:"oracle.keyFile" validate:"required"`
	KeyPassFile string `configkey:"oracle.keyPass.file"`
	KeyPass     string `configkey:"oracle.keyPass"`
}
