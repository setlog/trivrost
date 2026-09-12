package main

import (
	"os"
	"testing"
)

// TestReplacePlaceholders checks if the placeholders are replaced and escaped correctly
func TestReplacePlaceholders(t *testing.T) {
	os.Setenv("LAUNCHER_BINARY", "launcher&binary")
	launcherConfig = &config.LauncherConfig{
		VendorName:          "Example<Vendor>",
		BrandingName:        "Brand\"Name",
		BrandingNameShort:   "Short'Name",
		ReverseDnsProductId: "com.example>product",
		ProductVersion: config.ProductVersionStruct{
			Major: 1,
			Minor: 2,
			Patch: 3,
			Build: 4,
		},
	}

	versionSemantic = "1.2.3"
	versionFull = "1.2.3.4"

	// Simulating a template with all placeholders
	template := `${LAUNCHER_BINARY} ${LAUNCHER_BINARY_EXT} ${LAUNCHER_VENDOR_NAME} ${LAUNCHER_BRANDING_NAME} ${LAUNCHER_BRANDING_NAME_SHORT} ${LAUNCHER_REVERSE_DNS_PRODUCT_ID} ${LAUNCHER_VERSION_MAJOR} ${LAUNCHER_VERSION_MINOR} ${LAUNCHER_VERSION_PATCH} ${LAUNCHER_VERSION_BUILD} ${LAUNCHER_VERSION_SEMANTIC} ${LAUNCHER_VERSION_FULL}`

	expected := `launcher&amp;binary ${LAUNCHER_BINARY_EXT} Example&lt;Vendor&gt; Brand&quot;Name Short&apos;Name com.example&gt;product 1 2 3 4 1.2.3 1.2.3.4`
	result := replacePlaceholders(template)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
