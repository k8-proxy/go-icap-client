The Golang equivalent of c-icap-client:

******OPTIONS******

```go
 -i icap_servername
              The hostname of the icap server,port,service
              icap://<ICAP-Server-IP>:<port>/<service>
 -f filename
              Send this file to the icap server.
 -r filename
              Save output to this file.
 -c check response and file 
   check if icap server return response and file,if true must return response and file
   The default is false

```

**Initiating a plain ICAP request**

```go

./go-icap-client -i=icap://<ICAP-Server-IP>:<port>/<service> -f<input file path>
 -r<out file path> -c<out file path>
example : go run go-icap-client.go  -i=icap://34.242.219.224:1344/gw_rebuild -f=test.pdf -r=retest.pdf -c=true
return healthy server if file return 
example : go run go-icap-client.go  -i=icap://34.242.219.224:1344/gw_rebuild -f=test.pdf -r=retest.pdf -c=false
return healthy server even file not found
```

**Initiating a Secure ICAP request**

```go

./go-icap-client -i=icaps://<ICAP-Server-IP>:<port>/<service>
example : go run go-icap-client.go  -i=icaps://34.242.219.224:1345/gw_rebuild -f=test.pdf -r=retest.pdf -c=true

```



