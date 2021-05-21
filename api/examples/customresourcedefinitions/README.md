```shell
curl -X POST http://localhost:8000/apis/core.nrc.no/v1/customresourcedefinitions --data-binary @examples/customresourcedefinitions/bla.json -H "Content-Type: application/json"
curl -X POST http://localhost:8000/apis/bla.com/v1/blas --data-binary @examples/customresourcedefinitions/bla_payload.json -H "Content-Type: application/json"

```


