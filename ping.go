package main

import (
	"fmt"
	"net/http"
)

// Pinger проверяет доступность URL.

type Pinger struct {
	client HTTPCLient
}

type HTTPCLient interface {
	Head(url string) (resp *http.Response, err error)
}

// Ping запрашивает указанный URL.
// Возвращает true, если адрес доступен, и false в противном случае.

func (p Pinger) Ping(upl string) bool {
	resp, err := p.client.Head(upl)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func main() {
	client := &http.Client{}
	pinger := Pinger{client}
	url := "https://ya.ru"
	alive := pinger.Ping(url)
	fmt.Println(url, "is alive =", alive)
}
