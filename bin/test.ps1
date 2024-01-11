$ErrorActionPreference = "Stop";
trap { $host.SetShouldExit(1) }

Debug "$(gci env:* | sort-object name | Out-String)"

Invoke-Expression "go run github.com/onsi/ginkgo/v2/ginkgo $args --skip-package=db,timeouts"
if ($LastExitCode -ne 0) {
  throw "tests failed"
}
