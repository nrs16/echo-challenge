FROM golang:1.21

LABEL maintainer="noura.r.saad@gmail.com"
WORKDIR /app
COPY . .
RUN go mod download

USER root
RUN make

CMD  ["./main"]