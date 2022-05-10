FROM  golang:1.17-alpine3.15
RUN mkdir /dockerdir
ADD . dockerdir
WORKDIR /dockerdir
RUN go build -o main .
CMD ["/dockerdir/main"]