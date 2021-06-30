docker:
	sudo systemctl start docker

docker-down:
	sudo systemctl stop docker

spinup: docker
	sudo docker-compose -f artifacts/docker-compose.yaml up -d

spindown:
	sudo docker-compose -f artifacts/docker-compose.yaml down

serve:
	go run ./cmd --mongo-database=core --mongo-username=root --mongo-password=example --keycloak-base-url=http://localhost:8080 --keycloak-realm-name=nrc --keycloak-client-id=api --keycloak-client-secret=e6486272-039d-430f-b3c7-47887aa9e206 &> /dev/null &

stop-serve:
	kill `pgrep cmd`

test:
	go test ./pkg/testing


