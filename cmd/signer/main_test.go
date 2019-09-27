package main

import (
	"crypto/rsa"
	"math/big"
	"testing"

	"github.com/setlog/trivrost/cmd/launcher/resources"
	"github.com/setlog/trivrost/pkg/signatures"
)

// text that will be signed and is going to be used in tests below
var testContent = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."

// tests that use this private key for signing have to be OK
var testPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDoPNEZXNtQSULq
JSB/MlscSA13J+JefnXfdTNqHrvqKDKiJbcYwIAXy2qTA9jib3+wbjlOHHZCbdUo
uKuf1jaHvzcTmhk5YSE5IweiV6r/EmJGoog4bVEG5auR8nAkvkam3KAoRfQJN9xK
xgjcQTO7IUVaQdekVuVPLaRdfa0e5SN1IipEBIboKr9mz3taACiMWXbyUUE7WpEl
ctrLnnGTRmVGQPwLMeDw1ME+klpMma3I4ILdFg/UbbPoplC8NcluNFg/7+2uS3pU
SZvfLT352LcbCfYdXaES2kZUnMv+PUCjEprySzUKK9YD0A5xU/6VAR8GTQSJBxq2
XeyPwBEbAgMBAAECggEBAIatrj1dIjpfIhUTTtM06q1uA5EUaiyOfeEG4Lgr9qIG
icaKxLHwANjLuJRlaMN4Eb7JTSZFTzea5kDlR3I8EgeLFm+hr/scnt25uNWmrZ2a
la+M1h6TFqg/TM4ooGxOhD6EN8TjPHCUGoaqbbz9eviMhOGgyWOemQDf4S/ukBUX
hiTrLPwlJkwBtJagzfCfm9+1qBRy3ZgvkIxheSZldHsPiDsMgvjTsfDrd2FkjJNi
Hr79x+yW7dZ3XxZYHASh2io10IdpYgsV+yIS6CpRRWxsWp3k6tQBC39EvMZ3Wk90
tFRyJaKY467Dq8e6dIX1TxVJZsWt+NUmW5FOTYPozUECgYEA+YbUWc2Hblrzf6ae
15LQU7wZK7d+YpDD7upDVl9qw4x+5WcTzrS/87iLZH3SLat6p1+tFe+EtJemyYNA
CqdsiR+46OTxKVCDMEQDrDg5Vg2IpVpWiqD5b0ag1wkE6SEq5DC+hbX3WKAtYRxz
XPhoynKAdxFivOuJnWMgd/I4XDECgYEA7kMqdzEAX41q+c3eIth8w5aD5UNKe8O7
ZISzIzZ5JY5e1IzR4VXqRr9pCCw2uWhY7MfSUPxJVaaZ9sq2Bd7qDI/Ypza3hQiM
NGajmi4rqp24oOnOQ4laFKuN4zyjlnX+dw8DTpQ9S1PkOU9EQRAZA25bjhWC5UPo
6CubjH++CwsCgYEA7a6Ks4fdCzdDXkJ+Z2WHX1t6tnOwxX6TxA4NWkbFUcOQVD/d
VDZD6YnN7UkUXUBMMwYlvxFJ3SPfUW/eHsff0LYQ0nbRaMMyU1VWEkP0CY4WrTrh
2GcBcgdaybnjnZVkX7w2nvL3ysm4sBoDoXlViBGNYN2EqePKT8rOcLKfEOECgYEA
obg/I6XD7hdr6+B7DUXJ8WvBXKS+8qCZGhIkERuRQReQcE6gyoTpPln/bYetIU2d
RiIfM87568PoLyXKRNPYIuykDmNKT2bM22hrVWRPSUBCqB3qXdblqLAE358yHhc6
wA8VnIlrzSxE9U1DM7I8eCK4zAj3zqu4c5Xdv5CZKp8CgYEA2hfqWbev5HFta2kk
vHWdAylUvtBNDPJ7yFLH5XXARhTRkctr+rslr/KuNo5G/nTR6ozH+XpCueYyEeZF
p60bRu7q0C8r5Ameks9AGOhcoYGajADDW7P09jtC5OcC+Kyv18UlkWwzvlw9JR1C
YaHChpQ7+mdUR0V+Yn4Y0vCSuTs=
-----END PRIVATE KEY-----
`

// pair of testPrivateKey: tests that use this public key for verifying have to be OK
var testPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6DzRGVzbUElC6iUgfzJb
HEgNdyfiXn5133Uzah676igyoiW3GMCAF8tqkwPY4m9/sG45Thx2Qm3VKLirn9Y2
h783E5oZOWEhOSMHoleq/xJiRqKIOG1RBuWrkfJwJL5GptygKEX0CTfcSsYI3EEz
uyFFWkHXpFblTy2kXX2tHuUjdSIqRASG6Cq/Zs97WgAojFl28lFBO1qRJXLay55x
k0ZlRkD8CzHg8NTBPpJaTJmtyOCC3RYP1G2z6KZQvDXJbjRYP+/trkt6VEmb3y09
+di3Gwn2HV2hEtpGVJzL/j1AoxKa8ks1CivWA9AOcVP+lQEfBk0EiQcatl3sj8AR
GwIDAQAB
-----END PUBLIC KEY-----
`

