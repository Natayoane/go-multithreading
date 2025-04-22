# Desafio: Busca Rápida de Endereços com Multithreading e APIs

## Descrição do Desafio

Este projeto busca informações de endereço a partir de um CEP, utilizando duas APIs públicas (BrasilAPI e ViaCEP) de forma concorrente. O objetivo é exibir a resposta da API mais rápida, descartando a mais lenta, com os dados do endereço (rua, cidade, estado, CEP) e a origem da API, respeitando um limite de 1 segundo para timeout.

### Requisitos
- Fazer requisições simultâneas às APIs:
  - BrasilAPI: `https://brasilapi.com.br/api/cep/v1/{CEP}`
  - ViaCEP: `http://viacep.com.br/ws/{CEP}/json/`
- Usar a resposta mais rápida e descartar a mais lenta.
- Exibir os dados do endereço e a API de origem.
- Tratar timeout de 1 segundo com mensagem de erro.

## Como Executar

### Pré-requisitos
- Go 1.16 ou superior.
- Conexão com a internet.

### Passos
1. No terminal, navegue até o diretório do projeto.
2. Execute:
   ```bash
   go run main.go
   ```
3. O programa usa o CEP `89010-904` e exibe o resultado no terminal.

### Exemplo de Saída
```
Resposta de BrasilAPI:
Rua: Rua XV de Novembro
Cidade: Blumenau
Estado: SC
CEP: 89010-904
```
Ou, em caso de timeout:
```
Ocorreu timeout de 1 segundo
```

## Melhorias Implementadas Além do Desafio
- **Validação de CEP**: Verifica se o CEP tem 8 dígitos numéricos.
- **Cancelamento eficiente**: Usa `context.Context` para interromper requisições após a primeira resposta ou timeout.
- **Cliente HTTP otimizado**: Configura timeout no `http.Client` para evitar vazamento de recursos.
- **Parseamento de JSON**: Converte respostas das APIs em uma estrutura comum para exibição padronizada.
- **Tratamento de erros**: Lida com falhas de rede, HTTP e parseamento, exibindo mensagens claras.
- **Código modular**: Organiza a lógica em funções reutilizáveis para melhor manutenção.
