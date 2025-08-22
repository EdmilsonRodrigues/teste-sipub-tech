# Movies API
Essa é uma API criada com o objetivo de cumprir o teste técnico para a empresa Sipub.Tech.

O desafio consistia na criação de uma API que cumprisse os seguintes requisitos:
- API REST (Feito)
- Usando Docker (Feito)
- Usando Go (Feito)
- Usando MongoDB (Substituído)
- Documentação da Aplicação em Swagger ()
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
        "id": int,
        "title": string,
        "year": string,
    }
}
```

Filmes:
```json
{
    "data": [
        {
            "id": int,
            "title": string,
            "year": string,
        },
        ...
    ],
    "limit": int,
    "cursor": int,
}
```

Erros:
```json
{
    "details": {
        "message": string,
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
    "title": string,
    "year": string,
}
```

### DELETE /movies/:id
Deleta o filme com o ID passado.
Se o filme não existir, nada acontece
A requisição é processada em background, mas, mesmo que algo a impeça de ser processada no momento, 
ela volta para a fila até ser processada.


## Exemplos de uso via curl
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


