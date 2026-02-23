FROM golang:alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o xray-subscriptor

FROM scratch
COPY --from=build /app/xray-subscriptor /xray-subscriptor

ENTRYPOINT ["/xray-subscriptor"]
