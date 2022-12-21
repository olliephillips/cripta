package crypt

import (
	"crypto"
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

func MakeKeys(pubKeyPath, privKeyPath, friendSelf string) error {
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

	// duplicate pub key in friends folder
	friendSelfKey, err := os.Create(friendSelf)
	if err != nil {
		return err
	}
	defer friendSelfKey.Close()

	if err := pem.Encode(friendSelfKey, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyPEM}); err != nil {
		return err
	}

	return nil
}

// EncryptMessageWithPublicKey does what you might expect
func EncryptMessageWithPublicKey(payload []byte, toUser string, friendKeyPath string) ([]byte, error) {

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

	// sign the message

	return encMessage, nil
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

// SignPayload returns a signature of the payload using sender private key
func SignPayload(payload []byte, file string) ([]byte, error) {
	var sig []byte

	// get the user priv key
	privateKey, err := getPrivateKey(file)
	if err != nil {
		return sig, err
	}

	hashed := sha256.Sum256(payload)

	sig, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return sig, err
	}

	return sig, nil
}

// VerifySignature checks that the signture is valid for the sending user
func VerifySignature(sender string, msg []byte, friendKeyPath string, sig []byte) (bool, error) {
	// get the user pub key
	publicKey, err := getUserPublicKey(sender, friendKeyPath)
	if err != nil {
		return false, err
	}

	hashed := sha256.Sum256(msg)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], sig); err != nil {
		return false, err
	}

	return true, nil
}
