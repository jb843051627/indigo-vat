FROM golang:1.22-bookworm
ENV GOPROXY=https://goproxy.cn,direct
ENV GOTOOLCHAIN=local
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/indigo-vat .
EXPOSE 8080
ENTRYPOINT ["/app/indigo-vat"]
