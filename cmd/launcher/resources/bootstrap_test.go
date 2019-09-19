package resources_test

import (
	"crypto/rsa"
	"math/big"
	"testing"

	"github.com/setlog/trivrost/cmd/launcher/resources"
)

const publicRsaKey1 string = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAo78ibpJW/0BLzqFE1PX2
WHPuSVVkza5dnwe227dktobdDx0K3GD3Syg+RryA7GAwuGoqFLlWiQko2YjXu4ok
jjVM4MwEUuEpVDklSZ5WPMJLKLrSj1z5brgTudfYLtM2bBqHhvu8pPK1FjwhZgfa
vKROX25vNvwzD//3eAeWHZ4Y+R42LK0UcMZlznsLvAFbMcQ5Si8p/afrs9+9lbm5
5/554jE4fqpfKABElU6puWD1xl8TCkugiB/z/M86NQX4hFozgWKG3MoKlEj2x0mb
DX/5kqnBTt9/+6O+n/dMpenuOtdIfQA5c33Gz4svl7N+dlRyRgLQ9AZsM3ZRMXMq
9IojRBzm2DAMskWiy9CG+EjLm2iHI2QLtNvDfEJBwdw1r5mpQXK0cmo8cbrMwCcB
s3mchxJ3Ig/FG1hv+lJqNLOaTgDikNTUqHIznaUO+WAPP/EG4l8ajWnctDVB2vUT
TOdrsNHy5xiFpJ9kb5+sl+t7uw7DxTLrpZqG1LP4PXgRwh4zcc9HicGO1aI+vyoa
DKVKWb6fObnFB+oo193lvu2ZPsc2qGnjrDUmPdicmJ2+O4fQb32NSuiyi/Vr7Nfz
x4nCnGMhVSmkF3IqFyt8EGYeUPvkQECJF2qG7xWyJjqajrZqF2yrdXa1+UcHC1Hh
7cCGDqfV3FtQiwv6v72l1/cCAwEAAQ==
-----END PUBLIC KEY-----`
const exponentOfPublicRsaKey1 int = 65537
const modulusHexStringOfPublicRsaKey1 string = "00a3bf226e9256ff404bcea144d4f5f65873ee495564cdae5d9f07b6dbb764b686dd0f1d0adc60f74b283e46bc80ec6030b86a2a14b956890928d988d7bb8a248e354ce0cc0452e129543925499e563cc24b28bad28f5cf96eb813b9d7d82ed3366c1a8786fbbca4f2b5163c216607dabca44e5f6e6f36fc330ffff77807961d9e18f91e362cad1470c665ce7b0bbc015b31c4394a2f29fda7ebb3dfbd95b9b9e7fe79e231387eaa5f280044954ea9b960f5c65f130a4ba0881ff3fccf3a3505f8845a33816286dcca0a9448f6c7499b0d7ff992a9c14edf7ffba3be9ff74ca5e9ee3ad7487d0039737dc6cf8b2f97b37e7654724602d0f4066c33765131732af48a23441ce6d8300cb245a2cbd086f848cb9b688723640bb4dbc37c4241c1dc35af99a94172b4726a3c71baccc02701b3799c871277220fc51b586ffa526a34b39a4e00e290d4d4a872339da50ef9600f3ff106e25f1a8d69dcb43541daf5134ce76bb0d1f2e71885a49f646f9fac97eb7bbb0ec3c532eba59a86d4b3f83d7811c21e3371cf4789c18ed5a23ebf2a1a0ca54a59be9f39b9c507ea28d7dde5beed993ec736a869e3ac35263dd89c989dbe3b87d06f7d8d4ae8b28bf56becd7f3c789c29c63215529a417722a172b7c10661e50fbe4404089176a86ef15b2263a9a8eb66a176cab7576b5f947070b51e1edc0860ea7d5dc5b508b0bfabfbda5d7f7"

const publicRsaKey2 string = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxzAL6bnMfK2UhgoNYaOm
a1rsSShlGAONGClLlr00BMoT2zFKBc7c32Sy7fYh+6nnf1nmF+FPmkD+IpG9HQbj
5ZVpDczKoVssQfL/FytZYrck91vrtbdsBNPVteWN2RIpavDgcsnRT1TC9q358Z0p
qfA5I1qK8xGsaatqW3w12NvYcvIK/Jd490Dwkpyzi6iWDU4lmkBH9uJ592JQF4b3
AD4IXzQg2EUJ4mjORekuu/FkooSUD/pdbZgio1saHTNwhQv5xF9cuy1ENm7bV3e/
11WZiE7jVdy25NKlHeRuaLtHSp96QDngAcfK49bfXGhBSfuwELqFEwqGYKh6+uRA
FwIDAQAB
-----END PUBLIC KEY-----`
const exponentOfPublicRsaKey2 int = 65537
const modulusHexStringOfPublicRsaKey2 string = "00c7300be9b9cc7cad94860a0d61a3a66b5aec49286518038d18294b96bd3404ca13db314a05cedcdf64b2edf621fba9e77f59e617e14f9a40fe2291bd1d06e3e595690dcccaa15b2c41f2ff172b5962b724f75bebb5b76c04d3d5b5e58dd912296af0e072c9d14f54c2f6adf9f19d29a9f039235a8af311ac69ab6a5b7c35d8dbd872f20afc9778f740f0929cb38ba8960d4e259a4047f6e279f762501786f7003e085f3420d84509e268ce45e92ebbf164a284940ffa5d6d9822a35b1a1d3370850bf9c45f5cbb2d44366edb5777bfd75599884ee355dcb6e4d2a51de46e68bb474a9f7a4039e001c7cae3d6df5c684149fbb010ba85130a8660a87afae44017"

