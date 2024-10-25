# syntax=docker/dockerfile:1

FROM golang:1.23
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM scratch
WORKDIR /app
COPY --from=0 /app/templates /app/templates
COPY --from=0 /app/.env /app/.env
COPY --from=0 /app/static /app/static
COPY --from=0 /server /app/server
EXPOSE 8081
CMD ["/app/server"]
