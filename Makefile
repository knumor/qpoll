build:
	@npx tailwindcss -i views/css/styles.css -o public/styles.css
	@go build -ldflags "-s -w -X main.Env=prod" -o bin/qpoll main.go 
	@upx --lzma -6 bin/qpoll

test:
	@go test -v ./...
	
run: build
	@./bin/qpoll

tailwind:
	@npx tailwindcss -i views/css/styles.css -o public/styles.css --watch

