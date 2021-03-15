.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o health-check app/*.go

.PHONY: image
image:
	docker build . -t dms-sms-health-check:${VERSION}
	docker tag dms-sms-health-check:${VERSION} jinhong0719/dms-sms-health-check:${VERSION}.RELEASE

.PHONY: upload
upload:
	docker push jinhong0719/dms-sms-health-check:${VERSION}.RELEASE

.PHONY: stack
stack:
	env VERSION=${VERSION} docker stack deploy -c docker-compose.yml DSM_SMS
