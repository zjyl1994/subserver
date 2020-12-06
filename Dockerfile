FROM golang:buster AS builder
COPY . /code
WORKDIR /code
RUN go build -o subserver .
FROM debian:buster-slim
RUN mkdir -p /app/data
VOLUME /app/data
COPY --from=builder /code/subserver /app/subserver
COPY --from=builder /code/data.json /app/data/data.json
EXPOSE 8080
CMD ["/app/subserver"]
