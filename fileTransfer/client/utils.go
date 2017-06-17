package main

import (
	"fmt"
	"os"
)

// Habilitar mensagens de depuração
const DEBUG = false

// Quando um erro obtido
const ERRO_FATAL = "Erro fatal: %s"


func debug(s string) {
	if DEBUG {
		fmt.Println("DEBUG: " + s)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, ERRO_FATAL, err.Error())
		os.Exit(1)
	}
}