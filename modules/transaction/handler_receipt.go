package transaction

import (
	"fmt"
	"kayakaga-api/utils/helper"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *handler) ProcessReceipt(imageData []byte, accountID uint) (*ReceiptData, error) {
	return ProcessReceipt(imageData, accountID)
}

func (h *handler) ConfirmReceipt(userID uint, req *ConfirmReceiptRequest) (*TransactionResponse, error) {
	txn, err := ConfirmReceipt(userID, req)
	if err != nil {
		return nil, err
	}

	if err := h.repo.CreateTransaction(txn); err != nil {
		return nil, err
	}

	return h.buildResponse(txn), nil
}

// ImportReceiptHandler godoc
// @Summary Import transaction from receipt image
// @Description Extract transaction data from receipt image using AI OCR. Returns extracted data for user confirmation before saving.
// @Tags Transactions
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "Receipt image (JPG, JPEG, PNG, WEBP, max 10MB)"
// @Param account_id formData int true "Target Account ID"
// @Success 200 {object} helper.Response{data=ReceiptData}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 503 {object} helper.Response
// @Router /transactions/import/receipt [post]
func ImportReceiptHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountIDStr := c.PostForm("account_id")
		if accountIDStr == "" {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", "account_id is required")
			return
		}
		accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
		if err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", "invalid account_id")
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			helper.ErrorResponse(c, 400, "INVALID_FILE", "file is required")
			return
		}

		if fileHeader.Size > 10*1024*1024 {
			helper.ErrorResponse(c, 400, "INVALID_FILE", "File must be an image with max size 10MB")
			return
		}

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			helper.ErrorResponse(c, 400, "INVALID_FILE", "File must be JPG, JPEG, PNG, or WEBP")
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			helper.ErrorResponse(c, 500, "FILE_READ_ERROR", "Failed to read file")
			return
		}
		defer file.Close()

		imageData := make([]byte, fileHeader.Size)
		_, err = file.Read(imageData)
		if err != nil {
			helper.ErrorResponse(c, 500, "FILE_READ_ERROR", "Failed to read file content")
			return
		}

		if len(imageData) == 0 {
			helper.ErrorResponse(c, 400, "INVALID_FILE", "File is empty")
			return
		}

		result, err := uc.ProcessReceipt(imageData, uint(accountID))
		if err != nil {
			if strings.Contains(err.Error(), "AI service error") {
				helper.ErrorResponse(c, 503, "AI_SERVICE_ERROR", "Receipt scanning service unavailable. Please try again or enter manually.")
				return
			}
			if strings.Contains(err.Error(), "could not extract total amount") {
				helper.ErrorResponse(c, 400, "OCR_FAILED", err.Error())
				return
			}
			if strings.Contains(err.Error(), "invalid image format") {
				helper.ErrorResponse(c, 400, "INVALID_IMAGE", "Invalid image format or corrupted file")
				return
			}
			helper.ErrorResponse(c, 500, "PROCESSING_FAILED", fmt.Sprintf("Failed to process receipt: %s", err.Error()))
			return
		}

		helper.SuccessResponse(c, result)
	}
}

// ConfirmReceiptHandler godoc
// @Summary Confirm and save receipt transaction
// @Description Save the extracted receipt data as a transaction after user confirmation
// @Tags Transactions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ConfirmReceiptRequest true "Receipt transaction details"
// @Success 201 {object} helper.Response{data=TransactionResponse}
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /transactions/import/receipt/confirm [post]
func ConfirmReceiptHandler(uc UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var req ConfirmReceiptRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helper.ErrorResponse(c, 400, "INVALID_REQUEST", err.Error())
			return
		}

		if req.SourceID == 0 {
			req.SourceID = 3
		}

		resp, err := uc.ConfirmReceipt(userID, &req)
		if err != nil {
			if strings.Contains(err.Error(), "invalid date format") {
				helper.ErrorResponse(c, 400, "INVALID_DATE", err.Error())
				return
			}
			helper.ErrorResponse(c, 500, "CREATE_FAILED", err.Error())
			return
		}

		helper.CreatedResponse(c, resp)
	}
}
