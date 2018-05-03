package logcache

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	api "github.com/alekssaul/logcache-operator/pkg/apis/app/v1alpha1"
	"github.com/alekssaul/logcache-operator/pkg/tls"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultClusterDomain = "cluster.local"
	orgForTLSCert        = []string{"pivotal.io"}
)

// prepareDefaultLogcacheTLSSecrets creates the default secrets for the Log cache TLS assets.
// Currently we self-generate the CA, and use the self generated CA to sign all the TLS certs.
func prepareDefaultLogcacheTLSSecrets(l *api.LogCache) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("prepare default vault TLS secrets failed: %v", err)
		}
	}()

	// if TLS spec doesn't exist or secrets doesn't exist, then we can go create secrets.
	if api.IsTLSConfigured(l.Spec.TLS) {
		se := &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      l.Spec.TLS.Static.ServerSecret,
				Namespace: l.Namespace,
			},
		}
		err = query.Get(se)
		if err == nil {
			return nil
		}
		if !apierrors.IsNotFound(err) {
			return err
		}
	}

	/*caKey, caCrt, err := newCACert()
	if err != nil {
		return err
	}*/

	return err
}

func newCACert() (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	config := tls.CertConfig{
		CommonName:   "Log Cache Operator CA",
		Organization: orgForTLSCert,
	}

	cert, err := tls.NewSelfSignedCACertificate(config, key)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, err
}

func newKeyAndCert(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, config tls.CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	cert, err := tls.NewSignedCertificate(config, key, caCert, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, nil
}

// newTLSSecret is a common utility for creating a secret containing TLS assets.
func newTLSSecret(l *api.LogCache, caKey *rsa.PrivateKey, caCrt *x509.Certificate, commonName, secretName string,
	addrs []string, fieldMap map[string]string) (*v1.Secret, error) {
	tc := tls.CertConfig{
		CommonName:   commonName,
		Organization: orgForTLSCert,
		AltNames:     tls.NewAltNames(addrs),
	}
	key, crt, err := newKeyAndCert(caCrt, caKey, tc)
	if err != nil {
		return nil, fmt.Errorf("new TLS secret failed: %v", err)
	}
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: l.Namespace,
			Labels:    labelsForLogCache(l.Name),
		},
		Data: map[string][]byte{
			fieldMap["key"]:  tls.EncodePrivateKeyPEM(key),
			fieldMap["cert"]: tls.EncodeCertificatePEM(crt),
			fieldMap["ca"]:   tls.EncodeCertificatePEM(caCrt),
		},
	}
	return secret, nil
}
