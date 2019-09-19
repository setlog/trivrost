package signatures_test

import (
	"crypto/rsa"
	"math/big"
	"testing"

	"github.com/setlog/trivrost/pkg/signatures"
)

const validFile = "Hello World!\n"
const invalidFile = "Bad World!"
const signature = `XvRPGBxSzrtUBVJmQ1EPsWqFo7RpSGVjmbk1e+bKrY9E7OHXbEyGLAu74ik3Ij04
hUyKrSz6tPNsLLI6de5H/HyH1r1NxhZUpJELH90itQvq4SVhdDuz9g2H8iNEm76j
rO7slfTrPbJqz1xjxZ2+b1XEEcrTo6dtlWGo+lIORkcPEUjayem4KCHqwfwy/8ii
UjxWKz73s1EiK/r1N2y4N0BdoAFZkzhRSKgugV4khGEVScszcuCn/KQGWL9xwROM
VjxFoty0xDuDc6lq0htQcoLcv95m+ZsXop75Y1w9FPM0MhCZzI1/kBWRV3DotR4Z
lvCaHvlJXBnRUgQIwASQtQYzjsoEe/8zeQeuE0L8eaabpzAAWEMFPPME9KbeLY7O
PT2BI0C0E0lLnXhSNalACGhbrGpblEV37cpfDeATO/pSU13tT1PaT8QQjHZmIQcl
HwZoLIGsl+1K3esvPMPrcymTS+8Ibj/6/PmTfbIgRlZE2B37oaaSGdmlcDhJaGo5
1G3si5kDi16SVP5jnhEGw4BGtQpZfpDrYX33g69caCuEfKPwUUP9XO4wH+aVRBZQ
LpGEs2ZeqzwsAMFBz1cQIT3wOjDRb/UHsk2ND6txKHoCMypQucDaZbPb+e3hQoTy
/pCDwrBExUCEje3/omcSH94IP/oYuml/3xzI5eVK+SM=`
const publicKeyModulusHexStringMatching = "00a3bf226e9256ff404bcea144d4f5f65873ee495564cdae5d9f07b6dbb764b686dd0f1d0adc60f74b283e46bc80ec6030b86a2a14b956890928d988d7bb8a248e354ce0cc0452e129543925499e563cc24b28bad28f5cf96eb813b9d7d82ed3366c1a8786fbbca4f2b5163c216607dabca44e5f6e6f36fc330ffff77807961d9e18f91e362cad1470c665ce7b0bbc015b31c4394a2f29fda7ebb3dfbd95b9b9e7fe79e231387eaa5f280044954ea9b960f5c65f130a4ba0881ff3fccf3a3505f8845a33816286dcca0a9448f6c7499b0d7ff992a9c14edf7ffba3be9ff74ca5e9ee3ad7487d0039737dc6cf8b2f97b37e7654724602d0f4066c33765131732af48a23441ce6d8300cb245a2cbd086f848cb9b688723640bb4dbc37c4241c1dc35af99a94172b4726a3c71baccc02701b3799c871277220fc51b586ffa526a34b39a4e00e290d4d4a872339da50ef9600f3ff106e25f1a8d69dcb43541daf5134ce76bb0d1f2e71885a49f646f9fac97eb7bbb0ec3c532eba59a86d4b3f83d7811c21e3371cf4789c18ed5a23ebf2a1a0ca54a59be9f39b9c507ea28d7dde5beed993ec736a869e3ac35263dd89c989dbe3b87d06f7d8d4ae8b28bf56becd7f3c789c29c63215529a417722a172b7c10661e50fbe4404089176a86ef15b2263a9a8eb66a176cab7576b5f947070b51e1edc0860ea7d5dc5b508b0bfabfbda5d7f7"
const publicKeyModulusHexStringMismatching = "00c7300be9b9cc7cad94860a0d61a3a66b5aec49286518038d18294b96bd3404ca13db314a05cedcdf64b2edf621fba9e77f59e617e14f9a40fe2291bd1d06e3e595690dcccaa15b2c41f2ff172b5962b724f75bebb5b76c04d3d5b5e58dd912296af0e072c9d14f54c2f6adf9f19d29a9f039235a8af311ac69ab6a5b7c35d8dbd872f20afc9778f740f0929cb38ba8960d4e259a4047f6e279f762501786f7003e085f3420d84509e268ce45e92ebbf164a284940ffa5d6d9822a35b1a1d3370850bf9c45f5cbb2d44366edb5777bfd75599884ee355dcb6e4d2a51de46e68bb474a9f7a4039e001c7cae3d6df5c684149fbb010ba85130a8660a87afae44017"
const publicKeyExponent = 65537

func getPublicKey(publicKeyModulusHexString string) *rsa.PublicKey {
	modulus, _ := big.NewInt(0).SetString(publicKeyModulusHexString, 16)
	return &rsa.PublicKey{N: modulus, E: publicKeyExponent}
}

func TestValidSignatureWithOneKey(t *testing.T) {
	publicKeys := []*rsa.PublicKey{getPublicKey(publicKeyModulusHexStringMatching)}
	valid := signatures.IsSignatureValid([]byte(validFile), []byte(signature), publicKeys)
	if !valid {
		t.Error("Valid signature not recognized to be valid.")
	}
}

func TestValidSignatureWithFirstKeyMatching(t *testing.T) {
	publicKeys := []*rsa.PublicKey{getPublicKey(publicKeyModulusHexStringMatching), getPublicKey(publicKeyModulusHexStringMismatching)}
	valid := signatures.IsSignatureValid([]byte(validFile), []byte(signature), publicKeys)
	if !valid {
		t.Error("Valid signature not recognized to be valid.")
	}
}

func TestValidSignatureWithLastKeyMatching(t *testing.T) {
	publicKeys := []*rsa.PublicKey{getPublicKey(publicKeyModulusHexStringMismatching), getPublicKey(publicKeyModulusHexStringMatching)}
	valid := signatures.IsSignatureValid([]byte(validFile), []byte(signature), publicKeys)
	if !valid {
		t.Error("Valid signature not recognized to be valid.")
	}
}

func TestBadSignatureWithOneKey(t *testing.T) {
	publicKeys := []*rsa.PublicKey{getPublicKey(publicKeyModulusHexStringMatching)}
	valid := signatures.IsSignatureValid([]byte(invalidFile), []byte(signature), publicKeys)
	if valid {
		t.Error("Invalid signature not recognized to be invalid.")
	}
}

func TestBadSignatureWithTwoKeys(t *testing.T) {
	publicKeys := []*rsa.PublicKey{getPublicKey(publicKeyModulusHexStringMatching), getPublicKey(publicKeyModulusHexStringMismatching)}
	valid := signatures.IsSignatureValid([]byte(invalidFile), []byte(signature), publicKeys)
	if valid {
		t.Error("Invalid signature not recognized to be invalid.")
	}
}
