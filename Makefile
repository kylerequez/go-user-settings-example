install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss

start dev:
	air

templ:
	@templ generate -watch -proxy=http://localhost:1337

tailwind:
	@npx tailwindcss -i ./src/views/css/input.css -o ./src/public/styles.css --watch
