package main

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

var (
	endpoint = flag.String("e", "https://webservice.correios.com.br:443/service/rastro", "endpoint do webservice sro")
	objeto   = flag.String("o", "QK478713098BR", "objeto a ser rastreado")
	usuario  = flag.String("u", "", "nome de usuario")
	senha    = flag.String("p", "", "senha do usuario")
)

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
	Body struct {
		BuscaEventosResponse struct {
			Return struct {
				Objeto struct {
					Evento struct {
						Tipo      string `xml:"tipo"`
						Status    string `xml:"status"`
						Data      string `xml:"data"`
						Hora      string `xml:"hora"`
						Descricao string `xml:"descricao"`
						Codigo    string `xml:"codigo"`
						Cidade    string `xml:"cidade"`
						UF        string `xml:"uf"`
					} `xml:"evento"`
				} `xml:"objeto"`
			} `xml:"return"`
		} `xml:"buscaEventosResponse"`
	} `xml:"Body"`
}

func rastreia(wsdl string, objeto string, usuario string, senha string) (BuscaEventosResponse, error) {
	rastro := BuscaEventosResponse{}
	payload := fmt.Sprintf(`
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:res="http://resource.webservice.correios.com.br/">
			<soapenv:Header/>
			<soapenv:Body>
				<res:buscaEventos>
					<!--Optional:-->
					<usuario>` + usuario + `</usuario>
					<!--Optional:-->
					<senha>` + senha + `</senha>
					<!--Optional:-->
					<tipo>L</tipo>
					<!--Optional:-->
					<resultado>T</resultado>
					<!--Optional:-->
					<lingua>101</lingua>
					<!--Optional:-->
					<objetos>` + objeto + `</objetos>
				</res:buscaEventos>
			</soapenv:Body>
		</soapenv:Envelope>
 	`)
	req, err := http.NewRequest("POST", wsdl, strings.NewReader(payload))
	if err != nil {
		return rastro, err
	}
	http.DefaultClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		return rastro, errors.New("erro http status: " + res.Status)
	}
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
	err = xml.Unmarshal([]byte(b), &rastro)
	if err != nil {
		return rastro, err
	}
	return rastro, nil

}
func printJson(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
func main() {
	flag.Parse()
	rastro, err := rastreia(*endpoint, *objeto, *usuario, *senha)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(printJson(rastro))
}
