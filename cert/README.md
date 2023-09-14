# Certificates

This folder contains certificates and private-keys to start the logic_app server with HTTPS.

Source: <https://github.com/denji/golang-tls>

## Self-Signed Certificate for Development Server

A self-signed certificate is provided for running the development server / debugging.
Generate a self-signed certificate based on RSA 2048 as follows with `openssl`:

1. Generate a private key
2. Generate self-signed(x509) public key based on the private. Specify the server's hostname as `localhost`.

```bash
openssl genrsa -out cert/dev.key.pem 2048
openssl req -new -x509 -sha256 -key cert/dev.key.pem -out cert/dev.cert.pem -days 3650
```

The self-signed development certificate **must not be used** in a production environment!

## Certificate for Production Server

Generate a custom certificate for the domain under which the logic_app server is running in the production environment.
Let a public CA sign this certificate.

**You must not check in any private key other than the development server's private key into version control.**
