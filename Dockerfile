
ARG GO_VERSION=1.21.6
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Copia o go.mod e faz o download das dependencias.
COPY go.mod go.sum ./
RUN go mod download

# Copia o código da aplicação e compila o binario.
COPY . .

ENV PORT :8000

RUN CGO_ENABLED=0 GOOS=linux go build -o server
################################################

### Step 2: Copiar o binario do stage anterior para a imagem final.
FROM scratch

# Copia apenas o binario gerado no stage anterior.
COPY --from=builder /app/server /

# Expose the port
EXPOSE 8000

# Define o ponto de entrada para o container como /server.
# O binario será executado quando o container for iniciado.
ENTRYPOINT ["/server"]