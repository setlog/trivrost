package main

import (
	"crypto/rsa"
	"log"
	"math/big"
	"testing"

	"github.com/setlog/trivrost/cmd/launcher/resources"
	"github.com/setlog/trivrost/pkg/signatures"
)

var testContent = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."
var testSignature = `j6IM9PlV4VABelDIzRfZ5yd8xHTxeRXRJ8SCym+vW8/JY1zAzXEHB8M0oC+3Gbk/i32kOIpR0jDWhU0LZqxysQ7DtWGVAl463aD4skr/nrOK7rPb4uUfFr0FEQ1d6PZd6D4aecPprCTllD6TqDMZaMXf1Ng6D4LCi1N5lO28nfPTtGsnmk4TSzoQEXu6VYrZsxlbY9nSr/m2ko4A8BRtDSckaKjk3jRcCHPzvPFszuCYwiYXheRYRp5oyFrqgwt6V0sNVjaAbzTyeRQxnO1rJllQh2TrLHTxmuAzAIPWvj8wgWmW7wpvrr+23CTMhsw8VTvDnaI8Jk+GnM0hWuvtcw==`
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

const publicKeyExponent = 65537

func TestSignatureFromFile(t *testing.T) {
	key := readPrivateKey([]byte(testPrivateKey))
	s, err := createFileSignature(key, []byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}

	signature := []byte(s)

	if len(testSignature) != len(signature) {
		log.Println(string(signature))
		t.Errorf("Different length of signatures detected. Generated signature: %d bytes, expected: %d bytes", len(signature), len(testSignature))
		return
	}

	for i := range testSignature {
		if signature[i] != testSignature[i] {
			t.Errorf("Checking of signature failed at %d byte", i)
			return
		}
	}
}

func TestSignatureWithFunction(t *testing.T) {
	publicKeys := resources.ReadPublicRsaKeysAsset(testPublicKey)
	key := readPrivateKey([]byte(testPrivateKey))
	s, err := createFileSignature(key, []byte(testContent))
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(s)

	signature := []byte(s)

	valid := signatures.IsSignatureValid([]byte(testContent), []byte(signature), publicKeys)
	if !valid {
		t.Error("Valid signature not recognized to be valid.")
	}
}

func getPublicKey(publicKeyModulusHexString string) *rsa.PublicKey {
	modulus, _ := big.NewInt(0).SetString(publicKeyModulusHexString, 16)
	return &rsa.PublicKey{N: modulus, E: publicKeyExponent}
}
