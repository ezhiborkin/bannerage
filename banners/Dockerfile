FROM golang:latest
RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get install -y postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o banners ./cmd/banners/main.go

CMD ["./banners"]