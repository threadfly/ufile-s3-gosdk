package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	auth "github.com/threadfly/ufile-s3-gosdk/auth"
)

func main() {
	timeLayout := "20180826T000000Z"
	var u auth.UserAccountInfo
	u.PrivateKey = ""
	u.PublicKey = "27qXHOY3xPjKlpNjL0Aiou0hQ+LRnJwooDjSIvqeswP1MdQQwKVvBg=="

	extraParams := map[string]string{
		//	"policy":           "ewoiZXhwaXJhdGlvbiI6ICIyMDE4MDgyNlQwMDAwMDBaIiwKImNvbmRpdGlvbnMiOiBbCiAgICAgeyJidWNrZXQiOiAidGVzdC13ZmYtYmoifSwKICAgICB7IlVDbG91ZFBvbGljeSI6ICJld29KSW1OaGJHeGlZV05yVlhKc0lqb2lhSFIwY0Rvdkx6RXdOaTQzTlM0eU16RXVNakk1TDJOaGJHeGlZV05yWkdWdGJ5OWtaVzF2TDJobGJHeHZJaXdLQ1NKallXeHNZbUZqYTBKdlpIa2lPaUppZFdOclpYUTlKQ2hpZFdOclpYUXBKbXRsZVQwa0tHdGxlU2ttYm1GdFpUMHlNV0ZuWlNJS2ZRPT0ifQogXQp9",
		"x-amz-algorithm":  "AWS4-HMAC-SHA256",
		"x-amz-credential": fmt.Sprintf("%s/%s/%s/%s/aws4_request", u.PublicKey, timeLayout, "cn-bj", "s3"),
		"x-amz-date":       timeLayout,
		"x-amz-signature":  "",
		//	"uploadPolicy":     "ewoJImNhbGxiYWNrVXJsIjoiaHR0cDovLzEwNi43NS4yMzEuMjI5L2NhbGxiYWNrZGVtby9kZW1vL2hlbGxvIiwKCSJjYWxsYmFja0JvZHkiOiJidWNrZXQ9JChidWNrZXQpJmtleT0kKGtleSkmbmFtZT0yMWFnZSIKfQ==",
		"key": "20180728-xxxxxttt-0.go",
	}

	var authMeta auth.AuthMeta
	authMeta.Credential.PublicKey = u.PublicKey
	authMeta.Credential.Date = timeLayout
	authMeta.Credential.Region = "cn-bj"
	authMeta.Credential.Service = "s3"
	authMeta.Credential.Aws4Request = "aws4_request"
	authMeta.EncodeAWSPolicy = extraParams["policy"]

	sig := auth.CalculationsSignaturePostForm(&authMeta, &u)
	log.Printf("signature: %s \n ", sig)
	extraParams["x-amz-signature"] = sig

	bucket := "test-wff-bj"

	request, err := newfileUploadRequest("http://cn-bj-s3.ufileos.com/"+bucket, extraParams, "file", "./postObjcet.go")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
