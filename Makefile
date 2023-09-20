db-up:
	@docker-compose up db -d --force-recreate --no-deps

db-down:
	@docker-compose down db

build-oauth-server:
	@go build -o bin/oauth-server ./cmd/oauth-server

build-discord-runner:
	@go build -o bin/discord-runner ./cmd/discord-runner

build-raid-runner:
	@go build -o bin/raid-runner ./cmd/raid-runner

run-discord: build-discord-runner
	@./bin/discord-runner

run-raider: build-raid-runner
	@./bin/raid-runner

run-oauth-server: build-oauth-server
	@./bin/oauth-server

test:
	go test ./... -v

docker-build-oauth-server:
	@docker build -t oauth-server -f ./cmd/oauth-server/Dockerfile .

# Deployment

deploy-oauth-server:
	cp ./deploy/oauth-server.Dockerfile ./Dockerfile
	gcloud run deploy oauth-server --source . --project=zeero-marketplace --region=us-west1
	rm -f ./Dockerfile

deploy-discord-runner:
	cp ./deploy/discord-runner.Dockerfile ./Dockerfile
	gcloud run deploy discord-runner --source . --project=zeero-marketplace --region=us-west1
	rm -f ./Dockerfile

deploy-discord-raider:
	cp ./deploy/discord-raider.Dockerfile ./Dockerfile
	gcloud run deploy raid-runner --source . --project=zeero-marketplace --region=us-west1
	rm -f ./Dockerfile

# Socket
db-socket-production:
	./scripts/fetchCloudSqlProxy.sh && ./bin/cloud_sql_proxy -instances=zeero-marketplace:us-west1:discord-bot=tcp:5544
