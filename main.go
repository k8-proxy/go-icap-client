package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	ic "github.com/haitham911/icap-client"
)

func main() {
	fmt.Println("start")
	host := "34.242.219.224"
	result := Clienticap(host)
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
func Clienticap(server string) string {
	//ic.SetDebugMode(true)
	var requestHeader http.Header
	host := server

	port := "1345"
	service := "gw_rebuild"
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
	//	req, err := ic.NewRequest(ic.MethodRESPMOD, icap, nil, httpResp) req without tls
	req, err := ic.NewRequestTLS(ic.MethodRESPMOD, icap, nil, httpResp, "tls")

	if err != nil {
		fmt.Println(err)
		return "icap error: " + err.Error()

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

	p := new(strings.Builder)

	io.Copy(p, resp.ContentResponse.Body)

	return "0"
}
