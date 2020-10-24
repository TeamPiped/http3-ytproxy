package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

var hclient = &http.Client{
	Transport: &http3.RoundTripper{},
}

func main() {
	fmt.Println("Sending QUIC request to YouTube...")

	resp, err := hclient.Get("https://www.youtube.com/")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
