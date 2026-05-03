# Receipt OCR Feature

## Overview
The Receipt OCR feature allows users to upload receipt images for automatic data extraction using AI vision models.

## Endpoints

### 1. Process Receipt (OCR)
**POST** `/api/v1/transactions/import/receipt`

Upload a receipt image for AI-powered text extraction.

**Request:**
- Content-Type: `multipart/form-data`
- Authorization: `Bearer {token}`

**Parameters:**
- `file`: Image file (JPG, JPEG, PNG, WEBP, max 10MB)
- `account_id`: Target account ID for the transaction

**Response:**
```json
{
  "success": true,
  "data": {
    "merchant": "Indomaret Sudirman",
    "amount": -145000,
    "date": "2026-04-26",
    "category_id": 5,
    "category_name": "Belanja",
    "account_id": 1,
    "ai_categorized": true,
    "confidence": 0.92,
    "items": [
      {"name": "Indomie Goreng", "price": 7000, "qty": 2},
      {"name": "Aqua 600ml", "price": 5000, "qty": 1}
    ],
    "warning": null
  }
}
```

**Warning Response (Low Confidence):**
```json
{
  "success": true,
  "data": {
    "warning": "Low confidence scan. Please verify the details before saving.",
    "confidence": 0.42,
    ...
  }
}
```

### 2. Confirm Receipt
**POST** `/api/v1/transactions/import/receipt/confirm`

Save the extracted receipt data as a transaction.

**Request:**
```json
{
  "merchant": "Indomaret Sudirman",
  "amount": -145000,
  "date": "2026-04-26",
  "account_id": 1,
  "category_id": 5,
  "notes": "belanja bulanan",
  "source_id": 3
}
```

**Response:**
Standard transaction creation response.

## AI Configuration

The feature uses OpenRouter API with Google Gemini 2.0 Flash vision model.

**Environment Variables:**
```env
OPENROUTER_API_KEY=sk-or-v1-...
OPENROUTER_BASE_URL=https://openrouter.ai/api/v1
OPENROUTER_VISION_MODEL=google/gemini-2.0-flash-001
```

## Auto-Categorization

The system automatically categorizes transactions based on merchant keywords:

| Keywords | Category ID | Category Name |
|----------|-------------|---------------|
| gaji, salary, transfer masuk | 9 | Pemasukan |
| grab, gojek, ojek, taxi, parkir | 2 | Transportasi |
| netflix, spotify, youtube, steam, game | 3 | Hiburan |
| pln, listrik, pdam, air, internet, wifi, bpjs | 4 | Tagihan |
| indomaret, alfamart, supermarket, mall, tokopedia, shopee | 5 | Belanja |
| apotek, klinik, dokter, rumah sakit, obat | 6 | Kesehatan |
| investasi, saham, reksa, tabungan | 7 | Investasi |
| gopay, ovo, dana | 8 | Lainnya |

## Error Handling

### OCR Failed
```json
{
  "success": false,
  "error": {
    "code": "OCR_FAILED",
    "message": "Could not extract total amount from receipt. Please enter manually."
  }
}
```

### AI Service Error
```json
{
  "success": false,
  "error": {
    "code": "AI_SERVICE_ERROR",
    "message": "Receipt scanning service unavailable. Please try again or enter manually."
  }
}
```

### Invalid Image
```json
{
  "success": false,
  "error": {
    "code": "INVALID_IMAGE",
    "message": "Invalid image format or corrupted file"
  }
}
```

## Testing

### Using cURL:
```bash
# Upload receipt for OCR
curl -X POST http://localhost:8080/api/v1/transactions/import/receipt \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@struk.jpg" \
  -F "account_id=1"

# Confirm and save
curl -X POST http://localhost:8080/api/v1/transactions/import/receipt/confirm \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "merchant": "Indomaret Sudirman",
    "amount": -145000,
    "date": "2026-04-26",
    "account_id": 1,
    "category_id": 5,
    "source_id": 3
  }'
```

### Using Test Script:
```bash
# Update TOKEN in test_receipt_ocr.sh
chmod +x test_receipt_ocr.sh
./test_receipt_ocr.sh
```

## Features

✅ **Multi-format Support**: JPG, JPEG, PNG, WEBP
✅ **AI-Powered OCR**: Uses Google Gemini 2.0 Flash for accurate text extraction
✅ **Auto-Categorization**: Intelligent keyword-based categorization
✅ **Confidence Scoring**: Returns confidence level for extracted data
✅ **Low Confidence Warnings**: Alerts when OCR quality is poor
✅ **Item Extraction**: Extracts individual items with prices and quantities
✅ **User Confirmation**: Requires user confirmation before saving
✅ **Error Recovery**: Graceful handling of AI service failures

## Implementation Notes

- Maximum file size: 10MB
- AI model timeout: 30 seconds
- Confidence threshold: 0.5 (below this shows warning)
- Default source for receipt transactions: 3 (receipt)
- All amounts are converted to negative for expenses
- Unknown merchants default to "Unknown Merchant"
- Missing dates default to today
