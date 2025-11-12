# 🌤️ FC - Golang - Weather API Lab

> Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## 📌 Sobre

Este projeto foi desenvolvido como laboratório da pós-graduação, implementando uma API REST que:

- Recebe um CEP válido de 8 dígitos **brasileiros**
- Consulta a localização via ViaCEP (API brasileira)
- Busca a temperatura atual via WeatherAPI
- Retorna as temperaturas em 3 unidades: Celsius, Fahrenheit e Kelvin

> **⚠️ Limitação**: Apenas CEPs brasileiros são suportados, pois o ViaCEP é uma API exclusiva do Brasil.

## 🏗️ Arquitetura

O projeto segue os princípios da **Clean Architecture** com 3 camadas principais:

- **Handler**: Recebe requisições HTTP (Chi router)
- **UseCase**: Lógica de negócio (validação CEP, conversão temperaturas)
- **Repository**: Integração com APIs externas (ViaCEP, WeatherAPI)

## 🔧 Configuração

### Pré-requisitos

- Go 1.23.5+
- Conta na WeatherAPI (gratuita)

### 1. Configurar Ambiente

```bash
# Clonar o repositório
git clone https://github.com/dabraga/goexpert-weather-api-main.git
cd goexpert-weather-api-main

# Configurar dependências
go mod download

# Configurar API key
export WEATHER_API_KEY=sua-api-key-aqui
```

### 2. Obter API Key da WeatherAPI

1. Acesse: <https://www.weatherapi.com/>
2. Crie uma conta gratuita (1000 calls/mês)
3. Copie sua API key
4. Configure como variável de ambiente

## 🚀 Como Executar

### Desenvolvimento Local

```bash
# Executar aplicação
go run ./cmd/api

# Ou com variável de ambiente
WEATHER_API_KEY=sua-api-key go run ./cmd/api
```

### Testes

```bash
# Todos os testes
go test -v ./...

# Testes com cobertura
go test -v -cover ./...
```

## 📚 API

### 🌐 API em Produção (Cloud Run)

A API está disponível em produção no Google Cloud Run:

**URL**: <https://goexpert-weather-api-main-534292889217.us-central1.run.app>

#### Teste da API Deployada

```bash
# Exemplo de uso com CEP válido
curl "https://goexpert-weather-api-main-534292889217.us-central1.run.app/weather/26140040"

# Formato da resposta esperada:
# {"temp_C":17.2,"temp_F":63,"temp_K":290.2}
```

### Endpoint

```text
GET /weather/{cep}
```

### Respostas

#### ✅ Sucesso (200)

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### ❌ CEP Inválido (422)

```json
{
  "message": "invalid zipcode"
}
```

#### ❌ CEP Não Encontrado (404)

```json
{
  "message": "can not find zipcode"
}
```

### Exemplos de Uso

```bash
# CEP válido (Belford Roxo, RJ)
curl http://localhost:8080/weather/26140040

# CEP inválido (muito curto)
curl http://localhost:8080/weather/123

# CEP não encontrado
curl http://localhost:8080/weather/99999999
```

### CEPs para Teste

- `01310100` - Av. Paulista, São Paulo/SP
- `20040020` - Centro, Rio de Janeiro/RJ
- `26140040` - Belford Roxo/RJ

## 🧪 Testes

### Estrutura de Testes

- **Testes Unitários**: Cada camada testada isoladamente com mocks
- **Testes de Integração**: Endpoint completo com APIs mockadas
- **Cobertura**: 100% dos testes passando

### Executar Testes

```bash
# Todos os testes
go test -v ./...

# Testes unitários
go test -v ./internal/...

# Testes de integração
go test -v ./tests/...

# Com cobertura
go test -v -cover ./...
```

## ☁️ Deploy no Google Cloud Run

### Comandos de Deploy

```bash
# 1. Build da imagem Docker
docker build -t gcr.io/SEU_PROJECT_ID/fc-pos-golang-lab-weather-api .

# 2. Push da imagem para Google Container Registry
docker push gcr.io/SEU_PROJECT_ID/fc-pos-golang-lab-weather-api

# 3. Deploy no Cloud Run
gcloud run deploy fc-pos-golang-lab-weather-api \
  --image gcr.io/SEU_PROJECT_ID/steel-spark-477519-b5 \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars="WEATHER_API_KEY=sua_chave_aqui"
```

## 📊 Conversões de Temperatura

- **Celsius para Fahrenheit**: `F = C * 1.8 + 32`
- **Celsius para Kelvin**: `K = C + 273`

## 📝 Estrutura do Projeto

```bash
/
├── cmd/api/                    # Entry point
├── internal/
│   ├── domain/                 # Entidades e regras de negócio
│   ├── repository/             # Integração com APIs externas
│   ├── usecase/                # Casos de uso
│   ├── handler/                # HTTP handlers
│   └── dto/                    # Data Transfer Objects
├── tests/integration/          # Testes de integração
├── Dockerfile                  # Multi-stage build
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🚨 Troubleshooting

### Erro: "WEATHER_API_KEY é obrigatória"

- Verifique se a variável `WEATHER_API_KEY` está configurada
- Configure como variável de ambiente: `export WEATHER_API_KEY=sua-chave`

### Erro: "can not find zipcode"

- Verifique se o CEP tem 8 dígitos
- Confirme se o CEP existe no ViaCEP
- Teste com CEPs conhecidos (ex: 01310100)

### Erro: "weather not found"

- Verifique se a API key da WeatherAPI está válida
- Confirme se há créditos disponíveis na conta
- Teste a API diretamente no navegador

## 📄 Licença

Este projeto foi desenvolvido para fins acadêmicos.
