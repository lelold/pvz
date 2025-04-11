FROM golang:1.23

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /build ./cmd

EXPOSE 8080
CMD ["/build"]