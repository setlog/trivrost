package signatures

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/setlog/trivrost/pkg/misc"

	log "github.com/sirupsen/logrus"
)

func IsSignatureValid(message, signature []byte, publicKeys []*rsa.PublicKey) bool {
	hash := crypto.SHA256
	hashed := sha256.Sum256(message)
	log.Debugf("IsSignatureValid() with message sha256 %s and signature:\n%s",
		misc.ShortString(hex.EncodeToString(hashed[:]), 8, 8), misc.ShortString(string(signature), 10, 10))
	sig, _ := base64.StdEncoding.DecodeString(string(signature))
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}

	foundMatchingPublicKey := false
	for _, publicKey := range publicKeys {
		err := rsa.VerifyPSS(publicKey, hash, hashed[:], sig, opts)
		if err == nil {
			foundMatchingPublicKey = true
			log.Debugf(`Public key (E: %d, N: %v) matches.`, publicKey.E, *publicKey.N)
			break
		}
	}
	if !foundMatchingPublicKey { // Only log if none of the keys worked
		log.Warnf(`None of the public keys matched. They were the following:`)
		for _, publicKey := range publicKeys {
			log.Warnf(`Public key (E: %d, N: %v) does not match.`, publicKey.E, fmt.Sprintf("%.10s…", publicKey.N.String()))
		}
	}

	return foundMatchingPublicKey
}
