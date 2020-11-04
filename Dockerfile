FROM golang:1.14-alpine

WORKDIR /app
COPY . .
RUN go install

ENTRYPOINT ["iot"]