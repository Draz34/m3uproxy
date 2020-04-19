# M3U Proxy

Proxy written in Go. It will process an M3U response and transform the all channel addresses with new addresses. 
It was built taking with the purpose that all traffic goes through it.

## Endpoints

* GET /
```bash
$> curl http://${host}:${port}

Welcome to m3u proxy
```

* GET /ping
```bash
$> curl http://${host}:${port}/ping

pong
```

* GET /channels/{username}/{password}
```bash
$> curl http://localhost:9090/channels/{username}/{password}

#EXTM3U
#EXTINF:-1 tvg-id="ABC FHD" tvg-name="ABC FHD" tvg-logo="http://....png" group-title="group A"
http://localhost:9090/channels/{username}/{password}/5287
#EXTINF:-1 tvg-id="ABC HD" tvg-name="ABC HD" tvg-logo="http://....png" group-title="group A"
http://localhost:9090/channels/{username}/{password}/984
...
```

* GET /channels/{username}/{password}/{id}
```bash
$> curl http://localhost:9090/channels/{username}/{password}/984

... stream ...
```

* GET /channels/{username}/{password}/info/{id}
```bash
$> curl http://localhost:9090/channels/{username}/{password}/info/984

{
    Id: "984",
    Source: {
        Scheme: "http",
        Opaque: "",
        User: null,
        Host: "<originalHost>:<originalPort>,
        Path: "<original request uri to channel>",
        RawPath: "",
        ForceQuery: false,
        RawQuery: "",
        Fragment: ""
    }
}
``` 

## FAQ

### Requirements:
* Golang >= 1.11
* docker (if you whish to build the container)

**How to build the proxy**
```bash
go build m3uproxy/main.go
```

**How to run locally with config file**
Requires config file. Look for example in config/config-dev.yml
```bash
go run m3uproxy/main.go -file <path to config file>
```

**How to run locally with environment variables**
This is useful when running in a docker container
```bash
export M3U_PROXY_PORT="9090"
export M3U_PROXY_HOSTNAME="localhost"
export M3U_PROXY_CHANNELS_URL="<valid url to m3u list>"
go run m3uproxy/main.go

#or 

docker run -d \
--name m3uproxy \
-p 9090:9090 \
-v /data/m3uproxy/db:/var/lib/mysql \
-e WEBROOT=/srv/www \
-e MARIADB_ROOT_PASSWORD="root" \
-e M3U_PROXY_HOSTNAME="{my.proxy.com}" \
-e M3U_PROXY_XTREAM_PORT="{7713}" \
-e M3U_PROXY_ADMIN_LOGIN="{my_admin_login}" \
-e M3U_PROXY_ADMIN_PASSWORD="{my_admin_password}" \
-e M3U_PROXY_XTREAM_HOSTNAME="{iptv.server.com}" \
-e M3U_PROXY_XTREAM_USERNAME="{User}" \
-e M3U_PROXY_XTREAM_PASSWORD="{Password}" \
-e M3U_PROXY_XTREAM_VERSION="2.0" \
-e M3U_PROXY_CHANNELS_URL="{http://iptv.server.com/file.m3u}" \
draz34/m3uproxy:xtream-codes-api
 
```
