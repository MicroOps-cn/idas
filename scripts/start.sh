cd $(dirname $0)/..

isRunning=$(docker inspect openldap |jq ".[0].State.Running")
if [ "$isRunning" == "false" ];then
    docker start openldap
fi
make protos idas && (
    cd dist
    ./idas \
      --log.level=debug \
      --log.format=idas \
      --config idas-test.yaml \
      --http.openapi-path=/apidocs.json \
      --http.openapi-path=/apidocs.json \
      --swagger.file-path=./swagger-ui \
      --http.external-url=http://192.168.122.1:8000/
)
