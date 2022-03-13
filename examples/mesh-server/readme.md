# Example Key

```
curl --location --request POST 'http://localhost:8080/agents/key' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "pubkey:v1",
    "public_key": "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCfZSwSf28T9xoM1Hu2OETCLhEa\naEjbmEYXTRvA7aIS9aHg6ToEkz5J2SGxcnlUn\nS08AZMcOKgQSk/0x3f/Q13d9jx\ndLsgEWNJTj4r3iVtXQ2ZxRdeWqVpE8i4gpXo24iggoV8VLYZnUC/rGFbtXEhCgqc\niWEuM4+uLpLeqZAcLwIDAQAB\n-----END PUBLIC KEY-----",
    "expires_at": "2023-03-01T23:00:00+00:00"
}'
```