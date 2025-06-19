#!/bin/bash

# Constants
PRODUCT="goapiCA"
CN="GoApiCA"

# Settings
DOMAIN=$1
CERT_DAYS=$2

# === STEP 1: Generate CA (only once) ===
if [[ -f "$PRODUCT.pem" && -f "$PRODUCT.key" ]]; then
    echo -e "\n✅ $PRODUCT already exists. Skipping CA generation.\n"
else
    echo -e "\n🔧 Generating $PRODUCT..."
    openssl genrsa -out "$PRODUCT.key" 2048
    openssl req -x509 -new -nodes -key "$PRODUCT.key" -sha256 -days "$CERT_DAYS" -out "$PRODUCT.pem" -subj "/CN=$CN"
    echo -e "✅ $PRODUCT generated.\n"
fi

# === STEP 2: Generate server key & CSR ===
echo -e "🔧 Generating server key and CSR for $DOMAIN..."
openssl genrsa -out "$DOMAIN.key" 2048
echo -e "✅ Key $DOMAIN.key generated."
openssl req -new -key "$DOMAIN.key" -out "$DOMAIN.csr" -subj "/CN=$DOMAIN"
echo -e "✅ Signing request $DOMAIN.csr generated.\n"

# === STEP 3: Generate .ext file with SAN ===
cat > $DOMAIN.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = $DOMAIN
DNS.2 = localhost
EOF

# === STEP 4: Sign cert with CA ===
echo "🔧 Signing certificate..."
openssl x509 -req -in "$DOMAIN.csr" -CA "$PRODUCT.pem" -CAkey "$PRODUCT.key" -CAcreateserial \
-out $DOMAIN.crt -days $CERT_DAYS -sha256 -extfile $DOMAIN.ext

echo -e "✅ Certificate for $DOMAIN created and signed.\n"
