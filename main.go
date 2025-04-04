package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	brasilUrl string = "https://brasilapi.com.br/api/cep/v1/%s"
	viaUrl    string = "http://viacep.com.br/ws/%s/json/"
)

func cepFunc(url, cep string, ch chan<- map[string]any) {
	req, err := http.Get(fmt.Sprintf(url, cep))
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var response map[string]any
	err = json.Unmarshal(res, &response)
	if err != nil {
		panic(err)
	}
	ch <- response
}

func main() {
	brasilCepch := make(chan map[string]any)
	viaCepch := make(chan map[string]any)

	// Mude o cep se quiser testar outros ceps
	cep := "05541030"
	go cepFunc(brasilUrl, cep, brasilCepch)
	go cepFunc(viaUrl, cep, viaCepch)

	select {
	case msg := <-brasilCepch:
		fmt.Printf("URL: %s RESULT: %+v\n", brasilUrl, msg)
	case msg := <-viaCepch:
		fmt.Printf("URL: %s RESULT: %+v\n", viaUrl, msg)
	case <-time.After(time.Second * 1):
		fmt.Println("TIMEOUT!!")
	}
}
