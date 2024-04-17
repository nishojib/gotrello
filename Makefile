include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## templ: run the templ generate command
.PHONY: templ
templ:
	templ generate -watch -proxy=http://localhost:${PORT}

## tailwind: run the tailwindcss build command
.PHONY: tailwind
tailwind:
	npx tailwindcss -i ui/static/css/app.css -o ui/static/css/styles.css

## run/web: run the cmd/web application
.PHONY: run/web
run/web:
	go run ./cmd/web -db-dsn=${DB_DSN} -sb-url=${SUPABASE_URL} -sb-key=${SUPABASE_KEY}

# run/seed: seeds the database with data from testdata/fixtures
.PHONY: run/seed
run/seed:
	go run ./cmd/seed -db-dsn=${DB_DSN}


## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./internal/data/migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./internal/data/migrations -database ${DB_DSN} up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down:
	@echo 'Running down migrations...'
	migrate -path ./internal/data/migrations -database ${DB_DSN} down

## db/migrations/force n=$1: forcefully apply the database migration
.PHONY: db/migrations/force
db/migrations/force:
	@echo 'Running force migrations...'
	migrate -path ./internal/data/migrations -database ${DB_DSN} force ${n}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/web: build the cmd/web application
.PHONY: build/web
build/web:
	@echo 'Building cmd/web...'
	go build -ldflags='-s -w' -o=./bin/web ./cmd/web
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o=./bin/linux_amd64/web ./cmd/web

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh gotrello@${PRODUCTION_HOST_IP}

## production/deploy/web: deploy the web to production
.PHONY: production/deploy/web
production/deploy/web:
	rsync -P ./bin/linux_amd64/web gotrello@${PRODUCTION_HOST_IP}:~
	rsync -rP --delete ./internal/data/migrations gotrello@${PRODUCTION_HOST_IP}:~
	rsync -rP --delete testdata gotrello@${PRODUCTION_HOST_IP}:~
	rsync -P ./remote/production/web.service gotrello@${PRODUCTION_HOST_IP}:~
	rsync -P ./remote/production/Caddyfile gotrello@${PRODUCTION_HOST_IP}:~
	ssh -t gotrello@${PRODUCTION_HOST_IP} '\
		migrate -path ~/migrations -database $$DB_DSN up \
		&& sudo mv ~/web.service /etc/systemd/system/ \
		&& sudo systemctl enable web \
		&& sudo systemctl restart web \
		&& sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'