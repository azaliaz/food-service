FROM golang:latest

RUN apt-get update
RUN apt install -y libpcap-dev

RUN mkdir -p /usr/src/build
WORKDIR /usr/src/build

COPY go.mod go.sum ./

COPY . /usr/src/build

RUN go build -o main ./cmd/server/main.go

EXPOSE 3000
CMD [ "./main" ]