build: 
	@echo "Building Windows Target"
	@GOOS=windows go build -ldflags="-s -w" -o ./out/win/ ./cmd/...
	@echo "Building Linux Target"
	@GOOS=linux go build -ldflags="-s -w" -o ./out/linux/ ./cmd/... 
	@echo "Building Macos Target"
	@GOOS=darwin go build -ldflags="-s -w" -o ./out/macos/ ./cmd/...

run: build
	@echo "Running"
	@./out/linux/ts-go -d -i input.txt -o output.txt

