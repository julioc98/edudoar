package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Estrutura para a mensagem enviada à API
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Estrutura para a resposta recebida da API
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var chats = make(map[string][]ChatMessage)

func main() {
	http.HandleFunc("/", handleChatRequest)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

// Função para lidar com a solicitação de chat
func handleChatRequest(w http.ResponseWriter, r *http.Request) {

	// Verifica se é uma requisição OPTIONS (preflight request)
	if r.Method == http.MethodOptions {
		// Define os cabeçalhos CORS para a requisição OPTIONS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, User")
		return
	}

	// Define os cabeçalhos CORS para todas as outras requisições
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Obtém o cabeçalho "User" da solicitação
	user := r.Header.Get("User")
	if user == "" {
		http.Error(w, "Cabeçalho 'User' não encontrado", http.StatusBadRequest)
		return
	}

	// Verifica o método da solicitação
	if r.Method != http.MethodPost {
		http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
		return
	}

	// Lê o corpo da solicitação
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da solicitação", http.StatusInternalServerError)
		return
	}

	// Decodifica o corpo em uma estrutura
	var requestData struct {
		Message string `json:"message"`
	}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da solicitação", http.StatusBadRequest)
		return
	}

	// Chama a API do OpenAI
	response, err := callOpenAIChatAPI(user, requestData.Message)
	if err != nil {
		http.Error(w, "Erro na solicitação da API", http.StatusInternalServerError)
		return
	}

	// Monta a resposta
	responseData := struct {
		Reply string `json:"reply"`
	}{
		Reply: *response,
	}

	// Codifica a resposta em JSON
	responseBody, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Erro ao codificar a resposta", http.StatusInternalServerError)
		return
	}

	// Define os cabeçalhos da resposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// saveChat salva a conversa em memoria
func saveChat(user, role, msg string) {
	chats[user] = append(chats[user], ChatMessage{Role: role, Content: msg})
}

// Função para fazer a solicitação à API do OpenAI
func callOpenAIChatAPI(user, message string) (*string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	apiURL := "https://api.openai.com/v1/chat/completions"

	chatTmp := []ChatMessage{
		{
			Role:    "system",
			Content: "Contexto: " + os.Getenv("CONTEXT"),
		},
	}

	// Salva a conversa em memoria
	saveChat(user, "user", message)

	chatRec := append(chatTmp, chats[user]...)

	// Monta a estrutura da mensagem
	requestData := struct {
		Model    string        `json:"model"`
		Messages []ChatMessage `json:"messages"`
		User     string        `json:"user"`
	}{
		Model:    "gpt-3.5-turbo",
		Messages: chatRec,
		User:     user,
	}

	// Codifica a estrutura em JSON
	payload, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	// Cria a requisição HTTP
	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	// Define os cabeçalhos da requisição
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	// Envia a requisição à API
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Lê a resposta da API
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Decodifica a resposta em uma estrutura
	var chatResponse ChatResponse
	err = json.Unmarshal(responseData, &chatResponse)
	if err != nil {
		return nil, err
	}

	if len(chatResponse.Choices) == 0 {
		return nil, errors.New("respone is empty")
	}

	// Salva a conversa em memoria
	saveChat(user, "assistant", chatResponse.Choices[0].Message.Content)

	return &chatResponse.Choices[0].Message.Content, nil
}
