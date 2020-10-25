package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/lucas-clemente/quic-go/http3"
)

// http/3 client
var h3client = &http.Client{
	Transport: &http3.RoundTripper{},
}

// http/2 client
var h2client = &http.Client{}

// user agent to use
var ua = "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101"

func genericHTTPProxy(w http.ResponseWriter, req *http.Request) {

	q := req.URL.Query()
	host := q.Get("host")
	q.Del("host")

	if len(host) <= 0 {
		host = getHost(req.URL.Path)
	}

	if len(host) <= 0 {
		log.Panic("No host in query parameters.")
	}

	proxyURL, err := url.Parse("https://" + host + strings.Replace(req.URL.Path, "/ggpht", "", 1))

	if err != nil {
		log.Panic(err)
	}

	proxyURL.RawQuery = q.Encode()

	request, err := http.NewRequest("GET", proxyURL.String(), nil)

	copyHeaders(req.Header, request.Header)
	request.Header.Set("User-Agent", ua)

	if err != nil {
		log.Panic(err)
	}

	var client *http.Client

	if strings.HasPrefix(req.URL.Path, "/videoplayback") { // https://github.com/lucas-clemente/quic-go/issues/2836
		client = h2client
	} else {
		client = h3client
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Panic(err)
	}

	copyHeaders(resp.Header, w.Header())

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func copyHeaders(from http.Header, to http.Header) {
	// Loop over header names
	for name, values := range from {
		// Loop over all values for the name.
		for _, value := range values {
			to.Set(name, value)
		}
	}
}

func getHost(path string) (host string) {

	host = ""

	if strings.HasPrefix(path, "/vi/") {
		host = "i.ytimg.com"
	}

	if strings.HasPrefix(path, "/ggpht/") {
		host = "yt3.ggpht.com"
	}

	if strings.HasPrefix(path, "/a/") {
		host = "yt3.ggpht.com"
	}

	return host
}

func main() {
	http.HandleFunc("/videoplayback", genericHTTPProxy)
	http.HandleFunc("/vi/", genericHTTPProxy)
	http.HandleFunc("/a/", genericHTTPProxy)
	http.HandleFunc("/ggpht/", genericHTTPProxy)
	listener, err := net.Listen("unix", "http-proxy.sock")
	if err != nil {
		fmt.Println("Failed to bind to UDS, falling back to TCP/IP")
		fmt.Println(err.Error())
		http.ListenAndServe(":8080", nil)
	} else {
		http.Serve(listener, nil)
	}
}
