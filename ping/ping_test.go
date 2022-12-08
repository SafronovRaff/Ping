package ping

import (
	"net/http"
	"strconv"
	"strings"
)

/*
	func TestPing(t *testing.T) {
		client := &http.Client{}
		pinger := Pinger{client}
		got := pinger.Ping("https://www.google.com")
		if !got {
			t.Errorf(" https://www.google.com/404")
		}
		got = pinger.Ping("https://www.google.com/404")
		if got {
			t.Errorf("Expected google.com/404 to be unvailable")
		}
	}
*/
type MockClient struct{}

// Head возвращает http ответ со статусом указанным в url.
// пример:
// url = https://ya.ru/200 -> статус = 200
// url = https://ya.ru/404 -> статус = 404
func (client *MockClient) Head(url string) (resp *http.Response, err error) {
	parts := strings.Split(url, "/")
	last := parts[len(parts)-1]
	statusCode, err := strconv.Atoi(last)
	if err != nil {
		return nil, err
	}
	resp = &http.Response{StatusCode: statusCode}
	return resp, nil
}
