FROM golang:1.15 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.15 as tester

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /hello
ADD . /hello
WORKDIR /hello

RUN go test -v ./...

FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /hello
ADD . /hello
WORKDIR /hello

# Добавляем непривилегированного пользователя
RUN useradd -u 10001 helloworld

# Собираем бинарный файл
RUN make

FROM scratch as running

# Не забываем скопировать /etc/passwd с предыдущего стейджа
COPY --from=builder /etc/passwd /etc/passwd
USER helloworld

COPY --from=builder /app /app
#RUN chmod +x /hello

# Указываем порт
ENV PORT=8080

# Шарим порт
EXPOSE ${PORT}/tcp

# Запускаем
CMD ["/app"]