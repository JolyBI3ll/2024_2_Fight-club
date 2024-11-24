#!/bin/sh

# Проверка, установлен ли OpenSSL
if ! command -v openssl &> /dev/null; then
    echo "OpenSSL не найден. Устанавливаю..."
    apk update && apk add --no-cache openssl
else
    echo "OpenSSL уже установлен."
fi

# Создание сертификатов, если их нет
CERT_DIR="/app/ssl"
CERT_KEY="$CERT_DIR/pootnick.key"
CERT_CRT="$CERT_DIR/pootnick.crt"

if [ ! -f "$CERT_KEY" ] || [ ! -f "$CERT_CRT" ]; then
    echo "Создаю самоподписанные сертификаты..."
    mkdir -p "$CERT_DIR"
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout "$CERT_KEY" -out "$CERT_CRT" \
        -subj "/C=RU/ST=State/L=City/O=Organization/OU=Department/CN=localhost"
else
    echo "Сертификаты уже существуют."
fi
