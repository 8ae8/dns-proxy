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
      -v /path/to/config.json:/dns-proxy/config.json \
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
if it didn't work, try
```shell script
sudo lsof -i :53
```
then try to stop services listening to port 53(for example `sudo systemctl stop dbsmasq`) 

### Configuration

You create configuration like sample file existing named `config.sample.json`.
> Remember setup domain value as your server IP address so clients resolve the 
> host name with your server address. 

sample config file:

```json
{
  "dnsServerHost": "0.0.0.0:53",
  "dnsServers": [
    "8.8.8.8:53"
  ],
  "domains": {
    "test.com": "127.0.0.1"
  },
  "servers": [
    {
      "scheme": "http",
      "host": "0.0.0.0",
      "port": "80"
    },
    {
      "scheme": "https",
      "host": "0.0.0.0",
      "port": "443"
    }
  ]
}
```

1. dnsServerHost > `0.0.0.0:53`: dns server address
1. dnsServers > `["8.8.8.8:53"]`: secondary dns servers so be able to resolve all
other domains as it supposed to be resolved
1. domains > {`domain`: `your server ip`}: remember to set the server IP address
which you want to proxy request from. e.g. {"test.com": "192.168.1.6"}
`test.com` is the domain of website which you want to proxy through server and 
`192.168.1.6` would be IP address of the proxy server.
1. servers: would be proxy servers of service which enables configurations 
of exposing port and IP address for further configurations like putting nginx as
main router of requests.
