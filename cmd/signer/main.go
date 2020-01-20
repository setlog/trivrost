package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	privateKeyFile string
	targetFiles    []string
)

func main() {
	parseFlags()
	checkIfFilesExist()
	createSignatures()
}

func createSignatures() {
	pemString := readFile(privateKeyFile)
	key := readPrivateKey(pemString)

	for i := range targetFiles {
		fileContent := readFile(targetFiles[i])
		signature, err := createFileSignature(key, fileContent)
		if err != nil {
			fatalf("Creating of a signature for the file %s failed: %v", targetFiles[i], err)
		}
		if err = ioutil.WriteFile(targetFiles[i]+".signature", []byte(signature), 0644); err != nil {
			fatalf("Could not write a signature into the file %s.signature: %v", targetFiles[i], err)
		}
	}

}

func createFileSignature(key *rsa.PrivateKey, fileContent []byte) (string, error) {
	hashed := sha256.Sum256(fileContent)
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}

	signed, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, hashed[:], opts)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signed), nil
}

func readFile(fileName string) []byte {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fatalf("Could not read a file %s: %v", fileName, err)
	}
	return content
}

func readPrivateKey(pemString []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(pemString)
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		parseResult, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	}
	if err != nil {
		fatalf("Could not parse the private key as neither PKCS8- nor PKCS1-formed: %v.", err)
	}
	return parseResult.(*rsa.PrivateKey)
}

func parseFlags() {
	flag.Parse()
	if flag.NArg() < 2 {
		fatalf("Need at least 2 args: privateKeyFile targetFile1 targetFile2 ...")
	}

	privateKeyFile = flag.Arg(0)
	targetFiles = make([]string, len(flag.Args())-1)
	for i := 1; i < len(flag.Args()); i++ {
		targetFiles[i-1] = flag.Arg(i)
	}
}

func checkIfFilesExist() {
	if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
		fatalf("privateKeyFile %s does not exist", privateKeyFile)
	}
	for i := range targetFiles {
		if _, err := os.Stat(targetFiles[i]); os.IsNotExist(err) {
			fatalf("targetFile %s does not exist", targetFiles[i])
		}
	}
}

func fatalf(formatMessage string, args ...interface{}) {
	fmt.Printf("Fatal: "+formatMessage+"\n", args...)
	os.Exit(1)
}
