#!/bin/bash

# Test Receipt OCR Endpoint
# Make sure the API is running before executing this script

API_URL="http://localhost:8080"
TOKEN="your_jwt_token_here"

echo "=== Testing Receipt OCR Endpoint ==="
echo ""

# Test 1: Upload receipt image
echo "1. Uploading receipt for OCR processing..."
curl -X POST "$API_URL/api/v1/transactions/import/receipt" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@struk.jpg" \
  -F "account_id=1"

echo ""
echo ""

# Test 2: Confirm the receipt data (replace with actual data from step 1)
echo "2. Confirming receipt data..."
curl -X POST "$API_URL/api/v1/transactions/import/receipt/confirm" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "merchant": "Indomaret Sudirman",
    "amount": -145000,
    "date": "2026-04-26",
    "account_id": 1,
    "category_id": 5,
    "notes": "belanja bulanan",
    "source_id": 3
  }'

echo ""
echo ""

echo "=== Test Complete ==="
