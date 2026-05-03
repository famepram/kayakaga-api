package transaction

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"kayakaga-api/domain/mysql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ReceiptOCRResponse struct {
	Merchant   string              `json:"merchant"`
	Total      int64               `json:"total"`
	Date       string              `json:"date"`
	Items      []ReceiptItem       `json:"items"`
	Confidence float64             `json:"confidence"`
}

type ReceiptItem struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Qty   int    `json:"qty"`
}

type ReceiptData struct {
	Merchant      string        `json:"merchant"`
	Amount        int64         `json:"amount"`
	Date          string        `json:"date"`
	CategoryID    uint          `json:"category_id"`
	CategoryName  string        `json:"category_name"`
	AccountID     uint          `json:"account_id"`
	AiCategorized bool          `json:"ai_categorized"`
	Confidence    float64       `json:"confidence"`
	Items         []ReceiptItem `json:"items"`
	Warning       string        `json:"warning"`
}

type OpenRouterRequest struct {
	Model     string                `json:"model"`
	MaxTokens int                   `json:"max_tokens"`
	Messages  []OpenRouterMessage   `json:"messages"`
}

type OpenRouterMessage struct {
	Role    string                `json:"role"`
	Content []OpenRouterContent   `json:"content"`
}

type OpenRouterContent struct {
	Type      string `json:"type"`
	ImageURL  *OpenRouterImageURL `json:"image_url,omitempty"`
	Text      string `json:"text,omitempty"`
}

type OpenRouterImageURL struct {
	URL string `json:"url"`
}

type OpenRouterResponse struct {
	Choices []OpenRouterChoice `json:"choices"`
	Error   *OpenRouterError   `json:"error,omitempty"`
}

type OpenRouterChoice struct {
	Message OpenRouterMessageResponse `json:"message"`
}

type OpenRouterMessageResponse struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type OpenRouterError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    interface{} `json:"code"`
}

