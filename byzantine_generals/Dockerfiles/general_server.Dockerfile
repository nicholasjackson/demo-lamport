FROM golang:1.23.6

COPY . /app

WORKDIR /app/general
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/general-server

FROM debian:bookworm-slim

COPY --from=0 /bin/general-server /bin/general-server

ENTRYPOINT ["/bin/general-server"]