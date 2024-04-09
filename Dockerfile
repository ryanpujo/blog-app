FROM golang:1.22.1-alpine3.19

WORKDIR /ryanpujo/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /build ./cmd

EXPOSE 8080

CMD [ "/build" ]
