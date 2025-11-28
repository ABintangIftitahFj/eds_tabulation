# Tournament Testing Script
# This script runs all tests for the EDS UPI Tournament System

Write-Host "🧪 Running EDS UPI Tournament System Tests" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Green
Write-Host ""

# Frontend Tests
Write-Host "📱 Running Frontend Tests..." -ForegroundColor Yellow
Set-Location "frontend"

Write-Host "Installing dependencies..." -ForegroundColor Gray
npm install --silent

Write-Host "Running Jest tests..." -ForegroundColor Gray
$frontendResult = npm test -- --passWithNoTests --watchAll=false --coverage
$frontendExitCode = $LASTEXITCODE

Set-Location ".."

Write-Host ""
Write-Host "🔧 Running Backend Logic Tests..." -ForegroundColor Yellow

# Backend Logic Tests
Set-Location "Backend"

Write-Host "Running Go unit tests..." -ForegroundColor Gray
$backendResult = go test -v logic_test.go
$backendExitCode = $LASTEXITCODE

Set-Location ".."

# Test Summary
Write-Host ""
Write-Host "📊 Test Results Summary" -ForegroundColor Cyan
Write-Host "=====================" -ForegroundColor Cyan

if ($frontendExitCode -eq 0) {
    Write-Host "✅ Frontend Tests: PASSED" -ForegroundColor Green
} else {
    Write-Host "❌ Frontend Tests: FAILED" -ForegroundColor Red
}

if ($backendExitCode -eq 0) {
    Write-Host "✅ Backend Tests: PASSED" -ForegroundColor Green
} else {
    Write-Host "❌ Backend Tests: FAILED" -ForegroundColor Red
}

Write-Host ""
Write-Host "🔍 Test Coverage Information" -ForegroundColor Magenta
Write-Host "============================" -ForegroundColor Magenta
Write-Host "Frontend coverage report generated in: frontend/coverage/" -ForegroundColor Gray
Write-Host "View detailed coverage: frontend/coverage/lcov-report/index.html" -ForegroundColor Gray

Write-Host ""
Write-Host "📚 Available Test Commands" -ForegroundColor Blue
Write-Host "==========================" -ForegroundColor Blue
Write-Host ""
Write-Host "Frontend (from frontend/ directory):" -ForegroundColor White
Write-Host "  npm test                    # Run all tests in watch mode" -ForegroundColor Gray
Write-Host "  npm test -- --coverage      # Run tests with coverage" -ForegroundColor Gray
Write-Host "  npm test -- --watchAll=false # Run tests once" -ForegroundColor Gray
Write-Host ""
Write-Host "Backend (from Backend/ directory):" -ForegroundColor White
Write-Host "  go test -v logic_test.go     # Run unit tests" -ForegroundColor Gray
Write-Host "  go test -bench=.             # Run benchmarks" -ForegroundColor Gray
Write-Host "  go test -cover               # Run tests with coverage" -ForegroundColor Gray

if ($frontendExitCode -eq 0 -and $backendExitCode -eq 0) {
    Write-Host ""
    Write-Host "🎉 All tests passed successfully!" -ForegroundColor Green
    exit 0
} else {
    Write-Host ""
    Write-Host "⚠️ Some tests failed. Please check the output above." -ForegroundColor Red
    exit 1
}