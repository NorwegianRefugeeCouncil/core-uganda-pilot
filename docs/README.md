
# Run the server

1. Start the docker-compose

`docker-compose -f ./deployments/webapp.docker-compose.yaml`

2. Run the db migrations

`go run . migrate --config configs/config.yaml`

4. Start the server

`go run . serve all --config configs/config.yaml`

5. Start the frontend

`cd web/pwa && npm i && npm start`
