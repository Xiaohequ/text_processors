param(
    [switch]$Release
)

# Dossier de sortie "build" à la racine du projet
$OutputDir = Join-Path -Path $PSScriptRoot -ChildPath "build"
if (!(Test-Path -Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# Nom de l'artifact
$OutputName = "text_processors.exe"

# Construire la ligne de commande go build
$goArgs = @("build")
if ($Release) {
    # Flags de réduction de taille du binaire
    $goArgs += "-ldflags=-s -w"
}
$goArgs += @("-o", (Join-Path $OutputDir $OutputName), ".")

Write-Host "Exécution: go $($goArgs -join ' ')"

# Lancer la compilation
$processInfo = New-Object System.Diagnostics.ProcessStartInfo
$processInfo.FileName = "go"
$processInfo.Arguments = ($goArgs -join " ")
$processInfo.WorkingDirectory = $PSScriptRoot
$processInfo.RedirectStandardOutput = $true
$processInfo.RedirectStandardError = $true
$processInfo.UseShellExecute = $false

$process = New-Object System.Diagnostics.Process
$process.StartInfo = $processInfo
$process.Start() | Out-Null
$stdout = $process.StandardOutput.ReadToEnd()
$stderr = $process.StandardError.ReadToEnd()
$process.WaitForExit()

if ($stdout) { Write-Host $stdout }
if ($stderr) { Write-Error $stderr }

if ($process.ExitCode -ne 0) {
    Write-Error "La compilation a échoué avec le code $($process.ExitCode)."
    exit $process.ExitCode
}

Write-Host "Artifact généré: $(Join-Path $OutputDir $OutputName)"

