# Identity authentication service

### Introduction
Identity authentication service（身份认证服务）
Implement a single sign on system for OAuth 2.0 protocol. Your system can access the platform through OAuth2.0 to achieve single sign on.

### Framework
Based on the Go language development, the overall use of the go kit framework, the transport layer uses go restful to handle the HTTP protocol of the transport layer, and uses gogo/protobuf for serialization/deserialization.

### How to use?
#### Build
```bash
mkdir dist && make idas
```
#### Run
```bash
cd dist &&  ./idas \
      --log.level=debug \
      --log.format=idas \
      --http.external-url=http://127.0.0.1:8000/
```
