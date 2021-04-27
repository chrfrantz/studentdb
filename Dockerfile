FROM golang:1.16 as builder

LABEL maintainer "user@mail.com"

ADD ./cmd /app/cmd
ADD ./studentdb /app/studentdb
ADD ./mongodb /app/mongodb
ADD ./go.mod /app/go.mod
ADD ./go.sum /app/go.sum

WORKDIR /app/cmd/studentdb

# Alternative if you have exising repo
#RUN go get github.com/...

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o studentdb



FROM scratch

LABEL maintainer="user@mail.com"

WORKDIR /

COPY --from=builder /app/cmd/studentdb /app

CMD ["/app/studentdb"]

