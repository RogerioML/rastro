package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

const WSDL = "https://webservice.correios.com.br/service/rastro/Rastro.wsdl"

type rastreiaResponse struct {
	XML  xml.Name `xml:"Envelope"`
	Body struct {
	} `xml:"Body"`
}

func rastreia(wsdl string, objeto string, usuario string, senha string) (string, error) {
	payload := fmt.Sprintf(`
		<x:Envelope
		xmlns:x="http://schemas.xmlsoap.org/soap/envelope/"
		xmlns:res="http://resource.webservice.correios.com.br/">
		<x:Header/>
		<x:Body>
			<res:buscaEventos>
				<res:usuario>` + usuario + `</res:usuario>
				<res:senha>` + senha + `</res:senha>
				<res:tipo>L</res:tipo>
				<res:resultado>T</res:resultado>
				<res:lingua>101</res:lingua>
				<res:objetos>` + objeto + `</res:objetos>
			</res:buscaEventos>
		</x:Body>
	</x:Envelope>
	`)
	req, err := http.NewRequest("POST", WSDL, strings.NewReader(payload))
	if err != nil {
		return "", err
	}

}

func main() {

}
