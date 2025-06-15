FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .
RUN go mod tidy \
&& GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C cmd/previsao -o previsao \
&& GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C cmd/valida-cep -o valida-cep

FROM scratch AS servicea
WORKDIR /app
ENV AMBIENTE_PUBLICACAO=DEMO
COPY --from=build /app/cmd/valida-cep/valida-cep valida-cep
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/
ENTRYPOINT ["./valida-cep"]

FROM scratch AS serviceb
WORKDIR /app
ENV AMBIENTE_PUBLICACAO=DEMO
COPY --from=build /app/cmd/previsao/previsao previsao
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/
ENTRYPOINT ["./previsao"]