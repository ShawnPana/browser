$ErrorActionPreference = "Stop"

$Repo = "ShawnPana/browser"
$Binary = "browser"
$InstallDir = "$env:LOCALAPPDATA\browser"

# Detect architecture
$Arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { "amd64" }
    "ARM64" { "arm64" }
    default { Write-Error "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"; exit 1 }
}

# Get version
if (-not $env:VERSION) {
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $Release.tag_name
} else {
    $Version = $env:VERSION
}
if (-not $Version) {
    Write-Error "Failed to determine latest version"
    exit 1
}

Write-Host "Installing $Binary $Version (windows/$Arch)..."

$VersionTrimmed = $Version.TrimStart("v")
$Archive = "${Binary}_${VersionTrimmed}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/$Version/$Archive"
$ChecksumsUrl = "https://github.com/$Repo/releases/download/$Version/checksums.txt"

$TmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $TmpDir | Out-Null

try {
    # Download archive and checksums
    Invoke-WebRequest -Uri $Url -OutFile (Join-Path $TmpDir $Archive)
    Invoke-WebRequest -Uri $ChecksumsUrl -OutFile (Join-Path $TmpDir "checksums.txt")

    # Verify checksum
    $Checksums = Get-Content (Join-Path $TmpDir "checksums.txt")
    $Expected = ($Checksums | Where-Object { $_ -match $Archive }) -replace "\s+.*$", ""
    if (-not $Expected) {
        Write-Error "Checksum not found for $Archive"
        exit 1
    }

    $Actual = (Get-FileHash -Path (Join-Path $TmpDir $Archive) -Algorithm SHA256).Hash.ToLower()
    if ($Expected -ne $Actual) {
        Write-Error "Checksum mismatch!`n  expected: $Expected`n  actual:   $Actual"
        exit 1
    }

    # Extract
    Expand-Archive -Path (Join-Path $TmpDir $Archive) -DestinationPath $TmpDir -Force

    # Install
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir | Out-Null
    }
    Copy-Item (Join-Path $TmpDir "$Binary.exe") -Destination (Join-Path $InstallDir "$Binary.exe") -Force

    # Add to PATH if not already there
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Host "Added $InstallDir to user PATH."
    }

    Write-Host "Installed $Binary to $InstallDir\$Binary.exe"
} finally {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
}
