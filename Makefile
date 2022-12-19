run:
	go run ./cmd/cripta/*.go

build:
	go build ./cmd/cripta

test:
	go test -v --cover ./...

release:
	@echo "Enter the release version (format vx.x.x).."; \
	read VERSION; \
	git tag -a $$VERSION -m "Releasing "$$VERSION; \
	git push origin $$VERSION