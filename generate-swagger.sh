#!/bin/bash

# Script to generate Swagger documentation for Kayakaga API

echo "🚀 Generating Swagger documentation for Kayakaga API..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "📦 Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate documentation
echo "📝 Running swag init..."
swag init

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo ""
    echo "📚 Generated files:"
    echo "  - docs/docs.go"
    echo "  - docs/swagger.json"
    echo "  - docs/swagger.yaml"
    echo ""
    echo "🌐 Access Swagger UI at: http://localhost:8080/swagger/index.html"
    echo ""
    echo "💡 To import into Postman:"
    echo "  1. Start the API server"
    echo "  2. In Postman: File → Import → From URL"
    echo "  3. Enter: http://localhost:8080/swagger/doc.json"
else
    echo "❌ Failed to generate Swagger documentation"
    echo "Make sure you have swag installed: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi
