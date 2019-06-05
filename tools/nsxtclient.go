package tools

import (
	"crypto/tls"
	"io/ioutil"
	"ipManager/config"
	"log"
	"net/http"
)

type Request struct {
	Methoud  string
	Url      string
	UserName string
	PassWord string
}

var cfg = config.Cfg
var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func NewRequest() *Request {
	tmp := Request{
		Methoud:  "GET",
		Url:      " ",
		UserName: cfg.Username,
		PassWord: cfg.Password,
	}
	return &tmp
}

func (this *Request) ClientGet() []byte {
	request, err := http.NewRequest(this.Methoud, this.Url, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	request.SetBasicAuth(this.UserName, this.PassWord)

	client := http.Client{
		Transport: tr,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer resp.Body.Close()
	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return tmp
}
