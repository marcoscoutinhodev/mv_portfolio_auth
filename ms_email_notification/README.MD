# [ms_email_notification] - Serviço de Envio de E-mails

### Descrição
Este serviço é responsável por consumir a fila RabbitMQ e enviar e-mails utilizando o protocolo SMTP. Ele suporta failover automático para múltiplos servidores SMTP, garantindo a entrega confiável de e-mails.

### 🛠 Tecnologias
As seguintes ferramentas foram usadas na construção do projeto:
- Golang
- Go Routines
- RabbitMQ
- SMTP (Gmail)


### Funcionalidades
    Consumir mensagens da fila RabbitMQ
    Enviar e-mails usando o SMTP do Gmail
    Suporte a failover automático entre servidores SMTP

### Configuração e Uso
##### [PRECISA DO DOCKER INSTALADO]

- Clone o repositório
- Criar o arquivo .env e configurar as variáveis de ambientes que estão no arquivo .env.example
- Execute o script start_dev.sh para iniciar a aplicação em modo desenvolvimento
- Execute o serviço usando go main.go
- Se deseja apenas executar sem realizar alterações, utilize o Dockerfile para gerar a imagem com o executavel da aplicação, você vai precisar dos mesmos containers Postgres e RabbitMQ que está utilizando no micro serviço ms_auth
