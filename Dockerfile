# build env
FROM golang:1.20.3-alpine AS build-env
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o main .

# run env
FROM scratch
COPY --from=build-env /app /app
WORKDIR /app
ENTRYPOINT [ "/app/main" ]