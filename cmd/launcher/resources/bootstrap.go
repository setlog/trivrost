package resources

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"strings"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

const (
	defaultBrandingName = "Launcher"
	defaultVendorName   = "trivrost_vendor"
	defaultProductName  = "trivrost_launcher"
)

func readLauncherConfigAsset(a string) *config.LauncherConfig {
	launcherConfig := config.ReadLauncherConfigFromReader(strings.NewReader(a))
	if launcherConfig.BrandingName == "" {
		log.Printf("BrandingName empty. Setting it to \"%s\".", defaultBrandingName)
		launcherConfig.BrandingName = defaultBrandingName
	}
	if launcherConfig.VendorName == "" {
		log.Printf("VendorName empty. Setting it to \"%s\".", defaultVendorName)
		launcherConfig.VendorName = defaultVendorName
	}
	if launcherConfig.ProductName == "" {
		log.Printf("ProductName empty. Setting it to \"%s\".", defaultProductName)
		launcherConfig.ProductName = defaultProductName
	}
	return launcherConfig
}

func ReadPublicRsaKeysAsset(a string) []*rsa.PublicKey {
	publicKeys := []*rsa.PublicKey{}
	var block *pem.Block
	rest := []byte(a)
	for len(rest) > 0 {
		block, rest = pem.Decode(rest)
		if block == nil || block.Type != "PUBLIC KEY" {
			log.Panic("Invalid pem block or pem block without a public key.")
		}

		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			log.Panic(err)
		}

		pubRsa, isPublicKey := pub.(*rsa.PublicKey)
		if !isPublicKey {
			log.Panic("No RSA key.")
		}

		publicKeys = append(publicKeys, pubRsa)
	}
	return publicKeys
}

func readIconAsset(a string) []byte {
	return []byte(a)
}
