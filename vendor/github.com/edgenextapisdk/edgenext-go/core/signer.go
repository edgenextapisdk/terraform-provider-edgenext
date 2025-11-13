package core

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io"

	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	DateFormat           = "20060102T150405Z"
	SignAlgorithm        = "SDK-HMAC-SHA256"
	HeaderXDateTime      = "X-Sdk-Date"
	HeaderXHost          = "host"
	HeaderXContentSha256 = "X-Sdk-Content-Sha256"
)

func hmacsha256(keyByte []byte, dataStr string) ([]byte, error) {
	hm := hmac.New(sha256.New, []byte(keyByte))
	if _, err := hm.Write([]byte(dataStr)); err != nil {
		return nil, err
	}
	return hm.Sum(nil), nil
}

// CanonicalRequest Build a CanonicalRequest from a regular request string
func CanonicalRequest(request *http.Request, signedHeaders []string) (string, error) {
	var hexencode string
	var err error
	if hex := request.Header.Get(HeaderXContentSha256); hex != "" {
		hexencode = hex
	} else {
		bodyData, err := RequestPayload(request)
		if err != nil {
			return "", err
		}
		hexencode, err = HexEncodeSHA256Hash(bodyData)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s",
		request.Method,
		CanonicalURI(request),
		CanonicalQueryString(request),
		//CanonicalHeaders(request, signedHeaders),
		//strings.Join(signedHeaders, ";"),
		hexencode,
	), err
}

// CanonicalURI returns request uri
func CanonicalURI(request *http.Request) string {
	pattens := strings.Split(request.URL.Path, "/")
	var uriSlice []string
	for _, v := range pattens {
		uriSlice = append(uriSlice, escape(v))
	}
	urlpath := strings.Join(uriSlice, "/")
	if len(urlpath) == 0 || urlpath[len(urlpath)-1] != '/' {
		urlpath = urlpath + "/"
	}
	return urlpath
}

// CanonicalQueryString build query
func CanonicalQueryString(request *http.Request) string {
	var keys []string
	queryMap := request.URL.Query()
	for key := range queryMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var query []string
	for _, key := range keys {
		k := escape(key)
		sort.Strings(queryMap[key])
		for _, v := range queryMap[key] {
			kv := fmt.Sprintf("%s=%s", k, escape(v))
			query = append(query, kv)
		}
	}
	queryStr := strings.Join(query, "&")
	request.URL.RawQuery = queryStr
	return queryStr
}

func CanonicalHeaders(request *http.Request, signerHeaders []string) string {
	var canonicalHeaders []string
	header := make(map[string][]string)
	for k, v := range request.Header {
		header[strings.ToLower(k)] = v
	}
	for _, key := range signerHeaders {
		value := header[key]
		if strings.EqualFold(key, HeaderXHost) {
			value = []string{request.Host}
		}
		sort.Strings(value)
		for _, v := range value {
			canonicalHeaders = append(canonicalHeaders, key+":"+strings.TrimSpace(v))
		}
	}
	return fmt.Sprintf("%s\n", strings.Join(canonicalHeaders, "\n"))
}

// SignedHeaders 取出头部所有的key -> 转换成小写 -> 排序
func SignedHeaders(r *http.Request) []string {
	var signedHeaders []string
	for key := range r.Header {
		signedHeaders = append(signedHeaders, strings.ToLower(key))
	}
	sort.Strings(signedHeaders)
	return signedHeaders
}

// RequestPayload 读取body
func RequestPayload(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return []byte(""), nil
	}
	bodyByte, err := io.ReadAll(request.Body)
	if err != nil {
		return []byte(""), err
	}
	request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
	return bodyByte, err
}

// StringToSign Create a "String to Sign".
func StringToSign(canonicalRequest string, t time.Time) (string, error) {
	hashStruct := sha256.New()
	_, err := hashStruct.Write([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n%s\n%x",
		SignAlgorithm, t.UTC().Format(DateFormat), hashStruct.Sum(nil)), nil
}

// SignStringToSign Create the HWS Signature.
func SignStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hmsha, err := hmacsha256(signingKey, stringToSign)
	return fmt.Sprintf("%x", hmsha), err
}

// HexEncodeSHA256Hash returns hexcode of sha256
func HexEncodeSHA256Hash(body []byte) (string, error) {
	hashStruct := sha256.New()
	if len(body) == 0 {
		body = []byte("")
	}
	_, err := hashStruct.Write(body)
	return fmt.Sprintf("%x", hashStruct.Sum(nil)), err
}

// AuthHeaderValue Get the finalized value for the "Authorization" header. The signature parameter is the output from SignStringToSign
func AuthHeaderValue(signatureStr, accessKeyStr string) string {
	return fmt.Sprintf("Bearer %s", signatureStr)
}

// Signer access key from edgenext console
type Signer struct {
	AppId     string
	AppSecret string
}

func (s *Signer) sign(request *http.Request) (string, error) {
	var t time.Time
	var err error
	var date string
	if date = request.Header.Get(HeaderXDateTime); date != "" {
		t, err = time.Parse(DateFormat, date)
	}
	if err != nil || date == "" {
		t = time.Now()
		request.Header.Set(HeaderXDateTime, t.UTC().Format(DateFormat))
	}
	// 取出头部key
	signedHeaders := SignedHeaders(request)
	canonicalRequest, err := CanonicalRequest(request, signedHeaders)
	if err != nil {
		return "", err
	}
	stringToSignStr, err := StringToSign(canonicalRequest, t)
	if err != nil {
		return "", err
	}
	signatureStr, err := SignStringToSign(stringToSignStr, []byte(s.AppSecret))
	if err != nil {
		return "", err
	}
	return AuthHeaderValue(signatureStr, s.AppId), nil
}

// Sign SignRequest set Authorization header
func (s *Signer) Sign(request *http.Request) error {
	signatureStr, err := s.sign(request)
	if err != nil {
		return err
	}
	request.Header.Set("X-Auth-Sign", signatureStr)
	return nil
}

func (s *Signer) Verify(request *http.Request) error {
	signatureStr, err := s.sign(request)
	if err != nil {
		return err
	}
	// req sign
	reqAuth := request.Header.Get("X-Auth-Sign")
	if signatureStr != reqAuth {
		return fmt.Errorf("auth failed")
	}
	return nil
}
