setup:
	pnpm install
	go install github.com/gzuidhof/tygo@latest
	go install github.com/air-verse/air@latest

generate-types:
	cd server && tygo generate

build-client:
	export SITE_URL=$$(grep -m 1 'SITE_URL' ./server/.env | cut -d '=' -f 2) && \
	export REPORTING_ENDPOINT="$$SITE_URL" && \
	pnpm run build

build-server:
	cd server && go build -o overlord .

build: generate-types build-client build-server

watch: 
	$(MAKE) generate-types
	$(MAKE) build-client
	cd server && air -c ./.air.toml

run-server:
	./server/overlord