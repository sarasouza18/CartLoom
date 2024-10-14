FROM golang:1.21-alpine

WORKDIR /app

# Copiar go.mod e go.sum e baixar as dependências
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copiar o restante do código
COPY . .

# Compilar a aplicação
RUN go build -o /cartloom

# Definir o comando para rodar a aplicação
CMD ["/cartloom"]
