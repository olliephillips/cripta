package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func MakeKeys(pubKeyPath, privKeyPath string) error {
	// generate a new RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// encode the private key as a PEM block
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// encode the public key as a PEM block
	publicKey := privateKey.Public()
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	// write the private key to file
	privateKeyFile, err := os.Create(privKeyPath)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	publicKeyFile, err := os.Create(pubKeyPath)
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	if err := pem.Encode(publicKeyFile, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyPEM}); err != nil {
		return err
	}

	return nil
}

// EncryptMessageWithPublicKey does what you might expect
func EncryptMessageWithPublicKey(payload []byte, toUser string, friendKeyPath string) ([]byte, error) {
	/*var encMessage, label []byte

	// get the user pub key
	publicKey, err := getUserPublicKey(toUser, friendKeyPath)
	if err != nil {
		return encMessage, err
	}

	encMessage, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, payload, label)
	if err != nil {
		return encMessage, err
	}

	return encMessage, nil
	*/

	var encMessage, label []byte

	// get the user pub key
	publicKey, err := getUserPublicKey(toUser, friendKeyPath)
	if err != nil {
		return encMessage, err
	}

	hash := sha256.New()

	// chunk message
	msgLen := len(payload)
	step := publicKey.Size() - 2*hash.Size() - 2

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, payload[start:finish], label)
		if err != nil {
			return nil, err
		}

		encMessage = append(encMessage, encryptedBlockBytes...)
	}

	return encMessage, nil
	/*

		msgLen := len(msg)
		    step := public.Size() - 2*hash.Size() - 2
		    var encryptedBytes []byte

		    for start := 0; start < msgLen; start += step {
		        finish := start + step
		        if finish > msgLen {
		            finish = msgLen
		        }

		        encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		        if err != nil {
		            return nil, err
		        }

		        encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
		    }

		    return encryptedBytes, nil
	*/
}

// helper to get the users key
func getUserPublicKey(toUser string, friendKeyPath string) (*rsa.PublicKey, error) {
	var pubk *rsa.PublicKey

	pemData, err := ioutil.ReadFile(fmt.Sprintf("%s.txt", filepath.Join(".", friendKeyPath, strings.ToLower(toUser))))
	if err != nil {
		return pubk, err
	}

	// decode the PEM block to extract the public key
	block, _ := pem.Decode(pemData)
	if block == nil {
		return pubk, errors.New("Failed to decode PEM block")
	}

	// parse the public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return pubk, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

// DecryptMessageWithPrivateKey does what you might expect
func DecryptMessageWithPrivateKey(payload []byte, file string) ([]byte, error) {
	var decMessage, label []byte

	// get the user priv key
	privateKey, err := getPrivateKey(file)
	if err != nil {
		return decMessage, err
	}

	payloadLen := len(payload)
	step := privateKey.PublicKey.Size()
	//fmt.Println("decrypt:", step)

	for start := 0; start < payloadLen; start += step {
		finish := start + step
		if finish > payloadLen {
			finish = payloadLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, payload[start:finish], label)
		if err != nil {
			return nil, err
		}

		decMessage = append(decMessage, decryptedBlockBytes...)
	}

	return decMessage, nil
}

func getPrivateKey(file string) (*rsa.PrivateKey, error) {
	var privk *rsa.PrivateKey

	// read the PEM-encoded private key file
	pemData, err := ioutil.ReadFile(file)
	if err != nil {
		return privk, err
	}

	// decode the PEM block to extract the private key
	block, _ := pem.Decode(pemData)
	if block == nil {
		return privk, errors.New("Failed to decode PEM block")
	}

	// parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return privk, err
	}

	return privateKey, nil
}

/*
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
)

func main() {
	// Generate a new private key.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the message to be signed.
	message := []byte("hello, world")

	// Sign the message with the private key.
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, message)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the public key from the private key.
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	// Verify the signature with the public key.
	if ecdsa.Verify(publicKey, message, r, s) {
		fmt.Println("Signature is valid.")
	} else {
		fmt.Println("Signature is invalid.")
	}
}


*/
