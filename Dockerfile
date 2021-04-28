FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 1323/tcp

CMD ["go", "run", "main.go"]
