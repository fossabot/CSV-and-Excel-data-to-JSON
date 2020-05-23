FROM golang:latest
ADD . /go/src/app
WORKDIR /go/src/app
ENV PORT=8000
RUN go build main.go
CMD ["./main"]