type ConfirmReceiptRequest struct {
	Merchant   string `json:"merchant" binding:"required"`
	Amount     int64  `json:"amount" binding:"required"`
	Date       string `json:"date" binding:"required"`
	AccountID  uint   `json:"account_id" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
	Notes      string `json:"notes"`
	SourceID   uint   `json:"source_id"`
}

func ProcessReceipt(imageData []byte, accountID uint) (*ReceiptData, error) {
	log.Println("🔍 [OCR] Starting receipt processing...")
	log.Printf("📸 [OCR] Image size: %d bytes", len(imageData))

	mimeType, err := detectImageType(imageData)
	if err != nil {
		log.Printf("❌ [OCR] Image type detection failed: %v", err)
		return nil, errors.New("invalid image format")
	}
	log.Printf("✅ [OCR] Image type detected: %s", mimeType)

	base64Image := base64.StdEncoding.EncodeToString(imageData)
	log.Printf("📦 [OCR] Base64 encoded size: %d chars", len(base64Image))

	log.Println("🤖 [OCR] Calling OpenRouter API...")
	ocrResult, err := callOpenRouterOCR(mimeType, base64Image)
	if err != nil {
		log.Printf("❌ [OCR] OpenRouter API call failed: %v", err)
		return nil, fmt.Errorf("AI service error: %w", err)
	}
	log.Println("✅ [OCR] OpenRouter API call successful")

	log.Printf("📊 [OCR] Raw OCR result:")
	log.Printf("   - Merchant: %s", ocrResult.Merchant)
	log.Printf("   - Total: %d", ocrResult.Total)
	log.Printf("   - Date: %s", ocrResult.Date)
	log.Printf("   - Confidence: %.2f", ocrResult.Confidence)
	log.Printf("   - Items: %d", len(ocrResult.Items))

	if ocrResult.Total == 0 {
		log.Println("❌ [OCR] Total amount is zero, cannot process")
		return nil, errors.New("could not extract total amount from receipt")
	}

	merchant := ocrResult.Merchant
	if merchant == "" {
		merchant = "Unknown Merchant"
		log.Println("⚠️  [OCR] Merchant not detected, using default")
	}

	date := ocrResult.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
		log.Println("⚠️  [OCR] Date not detected, using today")
	}

	amount := ocrResult.Total
	if amount > 0 {
		amount = -amount
		log.Printf("💰 [OCR] Converted positive amount to negative: %d", amount)
	}

	categoryID := categorizeTransaction(merchant, amount)
	categoryName := getCategoryName(categoryID)
	log.Printf("🏷️  [OCR] Auto-categorized as: %s (ID: %d)", categoryName, categoryID)

	warning := ""
	if ocrResult.Confidence < 0.5 {
		warning = "Low confidence scan. Please verify the details before saving."
		log.Printf("⚠️  [OCR] Low confidence detected: %.2f", ocrResult.Confidence)
	}

	log.Println("✅ [OCR] Receipt processing completed successfully")
	return &ReceiptData{
		Merchant:      merchant,
		Amount:        amount,
		Date:          date,
		CategoryID:    categoryID,
		CategoryName:  categoryName,
		AccountID:     accountID,
		AiCategorized: true,
		Confidence:    ocrResult.Confidence,
		Items:         ocrResult.Items,
		Warning:       warning,
	}, nil
}

func detectImageType(data []byte) (string, error) {
	log.Println("🔍 [IMAGE] Detecting image type...")
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		log.Printf("❌ [IMAGE] Failed to decode image config: %v", err)
		return "", err
	}
	log.Printf("✅ [IMAGE] Format detected: %s", format)

	switch format {
	case "jpeg":
		return "image/jpeg", nil
	case "png":
		return "image/png", nil
	default:
		log.Printf("❌ [IMAGE] Unsupported format: %s", format)
		return "", errors.New("unsupported image format")
	}
}

func callOpenRouterOCR(mimeType, base64Image string) (*ReceiptOCRResponse, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	model := os.Getenv("OPENROUTER_VISION_MODEL")

	log.Println("🌐 [API] Initializing OpenRouter API call...")
	log.Printf("🔑 [API] Base URL: %s", baseURL)
	log.Printf("🤖 [API] Model: %s", model)

	if apiKey == "" {
		log.Println("❌ [API] API key not configured")
		return nil, errors.New("OpenRouter API key not configured")
	}
	log.Println("✅ [API] API key found")

	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	if model == "" {
		model = "google/gemini-2.0-flash-001"
	}

	payload := OpenRouterRequest{
		Model:     model,
		MaxTokens: 500,
		Messages: []OpenRouterMessage{
			{
				Role: "user",
				Content: []OpenRouterContent{
					{
						Type: "image_url",
						ImageURL: &OpenRouterImageURL{
							URL: fmt.Sprintf("data:%s;base64,%s", mimeType, base64Image),
						},
					},
					{
						Type: "text",
						Text: "Extract transaction data from this receipt. Return ONLY valid JSON, no explanation:\n{\"merchant\":\"store name\",\"total\":0,\"date\":\"YYYY-MM-DD\",\"items\":[{\"name\":\"item\",\"price\":0,\"qty\":1}],\"confidence\":0.0}\nRules:\n- total must be integer (IDR, no decimals)\n- date format: YYYY-MM-DD\n- confidence: 0.0-1.0\n- Use null for fields you cannot read clearly",
					},
				},
			},
		},
	}

	log.Println("📝 [API] Sending request to OpenRouter...")

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ [API] Failed to marshal payload: %v", err)
		return nil, err
	}
	log.Printf("📦 [API] Request payload size: %d bytes", len(jsonData))

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ [API] Failed to create request: %v", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://finai.app")
	req.Header.Set("X-Title", "Finai")

	client := &http.Client{Timeout: 30 * time.Second}
	log.Println("⏳ [API] Waiting for response...")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ [API] Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("📡 [API] Response status: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [API] Failed to read response body: %v", err)
		return nil, err
	}
	log.Printf("📦 [API] Response body size: %d bytes", len(body))

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		log.Printf("❌ [API] Failed to unmarshal response: %v", err)
		log.Printf("📄 [API] Raw response: %s", string(body))
		return nil, err
	}

	if openRouterResp.Error != nil {
		log.Printf("❌ [API] OpenRouter returned error: %s", openRouterResp.Error.Message)
		return nil, fmt.Errorf("OpenRouter API error: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		log.Println("❌ [API] No choices in response")
		return nil, errors.New("no response from AI model")
	}

	log.Printf("✅ [API] Received %d choice(s)", len(openRouterResp.Choices))

	// Extract text content from AI response
	var aiResponse string
	content := openRouterResp.Choices[0].Message.Content

	log.Printf("🔍 [API] Content type: %T", content)

	switch v := content.(type) {
	case string:
		aiResponse = strings.TrimSpace(v)
		log.Printf("📄 [API] Content is string type, length: %d", len(v))
	case []interface{}:
		log.Printf("📄 [API] Content is array type with %d items", len(v))
		for _, item := range v {
			if contentMap, ok := item.(map[string]interface{}); ok {
				if contentType, ok := contentMap["type"].(string); ok && contentType == "text" {
					if text, ok := contentMap["text"].(string); ok {
						aiResponse = strings.TrimSpace(text)
						log.Printf("📄 [API] Found text content, length: %d", len(text))
						break
					}
				}
			}
		}
	default:
		log.Printf("⚠️  [API] Unknown content type: %T", content)
	}

	if aiResponse == "" {
		log.Println("❌ [API] Empty AI response")
		return nil, errors.New("AI returned empty response")
	}

	log.Printf("📄 [API] AI Response: %s", aiResponse)

	startIdx := strings.Index(aiResponse, "{")
	endIdx := strings.LastIndex(aiResponse, "}")

	if startIdx == -1 || endIdx == -1 {
		log.Println("❌ [API] No JSON found in response")
		return nil, errors.New("AI did not return valid JSON")
	}

	jsonStr := aiResponse[startIdx : endIdx+1]
	log.Printf("📄 [API] Extracted JSON: %s", jsonStr)

	var ocrResult ReceiptOCRResponse
	if err := json.Unmarshal([]byte(jsonStr), &ocrResult); err != nil {
		log.Printf("❌ [API] Failed to parse JSON: %v", err)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	log.Println("✅ [API] OCR result parsed successfully")
	return &ocrResult, nil
}

func ConfirmReceipt(userID uint, req *ConfirmReceiptRequest) (*mysql.Transaction, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	txn := &mysql.Transaction{
		UserID:        userID,
		AccountID:     req.AccountID,
		CategoryID:    req.CategoryID,
		SourceID:      req.SourceID,
		Date:          date,
		Merchant:      req.Merchant,
		Amount:        req.Amount,
		Notes:         req.Notes,
		IsRecurring:   0,
		AiCategorized: 1,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	return txn, nil
}

func getCategoryName(categoryID uint) string {
	names := map[uint]string{
		1:  "Makanan & Minuman",
		2:  "Transportasi",
		3:  "Hiburan",
		4:  "Tagihan",
		5:  "Belanja",
		6:  "Kesehatan",
		7:  "Investasi",
		8:  "Lainnya",
		9:  "Pemasukan",
	}

	if name, ok := names[categoryID]; ok {
		return name
	}
	return "Lainnya"
}
