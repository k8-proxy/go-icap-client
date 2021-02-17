package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"time"

	ic "github.com/k8-proxy/icap-client"
)

type ReqParam struct {
	host     string
	port     string
	scheme   string
	services string
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		fmt.Println("error : args not found")
		os.Exit(1)
	}
	newreq, err := parsecmd(argsWithoutProg[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// "34.242.219.224"
	result := Clienticap(*newreq)
	if result != "0" {
		fmt.Println("not healthy server")
		os.Exit(1)

	} else {
		fmt.Println("healthy server")

		os.Exit(0)

	}
	fmt.Println("end")

}

//Clienticap icap client req
func Clienticap(newreq ReqParam) string {
	//ic.SetDebugMode(true)
	var requestHeader http.Header
	host := newreq.host

	port := newreq.port
	service := newreq.services //"gw_rebuild"
	fmt.Println("ICAP Scheme :" + newreq.scheme)
	fmt.Println("ICAP Server :" + host)
	fmt.Println("ICAP Port :" + port)
	fmt.Println("ICAP Service :" + strings.Trim(service, "/"))

	timeout := time.Duration(35000) * time.Millisecond

	handler := func(w http.ResponseWriter, r *http.Request) {
		// grab the generated receipt.pdf file and stream it to browser
		streamPDFbytes, err := ioutil.ReadFile("test.pdf")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		b := bytes.NewBuffer(streamPDFbytes)
		w.Header().Set("Content-type", "application/pdf")

		if _, err := b.WriteTo(w); err != nil { // <----- here!
			fmt.Fprintf(w, "%s", err)
		}

		w.Write([]byte("PDF Generated"))

	}

	reqtest := httptest.NewRequest("GET", "http://servertestpdf.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, reqtest)

	httpResp := w.Result()

	icap := "icap://" + host + ":" + port + "/" + service

	var req *ic.Request
	var reqerr error

	if newreq.scheme == "icaps" {
		req, reqerr = ic.NewRequestTLS(ic.MethodRESPMOD, icap, nil, httpResp, "tls")
	} else {
		req, reqerr = ic.NewRequest(ic.MethodRESPMOD, icap, nil, httpResp)
	}
	if reqerr != nil {
		fmt.Println(reqerr)
		return "icap error: " + reqerr.Error()

	}

	req.ExtendHeader(requestHeader)
	client := &ic.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "resp error: " + err.Error()

	}

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Status)
	fmt.Println(resp.Header)
	fmt.Println(resp.ContentRequest)
	fmt.Println(resp.ContentResponse)

	p := new(strings.Builder)

	io.Copy(p, resp.ContentResponse.Body)

	return "0"
}

func parsecmd(pram string) (*ReqParam, error) {
	// We'll parse this example URL, which includes a
	// scheme, authentication info, host, port, path,
	s := pram

	// Parse the URL and ensure there are no errors.
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	host, port, _ := net.SplitHostPort(u.Host)
	scheme := ""
	if host == "" {
		log.Fatal("invalid host")
	}
	if port == "" {
		log.Fatal("invalid port")
	}
	if u.Scheme == "" {
		log.Fatal("invalid scheme")
	}
	if u.Scheme == "icaps" || u.Scheme == "icap" {
		scheme = u.Scheme
	} else {
		log.Fatal("invalid scheme")
	}
	if u.Path == "" {
		log.Fatal("invalid services")
	}
	req := &ReqParam{
		host:     host,
		port:     port,
		scheme:   scheme,
		services: strings.Trim(u.Path, "/"),
	}

	return req, err

}
