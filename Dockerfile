FROM golang:1.10.4
ADD . /go/src/app
WORKDIR /go/src/app
RUN mkdir -p /go/src/app/data
COPY get.sh /go/src/app
RUN bash get.sh
ENV PORT=8000
CMD ["go", "run", "parser.go"]
