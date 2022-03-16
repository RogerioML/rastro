# rastro
Pequeno client do serviço de rastreio dos Correios

# Requisitos Necessários
- Runtime Go instalado
- Usuário e senha do webservice de rastreamento dos Correios

# Para executar

go run . -u=USUARIO -p=SENHA -o=CODIGO_RASTREIO

` Por default é utilizado o endpoint seguro: https://webservice.correios.com.br:443/service/rastro que pode ser alterado para acessos sem https, modificando o endpoint de acesso, da seguinte forma`

go run . -u=USUARIO -p=SENHA -o=CODIGO_RASTREIO -e=http://webservice.correios.com.br:80/service/rastro





