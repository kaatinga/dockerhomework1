FROM golang:1.15 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.15 as builder

RUN mkdir -p /hello
ADD . /hello
WORKDIR /hello

RUN go test -v ./...

# Добавляем непривилегированного пользователя
RUN useradd -u 10001 helloworld

# Собираем бинарный файл
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
   go build -o /hello ./cmd/hello

FROM scratch

# Не забываем скопировать /etc/passwd с предыдущего стейджа
COPY --from=builder /etc/passwd /etc/passwd
USER helloworld

COPY --from=builder /hello /hello

CMD ["/hello"]