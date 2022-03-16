# Example Agent

This code represents a minimal Poly agent

the private key (keys/private_key.pem) is used to register the agent with the Mesh server

## Adding the public key to the Mesh server

The public key needs to be added to the server first before it can be used

use the `/agents/key` endpoint with the following request body:

    {
        "name": "pubkey:v0.1",
        "public_key": "...",
        "expires_at": "2023-03-01T23:00:00+00:00"
    }

## Generating keys

Example RSA PKCS8 keys can be found under `/keys`

To create your own, run:

    # Private key
    openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048

    # Public key
    openssl rsa -pubout -in private.pem -out public_key.pem

    # Private key in pkcs8 format (for Java maybe :D)
    openssl pkcs8 -topk8 -in private.pem -out private_key.pem

    ## nocrypt (Private key does have no password)
    openssl pkcs8 -topk8 -in private.pem -nocrypt -out private_key.pem

