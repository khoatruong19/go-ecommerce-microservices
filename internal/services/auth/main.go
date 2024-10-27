package main

import (
	"github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/config/environment"
)

func main() {
	env := environment.ConfigAppEnv(environments...)
	print("hello")
}
