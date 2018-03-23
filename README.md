## Dynamically modify dns resolution content 

written in go language

### Build

1. go build -o dns-driver main.go
2. docker build . -t dns-driver:latest

### Run(docker)
docker run -d -p 8082:8082 -p 53:53 dns-driver:latest

### Restapi

1. ADD :
    http://localhost:8082/domain/add

    ```
    PUT /domain/add HTTP/1.1
    Content-Length: 43
    Host: localhost:8082
    Content-Type: application/json
    
    {"IP":"127.0.0.1", "Domain":"foo.com"}
    ```
    
2. DELETE :
    http://localhost:8082/domain/delete/{domain-name}

