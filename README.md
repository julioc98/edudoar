# EduDoar
Educando para doar!

# Assistente Virtual para Doações a Escolas Públicas

Este é um projeto de API escrito em Go que implementa um assistente virtual para ajudar as pessoas a doarem dinheiro para escolas públicas e, assim, abater parte das doações em seu imposto de renda. O assistente utiliza a API do OpenAI GPT-3.5 Turbo para interagir com os usuários.

## Pré-requisitos

Antes de executar o projeto, certifique-se de ter o seguinte instalado em seu ambiente de desenvolvimento:

- Go (versão 1.16 ou superior)
- Chave de API do OpenAI

## Instalação

Siga as etapas abaixo para configurar e executar o projeto:

1. Faça o clone deste repositório em sua máquina local.
2. No diretório raiz do projeto, execute o comando `go mod download` para baixar as dependências necessárias.

## Configuração

Antes de executar o projeto, você precisa configurar a chave de API do OpenAI. Para fazer isso, siga as etapas abaixo:

1. Acesse o site da OpenAI (https://www.openai.com/) e crie uma conta (se ainda não tiver uma).
2. Gere uma chave de API válida para usar a API do GPT-3.5 Turbo.
3. Defina a variável de ambiente `OPENAI_API_KEY` com a sua chave de API. Por exemplo, você pode executar o seguinte comando no terminal:

   ```bash
   export OPENAI_API_KEY=SUA_CHAVE_DE_API
   ```

   Certifique-se de substituir `SUA_CHAVE_DE_API` pela sua chave de API real.

## Executando o Projeto

No diretório raiz do projeto, execute o seguinte comando para iniciar o servidor:

```bash
go run main.go
```

O servidor estará ouvindo em `http://localhost:8000`.

## Utilizando a API

A API oferece um único endpoint `POST /` para receber as solicitações do assistente virtual. A estrutura da solicitação deve ser um objeto JSON com o seguinte formato:

```json
{
  "message": "MENSAGEM_DO_USUÁRIO"
}
```

Onde:
- `"message"` é a mensagem enviada pelo usuário para o assistente virtual.

A API responderá com um objeto JSON contendo a resposta do assistente, no seguinte formato:

```json
{
  "reply": "RESPOSTA_DO_ASSISTENTE"
}
```

## Exemplo de Uso

Aqui está um exemplo de como você pode chamar a API usando a ferramenta `curl`:

```bash
curl -X POST -H "Content-Type: application/json" -H "User: NOME_DO_USUÁRIO" -d '{"message": "Olá, como posso ajudar?"}' http://localhost:8000
```

Lembre-se de substituir `NOME_DO_USUÁRIO` pelo nome do usuário atual. A resposta do assistente será exibida no terminal.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para enviar pull requests para melhorias, correções de bugs ou novos recursos.

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).

