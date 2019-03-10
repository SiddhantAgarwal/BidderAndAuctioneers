package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func TestAdRequest(t *testing.T) {
	ts := httptest.NewServer(GetRouter())
	defer ts.Close()
	payload := strings.NewReader(fmt.Sprintf("{\n\t\"ad_placement_id\": \"%s\"\n}", "testing"))
	resp, err := http.Post(ts.URL+"/adrequest", "application/json", payload)
	if err != nil {
		t.Fail()
	} else {
		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode != http.StatusNoContent {
				t.Fail()
			}
		}
	}
}

func TestAdRequestCorrupt(t *testing.T) {
	ts := httptest.NewServer(GetRouter())
	defer ts.Close()
	payload := strings.NewReader("")
	resp, err := http.Post(ts.URL+"/adrequest", "application/json", payload)
	if err != nil {
		t.Fail()
	} else {
		if resp.StatusCode == http.StatusOK {
			t.Fail()
		}
		if resp.StatusCode == http.StatusNoContent {
			t.Fail()
		}
	}
}
