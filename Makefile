build:
	@npx tailwindcss -i views/css/styles.css -o public/styles.css
	@go build -ldflags "-X main.Env=prod" -o bin/qpoll main.go 

test:
	@go test -v ./...
	
run: build
	@./bin/qpoll

tailwind:
	@npx tailwindcss -i views/css/styles.css -o public/styles.css --watch

