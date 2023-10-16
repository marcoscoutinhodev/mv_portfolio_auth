# [ms_auth] - Micro Serviço de Autênticação

### Descrição
Este serviço é responsável por gerenciar a autenticação do usuário, incluindo operações como registro, login, confirmação de e-mail, recuperação de senha e atualização de tokens de autênticação.

### 🛠 Tecnologias
As seguintes ferramentas foram usadas na construção do projeto:
- Golang
- GO-Chi
- SQLC
- Golang-migrate
- Google Wire
- Postgres
- RabbitMQ
- K8S (Apenas criei os arquivos de configuração para simular um possível deploy. O micro serviço foi testado localmente usando o Minikube. Isso assegurou que todas as configurações estivessem corretas e que a aplicação estivesse pronta para ser implantada em um ambiente de produção, embora as configurações de secret que são informações sensíveis teriam que ser criadas no momento de criação do cluster, pois para essas informações não são feito commits)

### Endpoints
#### Cadastrar Usuário:

    Método: POST
    URL: /auth/signup
    Headers:
        Content-Type: application/json
    
Exemplo de corpo da requisição:
```json
{
  "name": "LOREM IPSUM",
  "email": "lorem@ipsum.com",
  "password": "P@ssw0rd"
}
```

#### Autenticar Usuário:

    Método: POST
    URL: /auth/signin
    Headers:
        Content-Type: application/json

#### Exemplo de corpo da requisição:

```json
{
  "email": "lorem@ipsum.com",
  "password": "P@ssw0rd"
}
```

#### Solicitação de Confirmação de E-mail:

    Método: POST
    URL: /auth/email-confirmation-request
    Headers:
        Content-Type: application/json
        email: lorem@ipsum.com

#### Confirmação de E-mail:

    Método: POST
    URL: /auth/confirm-email
    Headers:
        Content-Type: application/json
        x_access_token: {{Bearer Token}}
        
#### Recuperação de Senha:

    Método: POST
    URL: /auth/forgot-password
    Headers:
        Content-Type: application/json

#### Exemplo de corpo da requisição:
```json
{
  "email": "lorem@ipsum.com"
}
```

#### Atualização de Senha:

    Método: POST
    URL: /auth/update-password
    Headers:
        Content-Type: application/json
        x_access_token: {{Bearer Token}}

#### Exemplo de corpo da requisição:
```json
{
   "password": "P@ssw0rd"
}
```

#### Atualização de Token com Refresh Token:

    Método: POST
    URL: /auth/refresh-token
    Headers:
        Content-Type: application/json
        x_access_token: {{Bearer Token}}


### Configuração e Uso
##### [PRECISA DO DOCKER INSTALADO]

- Clone o repositório
- Criar o arquivo .env e configurar as variáveis de ambientes que estão no arquivo .env.example
- Execute o script start_dev.sh para iniciar a aplicação em modo desenvolvimento
- Execute o serviço usando go run cmd/http/main.go
- Se deseja apenas executar sem realizar alterações, utilize o Dockerfile.prd para gerar a imagem com o executavel da aplicação, você vai precisar dos containers do Postgres e do RabbitMQ, pode utilizar o docker-compose.dev.yaml e apenas alterar a imagem do APP para a imagem gerada
