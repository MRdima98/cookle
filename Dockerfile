# syntax=docker/dockerfile:1

FROM golang:1.23
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

FROM scratch
WORKDIR /app
COPY --from=0 /app/index.html /app/index.html
COPY --from=0 /app/.env /app/.env
COPY --from=0 /server /server
EXPOSE 8080
CMD ["/server"]
