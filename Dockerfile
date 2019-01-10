FROM golang:latest
ADD . /go/src/app
WORKDIR /go/src/app
COPY get.sh /go/src/app
RUN bash get.sh
ENV PORT=8000
CMD ["go", "build", "parser.go"]
CMD ["./parser"]
