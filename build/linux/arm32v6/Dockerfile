FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/robstrong/beacon-scan

RUN cd /go/src/github.com/robstrong/beacon-scan \
	&& GOOS=linux GOARCH=arm GOARM=6 go build . \
	&& mv beacon-scan /usr/bin/beacon-scan

FROM scratch

COPY --from=builder /usr/bin/beacon-scan /usr/bin/beacon-scan
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "beacon-scan" ]