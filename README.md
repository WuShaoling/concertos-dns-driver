## Dynamically modify dns resolution content 

written in go language

### Build

1. go build -o dns-driver main.go
2. docker build . -t dns-driver:latest

### Run(docker)
 docker run -d --name dnsmasq --net=host --dns=223.5.5.5 --dns=223.6.6.6  dns-driver:latest

### Restapi

1. ADD :
    http://localhost:40001/domain/

    ```
    PUT /domain/ HTTP/1.1
    Content-Length: 43
    Host: localhost:40001
    Content-Type: application/json
    
    {"IP":"127.0.0.1", "Domain":"localhost"}
    ```
    
2. DELETE :
    http://localhost:40001/domain/{domain-name}
    ```
    DELETE /domain/ HTTP/1.1
    ```

4. GETALL :
    http://localhost:40001/domain/
    ```
    GET /domain/ HTTP/1.1
    ```
