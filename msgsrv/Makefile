all: docker

msgsrv: main.go
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-s -w' -o $@ .

docker: .docker

.docker: msgsrv Dockerfile
	ifndef GCE_PROJECT_ID
	$(error GCE_PROJECT_ID is not set)
	endif
	docker build -t gcr.io/${GCE_PROJECT_ID}/msgsrv . && touch $@

clean:
	rm -f msgsrv .docker

.PHONY: clean docker
