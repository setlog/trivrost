package resources

import (
	"crypto/rsa"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

// Provisioned by launcher-config.json.gen.go
var LauncherConfig *config.LauncherConfig

// Provisioned by public-rsa-keys.pem.gen.go
var PublicRsaKeys []*rsa.PublicKey

// Provisioned by icon.png.gen.go
var LauncherIcon []byte
