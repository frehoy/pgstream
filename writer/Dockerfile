FROM golang:alpine3.16

WORKDIR /writer

COPY go.mod ./
RUN go mod download

COPY writer.go ./

RUN go build -o /writer

CMD [ "/writer/writer" ]