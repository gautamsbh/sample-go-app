# Build stage to build the application executable
FROM golang:1.19.2-alpine as build
LABEL Author "gautam.singh<gautam.singh.abes@gmail.com>"
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -o main main.go

# Final stage build having main executable artifact from stage build
FROM alpine:latest
WORKDIR /usr/bin
COPY --from=build /app/main /usr/bin
RUN chmod +x main
CMD ["main"]