package claude

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sirupsen/logrus"
)

type IClaude interface {
	CreateChat(req dto.ChatRequest) (dto.ChatResponse, error)
}

type claude struct {
	log *logrus.Logger
}

type content struct {
	Type        string `json:"type,omitempty"`
	Text        string `json:"text,omitempty"`
	Translation string `json:"translation,omitempty"`
}

type typeResponse struct {
	SourceLanguage string `json:"source_language"`
	Translation    string `json:"translation"`
	Explanation    string `json:"explanation"`
}

type ClaudeResponse struct {
	ID      string    `json:"id"`
	Content []content `json:"content"`
}

func NewClaude(log *logrus.Logger) IClaude {
	return &claude{
		log: log,
	}
}

func (c *claude) CreateChat(req dto.ChatRequest) (dto.ChatResponse, error) {
	client := anthropic.NewClient(os.Getenv("CLAUDE_API_KEY"))
	var modelPrompt string

	if req.Type != "image" {
		modelPrompt = "Translate the following Indonesian regional dialect or language to standard Bahasa Indonesia. Provide a response with: 'translation': the standard Indonesian translation. Format: {\"translation\": \"Saya pusing\"} Input text to translate: [INPUT_TEXT]"
	} else {
		modelPrompt = "You are an assistant that translates Indonesian regional languages and dialects into standard Bahasa Indonesia. Analyze input text and provide a JSON response with: 1) 'source_language': the detected dialect/language, 2) 'translation': the standard Indonesian translation, and 3) 'explanation': cultural context and usage explanation of key phrases. Format example: {\"source_language\": \"Sundanese\", \"translation\": \"Menurut saya, saya tidak bisa datang.\", \"explanation\": \"'Saur' berasal dari tradisi Sunda dalam berdialog sopan dengan orang yang lebih tua, menunjukkan rasa hormat. Contoh: 'Saur abdi, teu tiasa sumping.'\"} If language is unidentifiable, mark as 'Unidentified'. I understand various Indonesian languages including Javanese, Sundanese, Balinese, Madurese, Minangkabau, Batak, Bugis, Acehnese, Betawi, and regional dialects. Provide only the JSON response, no additional text."
	}

	userPrompt := req.Content

	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Dot5Sonnet20241022,
		Messages: []anthropic.Message{
			{
				Role:    anthropic.RoleAssistant,
				Content: []anthropic.MessageContent{{Type: "text", Text: &modelPrompt}},
			},
			{
				Role:    anthropic.RoleUser,
				Content: []anthropic.MessageContent{{Type: "text", Text: &userPrompt}},
			},
		},
		MaxTokens: 512,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			c.log.WithFields(map[string]interface{}{
				"error": e.Error(),
			}).Error("Claude API error")
		} else {
			c.log.WithFields(map[string]interface{}{
				"error": err.Error(),
			}).Error("Claude API error")
		}
		return dto.ChatResponse{}, &response.ErrChatFailed
	}

	c.log.WithFields(map[string]interface{}{
		"resp": resp.Content[0].GetText(),
	}).Info("Claude API response received")

	var typeResp typeResponse
	err = json.Unmarshal([]byte(resp.Content[0].GetText()), &typeResp)
	if err != nil {
		c.log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to unmarshal Claude response")
		return dto.ChatResponse{}, &response.ErrChatFailed
	}

	userResp := dto.ChatResponse{
		ID:             resp.ID,
		UserID:         req.UserID,
		Translation: 	typeResp.Translation,
		SourceLanguage: typeResp.SourceLanguage,
		Explanation:    typeResp.Explanation,
	}

	return userResp, nil
}
