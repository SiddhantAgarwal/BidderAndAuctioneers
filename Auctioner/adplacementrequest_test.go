package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestIndex(t *testing.T) {
	ts := httptest.NewServer(GetRouter())
	defer ts.Close()
	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fail()
	} else {
		if res.StatusCode != http.StatusOK {
			t.Fail()
		}
	}
}
func TestAdPlacementRequest(t *testing.T) {
	ts := httptest.NewServer(GetRouter())
	defer ts.Close()
	payload := strings.NewReader(fmt.Sprintf("{\n\t\"ad_placement_id\": \"%s\"\n}", "testing"))
	timeout := time.Duration(210 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Post(ts.URL+"/adplacement", "application/json", payload)
	if err != nil {
		errObj := err.(*url.Error)
		if errObj.Timeout() {
			fmt.Println("Timedout")
		}
		t.Fail()
	} else {
		if res.StatusCode != http.StatusOK {
			t.Fail()
		}
		buf := make([]byte, res.ContentLength)
		_, err = res.Body.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			if err.Error() != "EOF" {
				t.Fail()
			}
		}
	}
}

func TestAdPlacementRequestCorrupt(t *testing.T) {
	ts := httptest.NewServer(GetRouter())
	defer ts.Close()
	payload := strings.NewReader("")
	timeout := time.Duration(210 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Post(ts.URL+"/adplacement", "application/json", payload)
	if err == nil {
		if res.StatusCode == http.StatusOK {
			t.Fail()
		}
	}
}
