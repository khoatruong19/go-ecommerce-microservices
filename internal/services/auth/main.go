package main

import "github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/config/environment"

func main() {
	env := environment.ConfigAppEnv(environment.Development)
	print("hello ", env)
}
