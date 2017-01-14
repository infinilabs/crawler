## Security Gopa

### Gopa supported security communication

* generate cert files

```
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.crt -days 365
openssl rsa -in server.key -out server.public
```

or

```
go run $GOROOT/src/crypto/tls/generate_cert.go --host 127.0.0.1
```


* you should have three files:

    `server.crt`    `server.key`

* put them into  a folder, like `cert_files`

* start gopa with the parameter `cert_path`:

    `./gopa -cert_path=cert_files`


