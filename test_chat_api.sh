#!/bin/bash

# Test Chat Module Endpoints
# Make sure the API is running before executing this script

API_URL="http://localhost:8080"
TOKEN="paste_jwt_token_here"

echo "=== Testing Chat Module Endpoints ==="
echo ""

# Test 1: Health check
echo "1. Health check..."
curl "$API_URL/api/v1/chat/health" \
  -H "Authorization: Bearer $TOKEN"

echo ""
echo ""

# Test 2: Send message (without finai-agent-api, will return AGENT_UNAVAILABLE)
echo "2. Send message..."
curl -X POST "$API_URL/api/v1/chat/message" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Cek saldo semua akun",
    "conversation_history": []
  }'

echo ""
echo ""

# Test 3: Reset conversation
echo "3. Reset conversation..."
curl -X DELETE "$API_URL/api/v1/chat/reset" \
  -H "Authorization: Bearer $TOKEN"

echo ""
echo ""
echo "=== Test Complete ==="
