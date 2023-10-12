package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

/*
Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.
As duas requisições serão feitas simultaneamente para as seguintes APIs:

https://cdn.apicep.com/file/apicep/" + cep + ".json
http://viacep.com.br/ws/" + cep + "/json/

Os requisitos para este desafio são:
- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
- O resultado da request deverá ser exibido no command line, bem como qual API a enviou.
- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.
*/

func isValidZipCode(zip string) bool {
	r, err := regexp.Compile(`^[0-9]{5}-[0-9]{3}$`)

	if err != nil {
		log.Fatal("error occurred during regex compilation:", err)
		return false
	}

	return r.MatchString(zip)
}

func makeRequest(url string) map[string]any {
	formattedUrl := fmt.Sprintf("%v", url)
	request, err := http.NewRequest("GET", formattedUrl, nil)

	if err != nil {
		log.Fatal(fmt.Sprintf("error occurred while attempting to create the request to %v", formattedUrl))
	}

	client := &http.Client{}

	res, err := client.Do(request)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(fmt.Sprintf("error occurred while attempting to make the API call to %v", formattedUrl))
	}

	body, err := io.ReadAll(res.Body)

	var bodyParsed map[string]any
	err = json.Unmarshal(body, &bodyParsed)

	if err != nil {
		log.Fatal(fmt.Sprintf("error encountered while attempting to parse the response body: %v", formattedUrl))
	}

	return bodyParsed
}

func callApiCep(zip string, ch chan map[string]any) {
	res := makeRequest(fmt.Sprintf("https://cdn.apicep.com/file/apicep/%v.json", zip))
	ch <- res
}

func callViaCep(zip string, ch chan map[string]any) {
	res := makeRequest(fmt.Sprintf("https://viacep.com.br/ws/%v/json", zip))
	ch <- res
}

func main() {
	fmt.Println("Please enter the Zip code:")

	var zip string
	_, err := fmt.Scanln(&zip)

	if err != nil {
		log.Fatal(fmt.Sprintf("an error occurred while attempting to read the input. Error: %v", err.Error()))
	}

	if !isValidZipCode(zip) {
		log.Fatal("invalid zip code format. It should be 00000-000.")
	}

	c1 := make(chan map[string]any)
	c2 := make(chan map[string]any)

	go callApiCep(zip, c1)
	go callViaCep(zip, c2)

	select {
	case apiCepRes := <-c1:
		log.Printf("ApiCep: %v", apiCepRes)
	case viaCepRes := <-c2:
		log.Printf("ViaCep: %v", viaCepRes)
	case <-time.After(time.Second):
		log.Fatal("the request to ApiCep and ViaCep exceeded 1 second, resulting in a timeout.")
	}
}
