FROM golang:1.18.4-alpine3.16 AS builder
COPY . /go/src/github.com/flyer103/buildkit-demo
WORKDIR /go/src/github.com/flyer103/buildkit-demo
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /mocker ./cmd/frontend-mockerfile/main.go

FROM scratch
COPY --from=builder /mocker /bin/mocker
ENTRYPOINT [ "/bin/mocker" ]
