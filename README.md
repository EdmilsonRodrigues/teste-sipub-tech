# Movies API
Essa é uma API criada com o objetivo de cumprir o teste técnico para a empresa Sipub.Tech.

O desafio consistia na criação de uma API que cumprisse os seguintes requisitos:
- API REST (Feito)
- Usando Docker (Feito)
- Usando Go (Feito)
- Usando MongoDB (Substituído)
- Documentação da Aplicação em Swagger (Feita - Incompleta)
- Arquitetura Hexagonal (Feito)
- Microsserviços (Feito)
- Comunicação entre API e Movies via gRPC (Queries)
- Inicialização com apenas um comando: `make deploy-docker` ou `make deploy-microk8s` ou `make deploy-lxd`
- Extras:
  - Criar manifestos k8s dos serviços: (Feito)
  - Fazer as operações de Execute Event Driven (Feito)
  - Usar **LocalStack** no lugar do MongoDB (Feito)
  
  
## Como inicializar a aplicação?

Para qualquer uma das inicializações a seguir, recomendo que esteja em um ambiente GNU/Linux e necessita do GNU Make instalado.

### Para deploy com docker compose
Necessita de Docker e Docker Compose instalados na sua máquina:
```bash
make deploy-docker
```

### Para deploy com microk8s localmente
Necessita do microk8s instalado na sua máquina:
Não recomendo esse deploy pois ele irá aplicar outros manifestos na sua máquina local que podem ser indesejados. 
Para um maior isolamento, verificar o próximo deploy.
```bash
make deploy-microk8s
```

### Para deploy com microk8s uma vm gerida pelo LXD
Recomendo essa no lugar da última. Necessita do LXD instalado e configurado na máquina.
```bash
make deploy-lxd
```

## Documentação das rotas criadas

Corpo das respostas:
Filme: 
```json
{
    "data": {
        "id": 1,
        "title": "exemplo",
        "year": "1990"
    }
}
```

Filmes:
```json
{
    "data": [
        {
            "id": 4,
            "title": "exemplo 2",
            "year": "1895"
        },
        ...
    ],
    "limit": 5,
    "cursor": 3
}
```

Erros:
```json
{
    "details": {
        "message": "Isso é um erro"
    }
}
```

### GET /movies/
Permite realizar o fetch de múltiplos filmes na API.
Aceita 3 query paramenters
- year -> Um inteiro entre 1880 e o ano atual. Irá buscar somente os filmes lançados nesse ano.
- limit -> Um inteiro que limita o número de filmes buscados.
- cursor -> O id do último filme buscado pela query anterior. A próxima query pulará todos os filmes antes do 
  filme apontado pelo cursor.

### GET /movies/:id
Permite buscar um filme na API pelo id.


### POST /movies/
Recebe um JSON com o título e o ano. 
A requisição é processada em background, mas, mesmo que algo a impeça de ser processada no momento, 
ela volta para a fila até ser processada.
```json
{
    "title": "O labirinto do Fauno",
    "year": "2006"
}
```

### DELETE /movies/:id
Deleta o filme com o ID passado.
Se o filme não existir, nada acontece
A requisição é processada em background, mas, mesmo que algo a impeça de ser processada no momento, 
ela volta para a fila até ser processada.


## Exemplos de uso via curl
Para preencher automaticamente o repositório com os dados de input basta usar o comando:
```bash
make fill-db  # No caso de deploy com docker compose
make fill-db DYNAMO_DB_ENDPOINT=http:IP:NODE_PORT_LOCALSTACK   # No caso de deploy com k8s, sendo o IP localhost para cluster microk8s
                                                               # local e o IP da VM no caso de cluster isolado, e NODE_PORT_LOCALSTACK
                                                               # é a porta do serviço NodePort do localstack que transmite para a 
                                                               # porta 4566.
```

Para fazer requisições, primeiro deve-se pegar o IP e porta.
- Se o deploy escolhido foi via docker, o IP é `localhost` e a porta é `8080`.
- Se o deploy escolhido foi microk8s local, o IP é `localhost` e a porta é a NodePort do serviço de api-gateway.
- Se o deploy escolhido foi microk8s em um vm LXD, o IP é o IP padrão da VM (o primeiro que aparece), 
  e a porta é a NodePort do serviço de api-gateway.

Listar filmes:
```bash
curl http://IP:PORT/movies/                                 # Lista múltiplos filmes
curl http://IP:PORT/movies/?year=1992                       # Lista filmes de 1992
curl http://IP:PORT/movies/?limit=10&year=1940              # Lista até 10 filmes de 1940
curl http://IP:PORT/movies/?limit=10&year=1940&cursor=45    # Lista até 10 filmes de 1940, após o filme de id 45
curl http://IP:PORT/movies/?limit=15&cursor=155             # Lista até 10 filmes após o filme de ID 155

```

Pegar Filme:
```bash
curl http://IP:PORT/movies/45  # busca o filme com ID 45
```

Criar filme:
No exemplo adiciona o filme O labirinto do Fauno à API.
```bash
curl -X POST -H "Content-Type: application/json" -d '{"title": "O labirinto do Fauno", "year": "2006"}' http://ID:PORT
```

Deletar filme:
```bash
curl -X DELETE http://IP:PORT/movies/45  # deleta o filme com ID 45
```

## Espaço para melhorias:
### Documentação
A documentação Swagger da aplicação necessita de muitas melhorias, que virão logo, em próximas versões do projeto.
Lhe falta de essencial o corpo do método POST e o corpo das respostas, além de descrições mais detalhadas.
Isso contudo, será adicionado futuramente ao projeto.

### Indepotência
Até o momento, o método POST pode criar uma cópia de um filme já cadastrado. Há porém a necessidade de se adicionar mais campos, 
tendo em vista que é possível mais de um filme terem o mesmo nome e ano.


### Observabilidade
A aplicação não possui logs de forma satisfatória, e deve aumentar o número de logs para se ter melhor observablidade em produção.
Além de logs, não há integração com tracers, nem com serviço de métricas, como Prometheus. 

Em uma próxima versão será adicionado o suporte a observabilidade.

### Health Checks
Imagens Docker não possuem health checks, e tais health checks devem ser adicionados não somente nas imagens, mas também nos deploys,
incluindo esperas, para que os pods não fiquem falhando enquanto esperam suas dependências subirem.


Todas essas anotações foram observadas por mim e serão levadas em conta na próxima versão.

### Importer
O importer de dados do JSON não tem semáforos, então algumas goroutines podem acabar entrando em panico 
e não sendo possível salvar todos os filmes. 

Uma possível melhoria seria fazer um batch push de todos os movies e, caso não seja possível, adicionar um semáforo para impedir
que seja ultrapassado o rate limit da API.

