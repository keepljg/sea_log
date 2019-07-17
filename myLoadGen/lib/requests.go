package lib

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Crawler struct {
	userAgent map[string]string
	url       string
}

func NewCrawler(userAgent map[string]string, url string) *Crawler {
	return &Crawler{
		userAgent: userAgent,
		url:       url,
	}
}

func (this *Crawler) GetFakeHeader(request *http.Request, useAgent map[string]string) {
	for k, v := range useAgent {
		request.Header.Set(k, v)
	}
}

func (this *Crawler) Get(timeout time.Duration) ([]byte, error) {
	var client *http.Client
	var transport *http.Transport
	transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client = &http.Client{Transport: transport, Timeout: timeout}

	request, err := http.NewRequest("GET", this.url, nil)
	if err != nil {
		return nil, err
	}
	this.GetFakeHeader(request, this.userAgent)
	response, err := client.Do(request)
	if err != nil {
		//log.Printf("get html is err: %v", err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK || response.StatusCode == 201 {
		bodyReader := bufio.NewReader(response.Body)
		return ioutil.ReadAll(bodyReader)
	} else if response.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, errors.New("response status code is " + response.Status)
	}
}

// 表单提交
func (this *Crawler) PostFrom(req []byte, timeout time.Duration) ([]byte, error) {

	var (
		err     error
		request *http.Request
		resp    *http.Response
	)
	request, err = http.NewRequest("POST", this.url, bytes.NewReader(req))
	if err != nil {
		return nil, err
	}
	this.GetFakeHeader(request, this.userAgent)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{Timeout: timeout}
	resp, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return body, nil
}

func (this *Crawler) PostJson(req []byte, timeout time.Duration) ([]byte, error) {
	var (
		err     error
		request *http.Request
		resp    *http.Response
	)
	request, err = http.NewRequest("POST", this.url, bytes.NewReader(req))
	if err != nil {
		return nil, err
	}
	this.GetFakeHeader(request, this.userAgent)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{Timeout: timeout}
	resp, err = client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return body, nil
}
