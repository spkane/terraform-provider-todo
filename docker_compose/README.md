# Todo API usage

```sh
curl -i http://127.0.0.1:8080/
curl -i http://127.0.0.1:8080/ -X POST -H 'Content-Type: application/spkane.todo-list.v1+json' -d "{\"description\":\"message $RANDOM\", \"completed\":false}"
curl -i http://127.0.0.1:8080/ -X POST -H 'Content-Type: application/spkane.todo-list.v1+json' -d "{\"description\":\"message $RANDOM\",\"completed\":false}"
curl -i http://127.0.0.1:8080/ -X POST -H 'Content-Type: application/spkane.todo-list.v1+json' -d "{\"description\":\"message $RANDOM\",\"completed\":false}"
curl -i http://127.0.0.1:8080/3 -X PUT -H 'Content-Type: application/spkane.todo-list.v1+json' -d '{"description":"go shopping","completed":true}'
curl -i http://127.0.0.1:8080/1 -X DELETE
curl -i http://127.0.0.1:8080/3 -X DELETE
curl -i http://127.0.0.1:8080
curl -i http://127.0.0.1:8080/2 -X DELETE
```
