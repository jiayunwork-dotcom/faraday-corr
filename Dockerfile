FROM golang:1.21-alpine
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /faraday-corr .
EXPOSE 8080
CMD ["/faraday-corr", "-http", ":8080"]
