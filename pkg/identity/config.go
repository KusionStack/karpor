package identity

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

const (
	EnvUnifiedIdentityCertFile = "UNIFIED_IDENTITY_CERT_FILE"
	EnvUnifiedIdentityKeyFile  = "UNIFIED_IDENTITY_KEY_FILE"
)

var (
	unifiedIdentityCertFile = "/var/run/secrets/kubernetes.io/serviceaccount/app.crt"
	unifiedIdentityKeyFile  = "/var/run/secrets/kubernetes.io/serviceaccount/app.key"
)

func AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&unifiedIdentityCertFile, "unified-identity-cert-file", unifiedIdentityCertFile, "the filepath of unified-identity certificate")
	flags.StringVar(&unifiedIdentityKeyFile, "unified-identity-key-file", unifiedIdentityKeyFile, "the filepath of unified-identity private key")
}

func Validate() error {
	if len(unifiedIdentityCertFile) == 0 {
		return fmt.Errorf("--unified-identity-cert-file required")
	}

	if len(unifiedIdentityKeyFile) == 0 {
		return fmt.Errorf("--unified-identity-key-file required")
	}

	return nil
}

func GetCertFile() string {
	if unifiedIdentityCertFile != "" {
		return unifiedIdentityCertFile
	}
	return os.Getenv(EnvUnifiedIdentityCertFile)
}

func GetKeyFile() string {
	if unifiedIdentityKeyFile != "" {
		return unifiedIdentityKeyFile
	}
	return os.Getenv(EnvUnifiedIdentityKeyFile)
}
