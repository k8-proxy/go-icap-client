The Golang equivalent of c-icap-client

**Making Tls InsecureSkipVerify call**

```go

  req, err := ic.NewRequestTLS(ic.MethodRESPMOD, "icap://<host>:<port>/<path>", httpReq, httpResp,"tls")

  if err != nil {
    log.Fatal(err)
  }

  client := &ic.Client{
		Timeout: 5 * time.Second,
	}

  resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

```



