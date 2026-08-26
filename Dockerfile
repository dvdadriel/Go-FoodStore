FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/go-food-store .

FROM alpine:3.20
RUN adduser -D -u 10001 app
COPY --from=build /out/go-food-store /usr/local/bin/go-food-store
USER app
EXPOSE 8080
ENTRYPOINT ["go-food-store"]
