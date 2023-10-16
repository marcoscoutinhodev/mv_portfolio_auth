FROM golang:latest as builder

RUN apt update && \
  apt upgrade -y

WORKDIR /app

COPY . .
RUN go mod tidy && \
  GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./main.go 

FROM scratch as prd

WORKDIR /server
COPY --from=builder /app/server .

ENTRYPOINT ["./server"]
