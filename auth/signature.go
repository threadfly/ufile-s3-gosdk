package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

//Authorization
type Credential struct {
	PublicKey   string
	Date        string
	Region      string
	Service     string
	Aws4Request string
}

type AuthMeta struct {
	Credential
	SignedHeaders []string
	Signature     string
	// base64 encode
	EncodeAWSPolicy string
	// base64 encode
	EncodeUCloudPolicy string
}

type UserAccountInfo struct {
	PublicKey  string
	PrivateKey string
}

func CalculationsSignaturePostForm(authMeta *AuthMeta, useraccountinfo *UserAccountInfo) (signature string) {
	//stringToSign
	stringToSign := authMeta.EncodeAWSPolicy

	//Signature
	dateKey := []byte("AWS4" + useraccountinfo.PrivateKey)
	dateKeyHmac := hmac.New(sha256.New, dateKey)
	dateKeyHmac.Write([]byte(authMeta.Date))

	dataRegionKey := []byte(dateKeyHmac.Sum(nil))
	dataRegionKeyHmac := hmac.New(sha256.New, dataRegionKey)
	dataRegionKeyHmac.Write([]byte(authMeta.Region))

	dateRegionServiceKey := []byte(dataRegionKeyHmac.Sum(nil))
	dateRegionServiceKeyHmac := hmac.New(sha256.New, dateRegionServiceKey)
	dateRegionServiceKeyHmac.Write([]byte(authMeta.Service))

	signingKey := []byte(dateRegionServiceKeyHmac.Sum(nil))
	signingKeyHmac := hmac.New(sha256.New, signingKey)
	signingKeyHmac.Write([]byte(authMeta.Aws4Request))

	signatureKey := []byte(signingKeyHmac.Sum(nil))
	signatureKeyHmac := hmac.New(sha256.New, signatureKey)
	signatureKeyHmac.Write([]byte(stringToSign))
	return hex.EncodeToString(signatureKeyHmac.Sum(nil))
}
