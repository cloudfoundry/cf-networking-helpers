$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

cd pull-request/

Write-Host "PATH = $env:PATH"
Write-Host "Running unit tests for packages targeted to build on Windows..."
$env:SQL_FLAVOR="mysql"


go run github.com/onsi/ginkgo/v2/ginkgo -r -keep-going -trace -randomize-all -race -r `
  --skip-package=db,timeouts

if ($LastExitCode -ne 0) {
  Write-Host "cf-networking-helpers unit tests failed"
  exit 1
} else {
  Write-Host "cf-networking-helpers unit tests passed"
  exit 0
}
