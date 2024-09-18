FROM golang:1.22-bookworm
RUN apt update 
RUN apt install -y ssh

WORKDIR /app

COPY . .
RUN go mod download
RUN go mod tidy
RUN go mod vendor

RUN make build

EXPOSE 2222

CMD ["./gochat"]
# CMD ["tail", "-f", "/dev/null"]
