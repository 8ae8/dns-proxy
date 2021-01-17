# DNS Proxy

DNS proxy is a type of server you can connect to which spoofs your location, so 
you can access content which has been restricted in the region you're actually 
in.

## Run

### go build

If you have goland installed in your system, simply download packages and run 
application by: 

```shell script
go get ./...
go run main.go
```

### docker

1. build docker

    ```shell script
    docker build -t dns-proxy .
    ```

1.  create docker container:

    ```shell script
    docker run -itd \
      -p 53:53/udp \
      -p 80:80 \
      -p 443:443 \
      --name=dns dns-proxy
    ```
    > make sure port 53, 80 and 443 are free for service to listen

---

* If you're getting error for port 53 that is already in use try:
```shell script
sudo systemctl stop systemd-resolved
```