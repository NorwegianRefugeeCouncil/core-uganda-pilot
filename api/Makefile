.PHONY:
up:
	docker-compose -f ../artifacts/docker-compose.yaml up --build -d

.PHONY:
down:
	docker-compose -f ../artifacts/docker-compose.yaml down --remove-orphans

.PHONY:
gen-certs:
	${GOROOT}/bin/go run ${GOROOT}/src/crypto/tls/generate_cert.go --rsa-bits 2048 --host 127.0.0.1,::1,localhost,hydra --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
	mkdir certs
	mv cert.pem certs/cert.pem
	mv key.pem certs/key.pem

.PHONY:
serve:
	go run ./cmd \
		--mongo-database=core \
		--mongo-username=root \
		--mongo-password=example \
		--mongo-hosts=localhost:27017 \
		--redis-address=localhost:6379 \
		--environment=Development \
		--fresh=true \
		--seed=true \
		--tls-disable=true \
		--hydra-admin-url=http://localhost:4445 \
		--hydra-public-url=http://localhost:4444 \
		--login-templates-directory=pkg/apps/login/templates \
		--login-client-id=login \
		--login-client-name=login \
		--login-client-secret=somesecret \
		--login-iam-host=localhost:9000 \
		--login-iam-scheme=http \
		--web-templates-directory=pkg/apps/webapp/templates \
		--web-client-id=webapp \
		--web-client-secret=webapp \
		--web-client-name=webapp \
		--web-iam-host=localhost:9000 \
		--web-iam-scheme=http \
		--web-cms-host=localhost:9000 \
		--web-cms-scheme=http \
		--listen-address=:9000 \
		--base-url=http://localhost:9000

.PHONY:
docker-build:
	docker build -t digitalbackbonecr.azurecr.io/nrc-no/core:latest .

.PHONY:
docker-push: docker-build
	docker push digitalbackbonecr.azurecr.io/nrc-no/core:latest
