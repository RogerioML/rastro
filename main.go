package main

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

const WSDL = "https://webservice.correios.com.br:443/service/rastro/Rastro.wsdl"

//IsoUtf8 converte de ISO para UTF-8
func IsoUtf8(b []byte) ([]byte, error) {
	r := charmap.ISO8859_1.NewDecoder().Reader(strings.NewReader(string(b)))
	return ioutil.ReadAll(r)
}

//estrutura para tratar requests que tiverem erro
type fault struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName xml.Name
		Fault   struct {
			FaultCode   string `xml:"faultcode"`
			FaultString string `xml:"faultstring"`
		} `xml:"Fault"`
	} `xml:"Body"`
}

type BuscaEventosResponse struct {
	XML  xml.Name `xml:"Envelope"`
	Body struct {
		BuscaEventosResponse struct {
			Evento struct {
				Tipo      string `xml:"tipo"`
				Status    string `xml:"string"`
				Data      string `xml:"data"`
				Hora      string `xml:"hora"`
				Descricao string `xml:"descricao"`
				Local     string `xml:"local"`
				Codigo    string `xml:"codigo"`
				Cidade    string `xml:"cidade"`
				UF        string `xml:"uf"`
			} `xml:"evento"`
		} `xml:"buscaEventosResponse"`
	} `xml:"Body"`
}

func rastreia(wsdl string, objeto string, usuario string, senha string) (BuscaEventosResponse, error) {
	rastro := BuscaEventosResponse{}
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
		return rastro, err
	}
	http.DefaultClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return rastro, err
	}
	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return rastro, err
	}
	b, err = IsoUtf8(b)
	if err != nil {
		return rastro, err
	}
	if strings.Contains(string(b), "faultstring") {
		respError := fault{}
		_ = xml.Unmarshal([]byte(b), &respError)
		return rastro, errors.New(respError.Body.Fault.FaultString)
	}
	_ = xml.Unmarshal([]byte(b), &rastro)
	return rastro, nil

}

func main() {
	rastro, err := rastreia(WSDL, "QK478713098BR", "ECT", "SRO")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", rastro)
}
