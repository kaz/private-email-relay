TAG=asia.gcr.io/private-email-relay/private-email-relay

.PHONY: noop
noop:

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test -race ./...

.PHONY: build
build:
	docker build -t $(TAG) --platform linux/x86_64 .

.PHONY: push
push:
	docker push $(TAG)

.PHONY: deploy
deploy:
	gcloud run deploy --platform managed --project=private-email-relay --region=asia-northeast1 --image $(TAG) private-email-relay
