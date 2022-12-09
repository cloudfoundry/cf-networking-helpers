$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

cd pull-request/

Write-Host "PATH = $env:PATH"
Write-Host "Running unit tests for packages targeted to build on Windows..."
$env:SQL_FLAVOR="mysql"


go run github.com/onsi/ginkgo/ginkgo -r -keepGoing -trace -randomizeAllSpecs -progress -race `
  ./healthchecker

if ($LastExitCode -ne 0) {
  Write-Host "cf-networking-helpers unit tests failed"
  exit 1
} else {
  Write-Host "cf-networking-helpers unit tests passed"
  exit 0
}
