FROM golang:1.17.6 as builder

COPY . /app
WORKDIR /app
RUN go mod download

RUN go build -o /jaeger-example .

EXPOSE 8080

CMD [ "/jaeger-example" ]