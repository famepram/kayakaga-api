@echo off
REM Script to generate Swagger documentation for Kayakaga API (Windows)

echo 🚀 Generating Swagger documentation for Kayakaga API...

REM Check if swag is installed
where swag >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo 📦 Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
)

REM Generate documentation
echo 📝 Running swag init...
swag init

REM Check if generation was successful
if %ERRORLEVEL% EQU 0 (
    echo ✅ Swagger documentation generated successfully!
    echo.
    echo 📚 Generated files:
    echo   - docs\docs.go
    echo   - docs\swagger.json
    echo   - docs\swagger.yaml
    echo.
    echo 🌐 Access Swagger UI at: http://localhost:8080/swagger/index.html
    echo.
    echo 💡 To import into Postman:
    echo   1. Start the API server
    echo   2. In Postman: File → Import → From URL
    echo   3. Enter: http://localhost:8080/swagger/doc.json
) else (
    echo ❌ Failed to generate Swagger documentation
    echo Make sure you have swag installed: go install github.com/swaggo/swag/cmd/swag@latest
    exit /b 1
)

pause
