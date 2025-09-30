package handlers

import (
        "net/http"
        "os"
        "sync"

        "github.com/gin-gonic/gin"
        "github.com/openai/openai-go/v2"
)

var (
        openaiClient     openai.Client
        openaiClientOnce sync.Once
)

func getOpenAIClient() openai.Client {
        openaiClientOnce.Do(func() {
                openaiClient = openai.NewClient()
        })
        return openaiClient
}

func Chat(c *gin.Context) {
        var input struct {
                Message string `json:"message" binding:"required"`
        }

        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        apiKey := os.Getenv("OPENAI_API_KEY")
        if apiKey == "" {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI API key not configured"})
                return
        }

        client := getOpenAIClient()

        systemPrompt := `You are an AI assistant for an HR Management System. Your role is to help answer questions about:
- Employee information and records
- Attendance tracking and policies
- Leave requests and vacation policies
- Salary and payroll information
- HR best practices and guidelines

Be professional, helpful, and concise in your responses. If asked about specific employee data that you don't have access to, 
politely explain that you're an AI assistant and suggest they check the relevant section in the HCM system.`

        chatCompletion, err := client.Chat.Completions.New(c.Request.Context(), openai.ChatCompletionNewParams{
                Messages: []openai.ChatCompletionMessageParamUnion{
                        openai.SystemMessage(systemPrompt),
                        openai.UserMessage(input.Message),
                },
                Model: openai.ChatModelGPT4oMini,
        })

        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI response", "details": err.Error()})
                return
        }

        if len(chatCompletion.Choices) == 0 {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "No response from AI"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "response": chatCompletion.Choices[0].Message.Content,
                "message":  input.Message,
        })
}