const publicDsaKey string = `-----BEGIN PUBLIC KEY-----
MIIDRjCCAjkGByqGSM44BAEwggIsAoIBAQDO3PXgDK5txaU/eVRKgdZKQOJNryGk
qKJmHWukkg1WcRdIbCl2tGnCRTzvGi55XjNyzB8s11p9Z+NpBPOgOuIwbBakG7IB
Iy0skuYpK8s0j8RU+utXRSIaCTpEeyiE8+NTql7HhvqYN3GdsJ+JC2on1+vrOwPO
wmAGnsmnaeY3rB8fkB4UrbAuPU3/SZvnpNjx8QHkooGIQvyR/uL4+UCT2o4cf9HX
krQd1kTwqemz+DastW/88cCMjEOC7Hj5zR2xmymDTOGyIlaEq08hqmfe1OqHZqZf
w/3G0DRO3XDbo6otrFvZIFvlV7WkqUqahyOwqYCkeLTDCd87vALMD81LAiEA5uXt
bd+TtIr+uuFyPXyCCp81StE0BC06kt+Saxch7AcCggEAKR4uo+lOZ81ywZs5jnCI
a0OHMZKiljEuS9UlWs/LvexImOSBTD7RKzaX2psUfnlZ0mzDjrDBeyyZesKUNnVu
+K/EM8Sgy41FSTaUVc9U3PL1eGVSJ7OVrIDOWiLMWo3qIvmw4xrsLTi/9jrNrleJ
aOpTHkPXieDHp0Qz6kg7w9E7mu9cgUvkcp75fRlzAydJGd8iRQHXZUnzcsr8juiK
RCxWHpWFfc7SAbrQrBW6FTlr92tuN8VQq+xybFcPYXy3g6nHYnRnOHpm1AkTGnbD
pnIjOGuua2s5mG6rL/n3ySpPkbiMYrP9SOn6lHoXiGUrJTyf0uwHASVeh4pluV/A
7AOCAQUAAoIBAEoHiWw08f2fAXMNQcluCQi5z9J0KIvIdlzpnqyOFtr90CCVY26J
gyFbwkFTBI0Vqx50z6CsAAnlKNcaUmMriFlbsYDmY+Jr1Xps+t5c05KG0uvvUYWS
R1ME/XCACl96qRwmWvy5MzCizGRepEDNgyKyKBABCG66gN9nKGeRY30nzOqfKKjw
FMt3FxAuWej3aLlKIGmM+jHHj4y0Df5zz+CVK97nB9LV4XFiCuzQ965cKK3YOEmz
9lFTRzaxX3i/p5A1c5+nO59dy8MqXHZBvNxy4cIvq6XQFtCGCP26elqt5NN+2UvY
mRZmiBzyYuHLtZmcY1smzfpSp/V7nrdh4jE=
-----END PUBLIC KEY-----`

const publicNonsenseKey string = `Hallo Welt`

func TestToParseSingleRsaKey(t *testing.T) {
	publicKeys := resources.ReadPublicRsaKeysAsset(publicRsaKey1)
	if len(publicKeys) != 1 {
		t.Errorf("Read %d public Keys, but expected only 1.", len(publicKeys))
	}
	checkRsaKey(t, exponentOfPublicRsaKey1, modulusHexStringOfPublicRsaKey1, publicKeys[0])
}

func TestToParseTwoRsaKeys(t *testing.T) {
	publicKeys := resources.ReadPublicRsaKeysAsset(publicRsaKey1 + "\n" + publicRsaKey2)
	if len(publicKeys) != 2 {
		t.Errorf("Read %d public Keys, but expected exactly 2.", len(publicKeys))
	}
	checkRsaKey(t, exponentOfPublicRsaKey1, modulusHexStringOfPublicRsaKey1, publicKeys[0])
	checkRsaKey(t, exponentOfPublicRsaKey2, modulusHexStringOfPublicRsaKey2, publicKeys[1])
}

func TestToParseNotSupportedDsaKey(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The parsing did not panic with a non-RSA key.")
		}
	}()

	resources.ReadPublicRsaKeysAsset(publicDsaKey)
}

func TestToParseInvalidKeysString(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The parsing did not panic with an invalid key.")
		}
	}()

	resources.ReadPublicRsaKeysAsset(publicNonsenseKey)
}

func checkRsaKey(t *testing.T, expectedExponent int, expectedModulusHexString string, actualRsaKey *rsa.PublicKey) {
	if actualRsaKey.E != expectedExponent {
		t.Errorf("Didn't parse the correct exponent of the key.")
	}
	expectedModulus, _ := big.NewInt(0).SetString(expectedModulusHexString, 16)
	if expectedModulus.Cmp(actualRsaKey.N) != 0 {
		t.Errorf("Didn't parse the correct modulus of the key.")
	}
}
