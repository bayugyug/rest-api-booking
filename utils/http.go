package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	pHttpPoolList map[string][]*http.Client
	//ssl certs
	poolSSLCerts *x509.CertPool
)

const (
	pHttpPoolMax  = 100
	pHttpPoolName = "http-get-post"
)

//HttpInit initialize the http global
func HttpInit() {

	//init certs
	poolSSLCerts = x509.NewCertPool()
	poolSSLCerts.AppendCertsFromPEM(pemCerts)

	//init others here
	pHttpPoolList = make(map[string][]*http.Client)
	for i := 1; i <= pHttpPoolMax; i++ {
		httpClient := &http.Client{
			Timeout: time.Duration(30000 * time.Millisecond),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					RootCAs:            poolSSLCerts},
				Dial: (&net.Dialer{
					Timeout:   time.Duration(30000 * time.Millisecond),
					KeepAlive: 1,
				}).Dial,
				TLSHandshakeTimeout: 30 * time.Second,
			},
		}
		pHttpPoolList[pHttpPoolName] = append(pHttpPoolList[pHttpPoolName], httpClient)
	}
}

//HttpPost send request to remote end-point-urls (POST)
func HttpPost(url, body string, hdrs map[string]string) (string, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return "", -2, err
	}
	pool := pHttpPoolName
	p := rand.Intn(len(pHttpPoolList[pool]))
	if len(pHttpPoolList[pool]) <= 0 || pHttpPoolList[pool][p] == nil {
		return "", -2, fmt.Errorf("ERROR: httpPoolPick failed to get 1")
	}

	//settings
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for hdk, hdv := range hdrs {
		req.Header.Add(hdk, hdv)
	}
	resp, err := pHttpPoolList[pool][p].Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return "", -1, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	// read the body
	return strings.TrimSpace(string(contents)), resp.StatusCode, nil
}

//HttpGet send request to remote end-point-urls (GET)
func HttpGet(url string, hdrs map[string]string) (string, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", -2, err
	}
	pool := pHttpPoolName
	p := rand.Intn(len(pHttpPoolList[pool]))
	if len(pHttpPoolList[pool]) <= 0 || pHttpPoolList[pool][p] == nil {
		return "", -2, fmt.Errorf("ERROR: httpPoolPick failed to get 1")
	}
	//settings
	for hdk, hdv := range hdrs {
		req.Header.Add(hdk, hdv)
	}
	resp, err := pHttpPoolList[pool][p].Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return "", -1, err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	// read the body
	return strings.TrimSpace(string(contents)), resp.StatusCode, nil
}
