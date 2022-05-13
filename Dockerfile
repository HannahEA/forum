FROM golang:1.17.3-alpine3.13
RUN mkdir /dockerdir
ADD . /dockerdir
WORKDIR /dockerdir
RUN go build -o main .
CMD ["/dockerdir/main"]