// tests that use this private key for signing are supposed to fail
var testPrivateKey2Fail = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCm1IYb2ZydyNJc
R73uRvESjyqOlEtSQ3AYvM8dyaykW6kXUohWBNfh/EYE0nhhyPROtwmbaGjfGTEt
dxtYlfRHeLWod/dCpXxVISrZ9k3MplG22idPZvgxePSrdqWVk/2llOlm5AUYoU7B
voMveTuGgknyyssD6cu/mP8f4icUzNVTgYVxMug9lsr9x+IfArzvXSsnP5eY1V6c
pdbAtsb1RTNAbMZsqQzM453DY8rhqarSLtXp2Zxq6vmEvJN+BEF8sHF8yUiAQq6T
eGpyEb+oxy4wT+iPrlX0+o6Se0cfUvewa7+Ii8MV1f2I2kDvGVPvBJMhPi86FhmI
M2pxLVOFAgMBAAECggEAbI58GaE3lUB5Cc0xHNySv8XjJlX+0S/KwH4Ts8loiqaO
V/u/dWG/bHCwyzB9XvvZZWMbYEHHg+yroG8Rn0osY1l7s30kqvxt9CMZ9CyeoV1U
bMx1qehR9jdD1lLlGnjrIxTL78TOQCGu0sl6KakUf8lF/zPQeOJoT2tqD8AkOBao
XmOlnsEmCW/PJGM/e/oBsuxNq2uUAVOEiOxGhZqIN7HWfhTinKHX28FJ0pLCkuQ/
D7Si5icXEfrD1pD81Ft9v35sSQrU3YZk2jR0z/mDBfi/MPp2orTA5ZfHes+hpnEc
7yLS4wHNa+v760l3L7q1UGPKVycZv484zKmYFIQGHQKBgQDVZg9drlqBIgHXFE4x
wYXBZfD8t2ll/3lCN6xdBi/u1HP0ZGHICZ2l9lhHwXeeCEblY2Al2xjSdtnLfK5c
CnSvwu9UV8h4sB4kjn37PL0SohtZJlFZpt6vresoDpQQrMmnjhgr+ot9ODlH1bLF
yDnCtVXGlGYwPYAn40IxrmB5lwKBgQDIIolqxU7WkZ7nARhR++G3NeVqV8sB7GWX
LyPrnoKbb556BSM2Ourcu/3uQfWnMwrc5J53tVHZ6LjrB4pSPa546LooWv5vY9QD
ETwQzUKNfTiGT4qwkUdBs08fU9qmy0TT9y2SKbpe76GAS4oPbpYbk7W78jJRiWAZ
tZelUcSnQwKBgHJ/BwGRmdetQmV+7JF/rt9cbdd6JR/n2cywiFeFCVTQQsK+1UP5
/M7eBPHDGQX+lONg1WaaTpAl2qd2ZyrVJVRkd/q9+r7eZ93fYjLZnOyRc7D6gS1j
/hkubHyajdEAlFXFRKzcCdmOwBUN0JST4IHav4IDf2ykos1D/vEfCX5TAoGBAMYd
ykKzx3OI+/BZmSWvXqXq6Iv5FLF2vqqGs9xPMaOFPzAzXcQVVuHkB1+QVAmL8bjx
aB3AlKJOSp/++uKmxMxUNdQ1H6JNBFd0/Cz1xGgkCYyLuRNI/W0Af9bXP5/VoPDj
w2zpeeD4/rruDGFya44pDsJa44zrnQJWTSQOacnZAoGAM9XjNyoue1g7mrMCFaal
ff7KnzqzMn9QuzIP/C9ImuUB/z/FG7Rn98oK66RJuPu4/oix7AU8UDI5JrD7Evwr
WdEI/ZhvdTMMrho9fruTFINBU0LaGEBN/xt4cv9crGcEvenoPE0WRoxhZBb2disg
agNtVIy2S07B+Eixp/gqPW8=
-----END PRIVATE KEY-----
`

// tests that use this public key for verifying will fail
var testPublicKey2Fail = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAptSGG9mcncjSXEe97kbx
Eo8qjpRLUkNwGLzPHcmspFupF1KIVgTX4fxGBNJ4Ycj0TrcJm2ho3xkxLXcbWJX0
R3i1qHf3QqV8VSEq2fZNzKZRttonT2b4MXj0q3allZP9pZTpZuQFGKFOwb6DL3k7
hoJJ8srLA+nLv5j/H+InFMzVU4GFcTLoPZbK/cfiHwK8710rJz+XmNVenKXWwLbG
9UUzQGzGbKkMzOOdw2PK4amq0i7V6dmcaur5hLyTfgRBfLBxfMlIgEKuk3hqchG/
qMcuME/oj65V9PqOkntHH1L3sGu/iIvDFdX9iNpA7xlT7wSTIT4vOhYZiDNqcS1T
hQIDAQAB
-----END PUBLIC KEY-----
`

