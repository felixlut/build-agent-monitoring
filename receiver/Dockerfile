FROM golang:1.18-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
RUN go build -o /usr/local/bin/receiver .
# RUN go build -v -o /usr/local/bin/receiver .

EXPOSE 3000/tcp

CMD ["receiver"]
