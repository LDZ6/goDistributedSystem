FROM golang:alpine AS builder

# Set Go env
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/goodsinfo-client

# Install dependencies
RUN apk --update --no-cache add ca-certificates gcc libtool make musl-dev protoc git

# Build Go binary
COPY go.mod go.sum ./
RUN go mod download 
COPY . .
RUN make tidy build

# Deployment container
FROM scratch

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/goodsinfo-client/goodsinfo-client /goodsinfo-client
ENTRYPOINT ["/goodsinfo-client"]
CMD []
