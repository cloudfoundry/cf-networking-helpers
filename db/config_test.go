package db_test

import (
	"io/ioutil"

	"code.cloudfoundry.org/cf-networking-helpers/db"

	"net/url"

	"github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	DATABASE_CA_CERT = `-----BEGIN CERTIFICATE-----
MIIFjzCCA3egAwIBAgIURBuFfs2krcBqtKDwXVXsg0K9dQ4wDQYJKoZIhvcNAQEL
BQAwVzELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEQMA4GA1UEAwwHbXlzcWxDQTAeFw0y
MDA3MjExODMxNDlaFw0zMDA3MjIxODMxNDlaMFcxCzAJBgNVBAYTAkFVMRMwEQYD
VQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBM
dGQxEDAOBgNVBAMMB215c3FsQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIK
AoICAQDMqTAzSlG/akOb1ldMCFJUV+5CjPgXU3JB2vHd2PE3SDESnsairo8b2Mmd
aX4n50LUJaXhgCF5JtGkwAjLKuLn08TwbW/Km9gLLTnCr0gFu3Z5c0Q7aNs5hx00
ma0LuxZnfXWOv4+G26kURU2RniDRTG3L6OwP1cAwrzEykTv7TM4yN1+ZK7vqYt5W
Dd965uWRh+3wfF9UBxKaV2+64Fc+k6eXNl2RpBz++9gq5bXK1mIiD8tEQgNk2j5v
Z3QE0iusGSASiIyNLZ3O+VKGlDxez1QldlUu3WGGNeLCKDRUx+NWXQVXfeMW5RBA
CcuFRAbpBDF6eQj2yREKelj8jJNzQuNr3pqkGzJ0fvQRmI9wUNO34B2gFOMjojtc
PkT1HVVvi92/pOCzvgrXp+G/5DTKOhl2o7IvqEWqqK4i5i3RSstfw7y3QibXkG6g
JdO8emo9hm90AGUhtm6RZNcBozMg3WgSocPDAAJozLNc1MOz6WvV1a/hFeddjD9M
eO3qDuMnpnnN6aOLQFrc2ZVcGzWoGinj0pcbQDFeyFg+rA5LY+E66k9ghBaz2Vl3
PRXQUdMu3ghxcSMcf2Q2c9qs6PJA2e8RhIlCZXDwx8+6Et06SuBpf8hwAAdATaXb
edBwVYWBlDh4bIfAsAH5llxcDYTwn/kydlAAOrZEVKbraZ6S/QIDAQABo1MwUTAd
BgNVHQ4EFgQU4JAx6uVCU97mdF8z9PruwVKJEpIwHwYDVR0jBBgwFoAU4JAx6uVC
U97mdF8z9PruwVKJEpIwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AgEAvrdpdjElBrnW6Kx5wZ2wRNAeWQ+3iun0dJdpohwx8Q3+4zPqWfQR2FI9UOQn
ExPBhos/RM7I1EHL99ZkCzWMEEis7C/ZYLcoBsC3PsAqe2SvSTH8m+Sfn1SJdAp4
QVOT7qEAv3HMrPltq6KxxEeybZW4sV+4Q+0GKw8CXdhIkRtdv/w8suRyn9MszRJn
bvv95RJyz35pmj0R23/aHPjyNYjZQPu4yJbOViDbk5EH2jFc58WPQ9FYu3vLMMMU
BbTHy627LCeYoGJu8xFUzrwav/OeLXPasSsIJcM3Nqi65MWsUd/IxVTLdQhS+/od
R9KOuw5mmWi8dzgVeJN0+dw5KmuwTYuqMMFH9C3wizUS4+F/Of9qGv4zbatH57eM
fIYH/x0Y9qo5Yma7RHE8xG75OIUJ960KomyOKEROWcmIMua1dA+q1MK1G02oKJRs
95YH9uqc1E4ovZCZQZXNf+kQ1RFWf5KKpQxggmQD+SjY0AR7ZPnkYd6aV8U3w4Ku
21vOsp342CXSU4e7q2F4y/aRDOveSLUQ9mg0+aF69IhnGFU7ldL8+3fX/t97R++F
XfKgWjFKtCuyIiiNN9MOf46QowxCHcCu9/Ey0jNvzYfQZFTozo7ZkQuzvAC13BFT
pO9/Xu9jUtG2jFPH/7591N+Y3uJ0E+PcGPaii0fkQLrPs/M=
-----END CERTIFICATE-----`
)

var _ = Describe("Config", func() {
	var (
		config db.Config
	)

	BeforeEach(func() {
		config = db.Config{
			User:         "some-user",
			Password:     "some-password",
			Host:         "some-host",
			Port:         uint16(1234),
			DatabaseName: "some-database",
			Timeout:      5,
		}
	})

	Describe("ConnectionString", func() {
		Context("when the type is postgres", func() {
			BeforeEach(func() {
				config.Type = "postgres"
			})

			It("returns the connection string", func() {
				connectionString, err := config.ConnectionString()
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionString).To(Equal("postgres://some-user:some-password@some-host:1234/some-database?connect_timeout=5000&sslmode=disable"))
			})

			Context("when ssl is required", func() {
				BeforeEach(func() {
					config.RequireSSL = true
					config.CACert = "/tmp/cert"
				})

				Context("when skip_hostname_validation is set", func() {
					BeforeEach(func() {
						config.SkipHostnameValidation = true
					})
					It("sets sslmode to \"require\"", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						connUrl, err := url.Parse(connectionString)
						Expect(err).NotTo(HaveOccurred())
						connQuery := connUrl.Query()
						Expect(connQuery.Get("sslmode")).To(Equal("require"))
					})
				})

				Context("when skip_hostname_validation is not set", func() {
					Context("when ca_cert is empty", func() {
						BeforeEach(func() {
							config.CACert = ""
						})
						It("returns an error", func() {
							_, err := config.ConnectionString()
							Expect(err).To(HaveOccurred())
						})
					})

					It("sets sslmode to \"verify-full\"", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						connUrl, err := url.Parse(connectionString)
						Expect(err).NotTo(HaveOccurred())
						connQuery := connUrl.Query()
						Expect(connQuery.Get("sslmode")).To(Equal("verify-full"))
					})
					It("sets sslrootcert", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						connUrl, err := url.Parse(connectionString)
						Expect(err).NotTo(HaveOccurred())
						connQuery := connUrl.Query()
						Expect(connQuery.Get("sslrootcert")).To(Equal("/tmp/cert"))
					})
				})

			})
		})

		Context("when the type is mysql", func() {
			BeforeEach(func() {
				config.Type = "mysql"
			})

			It("returns the connection string", func() {
				connectionString, err := config.ConnectionString()
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&writeTimeout=5s&sql_mode=%28SELECT+CONCAT%28%40%40sql_mode%2C%27%2CANSI_QUOTES%27%29%29"))
			})

			Context("when require_ssl is enabled", func() {
				BeforeEach(func() {
					config.RequireSSL = true
				})

				AfterEach(func() {
					mysql.DeregisterTLSConfig("some-database-tls")
				})

				Context("success", func() {
					BeforeEach(func() {
						caCertFile, err := ioutil.TempFile("", "")
						Expect(err).NotTo(HaveOccurred())

						_, err = caCertFile.Write([]byte(DATABASE_CA_CERT))
						Expect(err).NotTo(HaveOccurred())

						config.CACert = caCertFile.Name()
					})

					It("returns the amended connection string", func() {
						connectionString, err := config.ConnectionString()
						Expect(err).NotTo(HaveOccurred())
						Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&tls=some-database-tls&writeTimeout=5s&sql_mode=%28SELECT+CONCAT%28%40%40sql_mode%2C%27%2CANSI_QUOTES%27%29%29"))
					})
				})

				Context("when reading the cert file fails", func() {
					BeforeEach(func() {
						config.CACert = "garbage"
					})

					It("returns an error", func() {
						_, err := config.ConnectionString()
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("reading db ca cert file: open garbage: no such file or directory"))
					})
				})

				Context("when adding the cert to the pool fails", func() {
					BeforeEach(func() {
						caCertFile, err := ioutil.TempFile("", "")
						Expect(err).NotTo(HaveOccurred())

						_, err = caCertFile.Write([]byte("garbage"))
						Expect(err).NotTo(HaveOccurred())

						config.CACert = caCertFile.Name()
					})

					It("returns an error", func() {
						_, err := config.ConnectionString()
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError("appending cert to pool from pem - invalid cert bytes"))
					})
				})
			})
		})

		Context("when the type is neither", func() {
			BeforeEach(func() {
				config.Type = "neither"
			})

			It("returns an error", func() {
				_, err := config.ConnectionString()
				Expect(err).To(MatchError("database type 'neither' is not supported"))
			})
		})

		Context("when the timeout is less than 1", func() {
			BeforeEach(func() {
				config.Type = "postgres"
				config.Timeout = 0
			})

			It("returns an error", func() {
				_, err := config.ConnectionString()
				Expect(err).To(MatchError("timeout must be at least 1 second: 0"))
			})
		})
	})
})
