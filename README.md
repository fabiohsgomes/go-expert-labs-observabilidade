# Go-expert Labs Observabilidade & OpenTelemetry

## Executando o projeto

Para subir os containers do projeto, na raiz do projeto, execute o comando abaixo:
```bash
docker compose up -d
```

Para realizar uma requisição para o Serviço A:
- Cep válido:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"cep": "01001000"}' http://localhost:3000/temperaturas
```

- Cep invalido (Formato Incorreto):
```bash
curl -X POST -H "Content-Type: application/json" -d '{"cep": "12345"}' http://localhost:3000/temperaturas
```

- Cep não encontrado:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"cep": "99999999"}' http://localhost:3000/temperaturas
```

Também é possível testar o projeto acessando os links, utlizando o arquivo api/api.http.

## Visualizando os Span

- http://localhost:9411/zipkin/

