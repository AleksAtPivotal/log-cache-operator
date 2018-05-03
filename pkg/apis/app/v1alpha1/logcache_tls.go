package v1alpha1

const (
	// CATLSCertName Name of CA cert file in the client secret
	CATLSCertName = "ca.crt"
)

// TLSPolicy defines the TLS policy of the Log Cache
type TLSPolicy struct {
	// StaticTLS enables user to use static x509 certificates and keys,
	// by putting them into Kubernetes secrets, and specifying them here.
	// If this is not set, operator will auto-gen TLS assets and secrets.
	Static *StaticTLS `json:"static,omitempty"`
}

// StaticTLS defines the TLS objects in Kubernetes
type StaticTLS struct {
	ServerSecret string `json:"serverSecret,omitempty"`
	// ClientSecret is the secret containing the CA certificate
	// that will be used to verify the above server certificate
	ClientSecret string `json:"clientSecret,omitempty"`
}

// IsTLSConfigured checks if the vault TLS secrets have been specified by the user
func IsTLSConfigured(tp *TLSPolicy) bool {
	if tp == nil || tp.Static == nil {
		return false
	}
	return len(tp.Static.ServerSecret) != 0 && len(tp.Static.ClientSecret) != 0
}
