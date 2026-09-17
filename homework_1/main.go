package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {

	fmt.Println("Имя пользователя:", getEnvName())

	args := os.Args[1:]
	fmt.Println("CLI аргументы:", args)

	fmt.Println("Текущая версия Go:", runtime.Version())
}

func getEnvName() string {
	value, exists := os.LookupEnv("USERNAME")
	if exists {
		return value
	}
	return "отсутствует"
}