const publicKeyExponent = 65537

func TestSignatureOK(t *testing.T) {
	publicKeys := resources.ReadPublicRsaKeysAsset(testPublicKey)
	key := readPrivateKey([]byte(testPrivateKey))
	signature, err := createFileSignature(key, []byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}

	valid := signatures.IsSignatureValid([]byte(testContent), []byte(signature), publicKeys)
	if !valid {
		t.Error("Valid signature not recognized to be valid.")
	}
}

func TestSignatureFail(t *testing.T) {
	publicKeys := resources.ReadPublicRsaKeysAsset(testPublicKey)
	key := readPrivateKey([]byte(testPrivateKey))
	signature, err := createFileSignature(key, []byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}

	out := []rune(signature)
	out[3] = out[3] + 1
	signature = string(out)

	valid := signatures.IsSignatureValid([]byte(testContent), []byte(signature), publicKeys)
	if valid {
		t.Error("Signature recognized as valid but it should not.")
	}
}

func TestSignatureWithWrongPrivateKey(t *testing.T) {
	publicKeys := resources.ReadPublicRsaKeysAsset(testPublicKey)
	key := readPrivateKey([]byte(testPrivateKey2Fail))
	s, err := createFileSignature(key, []byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}

	signature := []byte(s)

	valid := signatures.IsSignatureValid([]byte(testContent), []byte(signature), publicKeys)
	if valid {
		t.Error("Signature recognized as valid but it was created with wrong private key.")
	}
}

func TestSignatureWithWrongPublicKey(t *testing.T) {
	publicKeys := resources.ReadPublicRsaKeysAsset(testPublicKey2Fail)
	key := readPrivateKey([]byte(testPrivateKey))
	s, err := createFileSignature(key, []byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}

	signature := []byte(s)

	valid := signatures.IsSignatureValid([]byte(testContent), []byte(signature), publicKeys)
	if valid {
		t.Error("Signature recognized as valid but it was verified with wrong public key.")
	}
}

func getPublicKey(publicKeyModulusHexString string) *rsa.PublicKey {
	modulus, _ := big.NewInt(0).SetString(publicKeyModulusHexString, 16)
	return &rsa.PublicKey{N: modulus, E: publicKeyExponent}
}
