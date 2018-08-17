package s3test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"

	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"time"

	. "../Utilities"
)

func (suite *S3Suite) TestPresignRequest() {

	assert := suite
	region := viper.GetString("s3main.region")
	req, body := SetupRequest("S3", region, "{}")

	signer := SetupSigner(Creds)
	signer.Presign(req, body, "s3", region, 300*time.Second, time.Unix(0, 0))
	qry := req.URL.Query()
	var credentials string = viper.GetString("s3main.access_key") + "/" + "19700101" + "/" 
	credentials = credentials + viper.GetString("s3main.region") + "/" + "s3" + "/" + "aws4_request"
	assert.Equal(credentials, qry.Get("X-Amz-Credential"))
	assert.Equal("content-length;content-type;host;x-amz-meta-other-header;x-amz-meta-other-header_with_underscore", qry.Get("X-Amz-SignedHeaders"))
	assert.Equal("19700101T000000Z", qry.Get("X-Amz-Date"))
}

func (suite *S3Suite) TestSignRequest() {

	assert := suite
	region := viper.GetString("s3main.region")
	req, body := SetupRequest("S3", region, "{}")
	var credentials string = viper.GetString("s3main.access_key") + "/" + "19700101" + "/" 
	credentials = credentials + viper.GetString("s3main.region") + "/" + "s3" + "/" + "aws4_request"
	expectedauth := "AWS4-HMAC-SHA256 Credential=" + credentials + ", SignedHeaders=content-length;content-type;host;x-amz-content-sha256;x-amz-date;x-amz-meta-other-header;x-amz-meta-other-header_with_underscore;x-amz-target"
	signer := SetupSigner(Creds)

	signer.Sign(req, body, "s3", region, time.Unix(0, 0))

	qry := req.Header
	assert.Contains(qry.Get("Authorization"), expectedauth)
	assert.Equal("19700101T000000Z", qry.Get("X-Amz-Date"))
}

func (suite *S3Suite) TestSignBody() {

	assert := suite
	region := viper.GetString("s3main.region")
	req, body := SetupRequest("S3", region, "yello")

	signer := SetupSigner(Creds)
	signer.Sign(req, body, "s3", region, time.Now())

	hash := req.Header.Get("X-Amz-Content-Sha256")
	assert.Equal("0e6807fb3a06ab2a6ee35df3d89365b2af1266eb390e9e687e9a500de32571bd", hash)
}

func (suite *S3Suite) TestPresignEmptyBody() {

	assert := suite
	region := viper.GetString("s3main.region")
	req, body := SetupRequest("S3", region, "{}")

	signer := SetupSigner(Creds)
	signer.Presign(req, body, "s3", region, 5*time.Minute, time.Now())

	hash := req.Header.Get("X-Amz-Content-Sha256")
	assert.Equal("", hash)
}

func (suite *S3Suite) TestSignUnsignedpayload() {

	assert := suite
	region := viper.GetString("s3main.region")
	req, body := SetupRequest("S3", region, "yello")

	signer := SetupSigner(Creds)
	signer.Presign(req, body, "s3", region, 5*time.Minute, time.Now())

	hash := req.Header.Get("X-Amz-Content-Sha256")
	assert.Equal("", hash)
}

func (suite *S3Suite) TestSignWithRequestBody() {

	assert := suite
	signer := v4.NewSigner(Creds)

	expectBody := []byte("abc123")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		assert.Nil(err)
		assert.Equal(expectBody, b)
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("POST", server.URL, nil)

	_, err = signer.Sign(req, bytes.NewReader(expectBody), "service", "region", time.Now())
	assert.Nil(err)

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *S3Suite) TestSignWithRequestBodyOverwrite() {

	assert := suite
	signer := v4.NewSigner(Creds)

	var expectBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		assert.Nil(err)
		assert.Equal(len(expectBody), len(b))
		w.WriteHeader(http.StatusOK)
	}))

	req, err := http.NewRequest("GET", server.URL, strings.NewReader("invalid body"))

	_, err = signer.Sign(req, nil, "service", "region", time.Now())
	req.ContentLength = 0

	assert.Nil(err)

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func (suite *S3Suite) TestSignWithBodyReplaceRequestBody() {

	assert := suite
	region := viper.GetString("s3main.region")

	req, seekerBody := SetupRequest("S3", region, "{}")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	s := v4.NewSigner(Creds)
	origBody := req.Body

	_, err := s.Sign(req, seekerBody, "s3", "mexico", time.Now())
	assert.Nil(err)
	assert.NotEqual(req.Body, origBody)
	assert.NotNil(req.Body)
}

func (suite *S3Suite) TestSignWithBodyNoReplaceRequestBody() {

	assert := suite
	region := viper.GetString("s3main.region")

	req, seekerBody := SetupRequest("S3", region, "{}")
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	s := v4.NewSigner(Creds, func(signer *v4.Signer) {
		signer.DisableRequestBodyOverwrite = true
	})

	origBody := req.Body

	_, err := s.Sign(req, seekerBody, "s3", "mexico", time.Now())
	assert.Nil(err)
	assert.Equal(req.Body, origBody)
}

func (suite *S3Suite) TestPresignHandler() {

	assert := suite
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:             aws.String("bucket"),
		Key:                aws.String("key"),
		ContentDisposition: aws.String("a+b c$d"),
		ACL:                aws.String("public-read"),
	})

	req.Time = time.Unix(0, 0)
	urlstr, err := req.Presign(5 * time.Minute)

	assert.Nil(err)

	expectedHost := viper.GetString("s3main.endpoint")
	expectedDate := "19700101T000000Z"
	expectedHeaders := "content-disposition;host;x-amz-acl"
	var credentials string = viper.GetString("s3main.access_key") + "/" + "19700101" + "/" 
	credentials = credentials + viper.GetString("s3main.region") + "/" + "s3" + "/" + "aws4_request"
	expectedCred := credentials

	u, _ := url.Parse(urlstr)
	urlQ := u.Query()
	assert.Equal(expectedHost, u.Host)
	assert.Equal(expectedCred, urlQ.Get("X-Amz-Credential"))
	assert.Equal(expectedHeaders, urlQ.Get("X-Amz-SignedHeaders"))
	assert.Equal(expectedDate, urlQ.Get("X-Amz-Date"))
	assert.Equal("300", urlQ.Get("X-Amz-Expires"))

	assert.NotContains(urlstr, "+") // + encoded as %20
}

func (suite *S3Suite) TestStandaloneSignCustomURIEscape() {

	assert := suite
	var credentials string = viper.GetString("s3main.access_key") + "/" + "19700101" + "/" 
	credentials = credentials + viper.GetString("s3main.region") + "/" + "es" + "/" + "aws4_request"
	var expectedauth = "AWS4-HMAC-SHA256 Credential=" + credentials + ", SignedHeaders=host;x-amz-date"
	signer := v4.NewSigner(Creds, func(s *v4.Signer) {
		s.DisableURIPathEscaping = true
	})

	host := "https://subdomain.us-east-1.es.amazonaws.com"
	req, err := http.NewRequest("GET", host, nil)
	assert.Nil(err)

	req.URL.Path = `/log-*/_search`
	req.URL.Opaque = "//subdomain.us-east-1.es.amazonaws.com/log-%2A/_search"

	_, err = signer.Sign(req, nil, "es", "us-east-1", time.Unix(0, 0))
	assert.Nil(err)

	actual := req.Header.Get("Authorization")
	assert.Contains(actual, expectedauth)
}
