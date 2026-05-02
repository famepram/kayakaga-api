package helper

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MetaInfo struct {
	Timestamp string `json:"timestamp"`
	Total     *int64 `json:"total,omitempty"`
	Page      *int   `json:"page,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    &MetaInfo{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
		Meta:    &MetaInfo{Timestamp: time.Now().UTC().Format(time.RFC3339)},
	})
}

func ErrorResponse(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	})
}

func PaginatedResponse(c *gin.Context, data interface{}, total *int64, page *int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &MetaInfo{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Total:     total,
			Page:      page,
		},
	})
}

func FormatRupiah(amount int64) string {
	if amount == 0 {
		return "Rp0"
	}

	negative := amount < 0
	if negative {
		amount = -amount
	}

	str := ""
	for i := 0; i < len(string(amount)); i++ {
		pos := len(string(amount)) - i
		if pos%3 == 0 && i != 0 {
			str += "."
		}
		str += string(string(amount)[i])
	}

	if negative {
		str = "-Rp" + str
	} else {
		str = "Rp" + str
	}

	return str
}

func RoundToFloat64(val int64) float64 {
	return float64(val) / math.Pow10(2)
}
