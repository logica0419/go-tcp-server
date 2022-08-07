FROM golang:alpine
WORKDIR /run

RUN go install github.com/cosmtrek/air@latest

COPY go.* .
RUN go mod download

EXPOSE 8080

ENTRYPOINT ["air", "--tmp_dir", "../../tmp", "--build.bin", "/tmp/main", "--build.cmd", "go build -o /tmp/main ."]
