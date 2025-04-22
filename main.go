package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type apiResponse struct {
	body   string
	source string
	err    error
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipcode"`
}

func validateCEP(cep string) error {
	cleaned := strings.ReplaceAll(cep, "-", "")
	if len(cleaned) != 8 || !regexp.MustCompile(`^[0-9]+$`).MatchString(cleaned) {
		return fmt.Errorf("CEP inválido: %s", cep)
	}
	return nil
}

func fetchAPI(ctx context.Context, url, source string, ch chan<- apiResponse) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- apiResponse{err: err}
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		ch <- apiResponse{err: err}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		ch <- apiResponse{err: fmt.Errorf("HTTP error: %d", res.StatusCode)}
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		ch <- apiResponse{err: err}
		return
	}
	ch <- apiResponse{body: string(body), source: source}
}

func parseResponse(body []byte, source string) (Address, error) {
	var addr Address
	switch source {
	case "VIACEP":
		var viaResp struct {
			Logradouro string `json:"logradouro"`
			Localidade string `json:"localidade"`
			Uf         string `json:"uf"`
			Cep        string `json:"cep"`
		}
		if err := json.Unmarshal(body, &viaResp); err != nil {
			return Address{}, err
		}
		addr = Address{
			Street:  viaResp.Logradouro,
			City:    viaResp.Localidade,
			State:   viaResp.Uf,
			ZipCode: viaResp.Cep,
		}
	case "BrasilAPI":
		var brResp struct {
			Street string `json:"street"`
			City   string `json:"city"`
			State  string `json:"state"`
			Cep    string `json:"cep"`
		}
		if err := json.Unmarshal(body, &brResp); err != nil {
			return Address{}, err
		}
		addr = Address{
			Street:  brResp.Street,
			City:    brResp.City,
			State:   brResp.State,
			ZipCode: brResp.Cep,
		}
	default:
		return Address{}, fmt.Errorf("API desconhecida: %s", source)
	}
	return addr, nil
}

func main() {
	cep := "89010-904"
	if err := validateCEP(cep); err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan apiResponse, 2)
	go fetchAPI(ctx, "https://brasilapi.com.br/api/cep/v1/"+strings.ReplaceAll(cep, "-", ""), "BrasilAPI", ch)
	go fetchAPI(ctx, "http://viacep.com.br/ws/"+cep+"/json/", "VIACEP", ch)

	select {
	case res := <-ch:
		if res.err != nil {
			fmt.Printf("Erro na API %s: %v\n", res.source, res.err)
			return
		}
		addr, err := parseResponse([]byte(res.body), res.source)
		if err != nil {
			fmt.Printf("Erro ao parsear resposta da API %s: %v\n", res.source, err)
			return
		}
		fmt.Printf("Resposta de %s:\nRua: %s\nCidade: %s\nEstado: %s\nCEP: %s\n",
			res.source, addr.Street, addr.City, addr.State, addr.ZipCode)
	case <-ctx.Done():
		fmt.Println("Ocorreu timeout de 1 segundo")
	}
}
