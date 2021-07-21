FROM golang:1.12.0-alpine3.9
#copy source code
COPY . /app

# set go path
ENV GOPATH=/app

#build docker binary
WORKDIR /app/src/server
RUN go build -o main server.go
EXPOSE 8080
CMD ["./main"]
