$ErrorActionPreference = "Stop"

$repository = "nihitdev/drop"
$installDirectory = Join-Path $env:LOCALAPPDATA "Programs\drop"

$architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
switch ($architecture) {
    "X64" { $assetName = "drop-windows-amd64.exe" }
    "Arm64" { $assetName = "drop-windows-arm64.exe" }
    default { throw "Unsupported Windows architecture: $architecture" }
}

$releaseBase = "https://github.com/$repository/releases/latest/download"
$temporaryDirectory = Join-Path ([System.IO.Path]::GetTempPath()) ("drop-install-" + [guid]::NewGuid())
$downloadPath = Join-Path $temporaryDirectory $assetName
$checksumPath = Join-Path $temporaryDirectory "checksums.txt"

try {
    New-Item -ItemType Directory -Path $temporaryDirectory | Out-Null
    Write-Host "Downloading $assetName..."
    Invoke-WebRequest -UseBasicParsing -Uri "$releaseBase/$assetName" -OutFile $downloadPath
    Invoke-WebRequest -UseBasicParsing -Uri "$releaseBase/checksums.txt" -OutFile $checksumPath

    $checksumLine = Get-Content $checksumPath | Where-Object { $_ -match "\s+$([regex]::Escape($assetName))$" }
    if (-not $checksumLine) {
        throw "No checksum found for $assetName"
    }
    $expectedChecksum = ($checksumLine -split "\s+")[0]
    $actualChecksum = (Get-FileHash -Algorithm SHA256 $downloadPath).Hash
    if ($actualChecksum -ne $expectedChecksum) {
        throw "Checksum verification failed for $assetName"
    }

    New-Item -ItemType Directory -Force -Path $installDirectory | Out-Null
    Copy-Item -Force $downloadPath (Join-Path $installDirectory "drop.exe")

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $pathEntries = @($userPath -split ";" | Where-Object { $_ })
    if ($pathEntries -notcontains $installDirectory) {
        $newUserPath = (($pathEntries + $installDirectory) -join ";")
        [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
        $env:Path = "$env:Path;$installDirectory"
        Write-Host "Added $installDirectory to your user PATH."
    }

    Write-Host "drop installed successfully. Open a new terminal and run: drop <file>"
}
finally {
    if (Test-Path $temporaryDirectory) {
        Remove-Item -Recurse -Force $temporaryDirectory
    }
}
