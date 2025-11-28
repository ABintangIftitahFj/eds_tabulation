# Get the directory where this script is located
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Change to the script directory
Push-Location $ScriptDir

try {
    # Check if .env file exists
    if (Test-Path ".env") {
        Get-Content .env | ForEach-Object {
            if ($_ -match '^(?<k>[^#=\s]+)=(?<v>.+)$') {
                $key = $matches['k'].Trim()
                $value = $matches['v'].Trim()
                Set-Item -Path "Env:$key" -Value $value
            }
        }
    } else {
        Write-Warning ".env file not found in $ScriptDir"
    }
    
    # Check if main.go exists
    if (Test-Path "main.go") {
        go run main.go
    } else {
        Write-Error "main.go file not found in $ScriptDir"
    }
} finally {
    # Return to original directory
    Pop-Location
}