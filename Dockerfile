FROM golang:1.17.3-alpine3.13

RUN mkdir /docker-practice

ADD . /docker-practice


WORKDIR /docker-practice

RUN go build -o main .

CMD ["/docker-practice/main"]