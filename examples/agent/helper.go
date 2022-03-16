package main

import "os"

func getPrivateKey() []byte {
	k, err := os.ReadFile("keys/private_key.pem")
	if err != nil {
		panic(err)
	}
	return k
}
