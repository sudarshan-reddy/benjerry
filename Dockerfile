FROM golang:1.10.0-alpine AS go
RUN apk --no-cache add git
WORKDIR /go/src/github.com/sudarshan-reddy/benjerry

COPY . /go/src/github.com/sudarshan-reddy/benjerry

RUN go build -ldflags "-X 'main.buildTimestamp=$(date '+%b %d %Y %T')' -X main.commitID=`git describe --match=NeVeRmAtCh --always --abbrev --dirty`"

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=go /go/src/github.com/sudarshan-reddy/benjerry .

ENTRYPOINT ["./benjerry"]
