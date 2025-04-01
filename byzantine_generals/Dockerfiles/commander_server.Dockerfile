FROM golang:1.23.6

COPY . /app

WORKDIR /app/commander
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/commander-server

FROM debian:bookworm-slim

COPY --from=0 /bin/commander-server /bin/commander-server

ENTRYPOINT ["/bin/commander-server"]