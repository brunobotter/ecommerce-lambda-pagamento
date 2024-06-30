package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	sqsSvc *sqs.SQS
)

func init() {
	log.Println("Criando sessão da AWS!")
	// Cria uma nova sessão AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Lê a região da variável de ambiente AWS_REGION
	})
	if err != nil {
		log.Fatalf("Erro ao criar sessão AWS: %v", err)
	}

	// Cria serviços SQS e SNS
	sqsSvc = sqs.New(sess)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Println("Iniciando processamento do evento SQS")
	for _, message := range sqsEvent.Records {
		// Decodifica o JSON do corpo da mensagem
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(message.Body), &payload); err != nil {
			log.Printf("Erro ao decodificar JSON: %v", err)
			continue
		}

		nomeCliente, ok := payload["nome_cliente"].(string)
		if !ok {
			log.Println("Campo nome_cliente não encontrado ou inválido")
			continue
		}

		/*emailCliente, ok := payload["email_cliente"].(string)
		if !ok {
			log.Println("Campo email_cliente não encontrado ou inválido")
			continue
		}*/

		nomeProduto, ok := payload["nome_produto"].(string)
		if !ok {
			log.Println("Campo nome_produto não encontrado ou inválido")
			continue
		}

		valor, ok := payload["valor"].(float64)
		if !ok {
			log.Println("Campo valor não encontrado ou inválido")
			continue
		}

		quantidade, ok := payload["quantidade"].(float64)
		if !ok {
			log.Println("Campo quantidade não encontrado ou inválido")
			continue
		}

		tipoPagamento, ok := payload["tipo_pagamento"].(string)
		if !ok {
			log.Println("Campo tipo_pagamento não encontrado ou inválido")
			continue
		}

		// Lógica para processar o tipo de pagamento
		switch tipoPagamento {
		case "cartao_credito":
			log.Println("Processando pagamento com cartão de crédito...")
			log.Println("Olá %s,\n\nSeu pagamento foi processado com sucesso.\n\nDetalhes do pagamento:\nTipo Pagamento: %s\nProduto: %s\nValor: %.2f\nQuantidade: %.2f\n\nObrigado por sua compra!",
				nomeCliente, tipoPagamento, nomeProduto, valor, quantidade)
		case "pix":
			log.Println("Processando pagamento via PIX...")
			log.Println("Olá %s,\n\nSeu pagamento foi processado com sucesso.\n\nDetalhes do pagamento:\nTipo Pagamento: %s\nProduto: %s\nValor: %.2f\nQuantidade: %.2f\n\nObrigado por sua compra!",
				nomeCliente, tipoPagamento, nomeProduto, valor, quantidade)
		case "boleto":
			log.Println("Processando pagamento via boleto...")
			log.Println("Olá %s,\n\nSeu pagamento foi processado com sucesso.\n\nDetalhes do pagamento:\nTipo Pagamento: %s\nProduto: %s\nValor: %.2f\nQuantidade: %.2f\n\nObrigado por sua compra!",
				nomeCliente, tipoPagamento, nomeProduto, valor, quantidade)
		default:
			log.Printf("Tipo de pagamento não reconhecido: %s", tipoPagamento)
			continue
		}
		// Exemplo de exclusão da mensagem do SQS após processamento
		_, err := sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(os.Getenv("https://sqs.us-east-1.amazonaws.com/730335442778/Pagamento.fifo")),
			ReceiptHandle: aws.String(message.ReceiptHandle),
		})
		if err != nil {
			log.Printf("Erro ao excluir mensagem do SQS: %v", err)
			return err
		}
		log.Println("Mensagem do SQS excluída com sucesso")
	}

	return nil
}
func main() {
	lambda.Start(handler)

}
