package db_test

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"code.cloudfoundry.org/cf-networking-helpers/db"
	"code.cloudfoundry.org/cf-networking-helpers/fakes"
	"github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	DATABASE_CLIENT_CERT = `-----BEGIN CERTIFICATE-----
MIIEOTCCAiECFFQB88eMvRFzig5vh+MJyi0LpnODMA0GCSqGSIb3DQEBCwUAMFcx
CzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRl
cm5ldCBXaWRnaXRzIFB0eSBMdGQxEDAOBgNVBAMMB215c3FsQ0EwHhcNMjAwNzIx
MTgzNjE0WhcNMzAwNzIxMTgzNjE0WjBbMQswCQYDVQQGEwJBVTETMBEGA1UECAwK
U29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMRQw
EgYDVQQDDAtteXNxbENsaWVudDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBAKAg5PUELLrecdLXus0XPVJieZU+KjJNukq0PoRlMlWI+HT+Ibc7IMKqD0XQ
xjXrp1PpQwnRlntCKW51OS4dpBEGjReBNPrzShfgL1AcOQXJtfhQ1yW/KCaHjXNn
ICyhOHV8JgZzUdTanN9iv6SHOGHyRcFtU5pYYLj39LesczXYIABZYV6wj/BDLDo4
5eLnCdGMLWEKBnuUN/5BrIeYxiI/GlT2+zpQ5LnJqcQF4GVq3XSXpfN1WI/idDUV
nej4w1KUe2dPxXgLV7vvpjpxLJrjdLW4NVoV+FOfHNtIshGH1XNa1Rt9d40qHlmJ
80/A0qJmgGAk4+5AU1Mp5MQ/iFkCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAnFbU
k4r96qIsQyPOtk5+MO6HzH+jzNVbdD9+4Wh4vJ14nNQ2jNsBh90A1dpDYd9aIVbq
iYxlA0//CpU0Pj9t/3ymrkBZ6AQnkNZuX/x5yYo5n2AJY0RwIxwviAyWrBqeGIsU
HCz3gHxxk8RpNVIM+PHHAqbvHTqpJJKnH0/1GJZH0yQN3Md+ahqUho9qmEJ4BaPw
tU4QmaRF0TF1PsfgE+2e3WsO4K6L8Nr4TkfzwPuc2q1RUE9ABrLrIb6Z8glck5TO
nAaCsp9qUARbnGUCEQBq85EKXhKN5xlbL10XPbuGtasUdTC+PKV4FVDtSD9oylLt
HruDcMlM5GHNACfvdd6Nz56VvrG9WviAYhsKNTuBIYoNGOfjoV7NADkyqTTT4l7A
JewANT8l9ntuETg6ltbamspCPXcV0OPe8CqY51+nw34arKk6sjk0qUIW8opQBaCA
59zc27dEV+KgKNkjMOUqzcEIaZmqt+sMOfY6deJ62ZeOsuZF5se4Lz+XDmmqpTXn
Gim3GXnxjUDAUZOd88665Y2iirAmG1TcMDek0lBu7/ysuxjBK+Ef3BQ0YONQvzmn
8OgCqZOlf9gdcBb1P6yIT2tIHUJu8D6C0PJuZi2+N6W8D+3tuGdaGxz+jOEDFcEc
mSvPMfv+Qs4rTUvQi9ISXSWS9WDxye35Y/H5Zas=
-----END CERTIFICATE-----`

	DATABASE_CLIENT_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAoCDk9QQsut5x0te6zRc9UmJ5lT4qMk26SrQ+hGUyVYj4dP4h
tzsgwqoPRdDGNeunU+lDCdGWe0IpbnU5Lh2kEQaNF4E0+vNKF+AvUBw5Bcm1+FDX
Jb8oJoeNc2cgLKE4dXwmBnNR1Nqc32K/pIc4YfJFwW1TmlhguPf0t6xzNdggAFlh
XrCP8EMsOjjl4ucJ0YwtYQoGe5Q3/kGsh5jGIj8aVPb7OlDkucmpxAXgZWrddJel
83VYj+J0NRWd6PjDUpR7Z0/FeAtXu++mOnEsmuN0tbg1WhX4U58c20iyEYfVc1rV
G313jSoeWYnzT8DSomaAYCTj7kBTUynkxD+IWQIDAQABAoIBADv4tedx/dKq9xRR
SZwARS4yxXh4xNL8O3Y5FWT+haB1YSBoAXafhYCCsp+iNmoBbTTHGx11Spe1StEc
xCKaZVUfD/6qnI4USj6w7udg+pZJWoa6uClh84aio/6BtBIi+4M80Pz/eblveutL
D51uK1a7pxZvfK1ExkzgSN31j1ytkl5zkSI/BPS26UqshdW599iCCT8AR/bs3r5q
9Umn2VgyH1R2p23aycCoJ5vZOm9P6PNQqG/aMdRvY2zPLvm53fG9VKLPvZHV00Xp
7Ld6E3l1IWjmzW+XaVUdWDC3+lSnJ5JgoyyWwi/49A3aj/79Pg8Wszsc6Zwu2aNV
48DsdIECgYEA0s1frMgHRpkOY0II5ri5DTIrS5Ajks4MzMJQu1LaEwEeSxOD4Vfh
tlLpCTmG+UbYQbP3CjFd44klLNW3Lf5SBqImDsOsc2pHsMYsD9ah9YLn26OYaVhf
lHf4oYu+MPMeiQA9ObuQndBryIhS1i79j364AAO5Wyvu46qup4C9LxECgYEAwnYc
47kL+1a6UWgEmnZUb4M1bxvCjYKnWRf7dWMf+oq8VWiKbwnorq8diNV/jpM0r2Oz
pTvKZ6i+WBboFPAuRqteuTv+LmC1YhA2RGhD34+gnraG/G2hpVy7sVjSYt2gavuO
cQGj6shMW5uZfUut0Wto+Q6PV7o67jOED+AEVMkCgYBqny4ROmtTrf60/aV67VvL
7OAxLAqSKl4XAwHKfbGHgz1LQ+ekhsrHaVAtNkeqtcaUFO6S3T1b5PZEoyQvwz7B
1Cnjtdz3033HT7ThnfH1N+0GDsz8G5LciYHcU84y/AUpzlEpblFLQSwDFdCwPLxL
ZBC1ES0jjCkcRixo1JjbwQKBgQCblQhUIe82LhNvojtcoaO4fE+6D4m+0nO10nw7
VQ121J0N8FAhutnROQX2PmqJ1bjnQmuunYG5Icb4j4srhWZg4CcvKJKa5ID6bmIc
pRb4vN8TXJHvUH9t4B3DLH9W3l7EeXNjcp6E77A38uwA1RXCYZ9g9Ic29ybDAbB9
SwvsEQKBgQDRso0Q2o8YsYIfzdUUdEk/8d0gCVF1LAUyQ63ewAGSKe+iwOkxNhIW
0aF39KEPu+qLIrK9ckHKqKR2FjIBDMTV3kzJoqu6fKFIJ119uQPxqenDWVOUzi/F
M6/bsoX6g47QqJrselETFPnQClDCMJVJq6LRqeeEhqL6EYGbbzumcA==
-----END RSA PRIVATE KEY-----`

	CERTIFICATE_FROM_ANOTHER_CA = `-----BEGIN CERTIFICATE-----
MIIEIzCCAgugAwIBAgIRAKRY04+EGcnAfPpeLf3dxZcwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEAxMHbXlzcWxDQTAeFw0xODEwMTYyMjU2NDdaFw0xODEwMTYyMzU2
NDdaMBUxEzARBgNVBAMTCmV4cGlyZVNvb24wggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCZ/8fZc0q1I03L78hro9jr987Tn6hsJNGod61GiWOHybHezt5i
+SOp7S/fdmLQsRSopUOSmlAH6ta5QbffGbtHY2NJQKJq7N8KRt6aSfbHxPDG96Rp
Q0OZZLyiEaFz2jECoTjqwZX9duG5wA1/AVnZEKqnbAWdIWP9AOTzwdJ/ne4CLzyj
Lm/HUNi9xsZvU5xgb8ZSW3z8SOf39UedocmDcA/rTZWAkO6ELPvx4KD6t5aBC4ir
k7tGveQFxTvziZr3lNZk+NTX2OWUrz5yoH/nMiXtHe4JuytFsN5DYF1f6/3Fxl19
AhkCkxTj238/FFLID34W7mfZbgN59ByBgnPxAgMBAAGjcTBvMA4GA1UdDwEB/wQE
AwIDuDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYEFJuI
c/+CukBgDTwmU6A+6++WwzaKMB8GA1UdIwQYMBaAFFmJZ4pg0d2UuDKpQzF3XQAM
LandMA0GCSqGSIb3DQEBCwUAA4ICAQAh7+THwf0fe3syAaPVqnpx2kswUAqP9VTw
waxXswwp632JnQa9vctuVBQ7DNwOHSixaNlM7yR+w1FlubwLzNRR5EXOgi2kl5Le
mewKBmJLpMwkmAbpCUB2B2ofJJguMe0JVQC6OC3eA3JsTc1/FtqJ4H1+RD5xT6hx
uOxla3zwfynYD4WdRMAosYVJouCScgWJpK+MWEkMCx94GUcO4Ik9acWhzBcdgaUG
qjbtTq5dHgVwernhJaiuUC2R5wEvb3rkhav2TYHJucFm0NHFbMCCYNbFAp1t1OyW
hiNrGtUGN2jBoFZ9OEZaWuY00mKs0Elp5/ugHQ5hW6HXam/4Fh95PMBR1QC+c5AC
AhdCYEXpZXkjCe5vnXHegBxAMV2FU33G9rPWWAi76sBlqjApGaYfbYJW63bhEOZT
AtnHlrPVw/GM16KkzMEEbi4lRvY4F3F2FJ+LZSMKMNs9aX/CAAWs9up3n7PcePP0
fV70C2hVtCJbIfRPaWvrVAAktBP9xLTnzUvzijPLMEJ9o45vWdrtvyBFknQCpMts
lw6sWU26m2gvxs6CcX3yt0bt8SxjqyulqrOdFCVSjZbGMDaIamdEKnC6k5ySyizn
SM2qNm+nV5FhjsyMyzs6OuCNEZGDAqklWBAHHqLncb6elO9NZgDysB/xn6jS+zqT
F1Y5M6wvLA==
-----END CERTIFICATE-----`
)

var _ = Describe("MySQLConnectionStringBuilder", func() {
	Describe("Build", func() {
		var (
			mysqlConnectionStringBuilder *db.MySQLConnectionStringBuilder
			mySQLAdapter                 *fakes.MySQLAdapter
			config                       db.Config
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

			mySQLAdapter = &fakes.MySQLAdapter{}

			mysqlConnectionStringBuilder = &db.MySQLConnectionStringBuilder{
				MySQLAdapter: mySQLAdapter,
			}
			mySQLAdapter.ParseDSNStub = func(dsn string) (cfg *mysql.Config, err error) {
				return mysql.ParseDSN(dsn)
			}
		})

		It("builds a connection string", func() {
			connectionString, err := mysqlConnectionStringBuilder.Build(config)
			Expect(err).NotTo(HaveOccurred())
			Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&writeTimeout=5s&sql_mode=%28SELECT+CONCAT%28%40%40sql_mode%2C%27%2CANSI_QUOTES%27%29%29"))
		})

		Context("when mysql.ParseDSN can't parse the connection string", func() {
			BeforeEach(func() {
				mySQLAdapter.ParseDSNReturns(nil, errors.New("foxtrot"))
			})

			It("returns an error", func() {
				_, err := mysqlConnectionStringBuilder.Build(config)
				Expect(err).To(MatchError("parsing db connection string: foxtrot"))
			})
		})

		Context("when requiring ssl", func() {
			var (
				caCertPool *x509.CertPool
			)

			BeforeEach(func() {
				caCertFile, err := ioutil.TempFile("", "")
				_, err = caCertFile.Write([]byte(DATABASE_CA_CERT))
				Expect(err).NotTo(HaveOccurred())

				config.RequireSSL = true
				config.CACert = caCertFile.Name()

				caCertPool = x509.NewCertPool()
				ok := caCertPool.AppendCertsFromPEM([]byte(DATABASE_CA_CERT))
				Expect(ok).To(BeTrue())
			})

			It("builds a tls connection string", func() {
				connectionString, err := mysqlConnectionStringBuilder.Build(config)
				Expect(err).NotTo(HaveOccurred())
				Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&tls=some-database-tls&writeTimeout=5s&sql_mode=%28SELECT+CONCAT%28%40%40sql_mode%2C%27%2CANSI_QUOTES%27%29%29"))

				Expect(mySQLAdapter.RegisterTLSConfigCallCount()).To(Equal(1))
				passedTLSConfigName, passedTLSConfig := mySQLAdapter.RegisterTLSConfigArgsForCall(0)
				Expect(passedTLSConfigName).To(Equal("some-database-tls"))
				Expect(passedTLSConfig.InsecureSkipVerify).To(Equal(false))
				Expect(passedTLSConfig.RootCAs.Subjects()).To(Equal(caCertPool.Subjects()))
			})

			Context("when SkipHostnameValidation is true", func() {
				BeforeEach(func() {
					config.SkipHostnameValidation = true
				})

				It("builds tls config skipping hostname", func() {
					connectionString, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).NotTo(HaveOccurred())
					Expect(connectionString).To(Equal("some-user:some-password@tcp(some-host:1234)/some-database?parseTime=true&readTimeout=5s&timeout=5s&tls=some-database-tls&writeTimeout=5s&sql_mode=%28SELECT+CONCAT%28%40%40sql_mode%2C%27%2CANSI_QUOTES%27%29%29"))

					Expect(mySQLAdapter.RegisterTLSConfigCallCount()).To(Equal(1))
					passedTLSConfigName, passedTLSConfig := mySQLAdapter.RegisterTLSConfigArgsForCall(0)
					Expect(passedTLSConfigName).To(Equal("some-database-tls"))
					Expect(passedTLSConfig.InsecureSkipVerify).To(BeTrue())
					Expect(passedTLSConfig.RootCAs.Subjects()).To(Equal(caCertPool.Subjects()))
					Expect(passedTLSConfig.Certificates).To(BeNil())
					// impossible to assert VerifyPeerCertificate is set to a specfic function
					Expect(passedTLSConfig.VerifyPeerCertificate).NotTo(BeNil())
				})
			})

			Context("when it can't read the ca cert file", func() {
				BeforeEach(func() {
					config.CACert = "/foo/bar"
				})

				It("returns an error", func() {
					_, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).To(MatchError("reading db ca cert file: open /foo/bar: no such file or directory"))
				})
			})

			Context("when it can't append the ca cert to the cert pool", func() {
				BeforeEach(func() {
					caCertFile, err := ioutil.TempFile("", "")
					_, err = caCertFile.Write([]byte("bad cert"))
					Expect(err).NotTo(HaveOccurred())

					config.CACert = caCertFile.Name()
				})

				It("returns an error", func() {
					_, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).To(MatchError("appending cert to pool from pem - invalid cert bytes"))
				})
			})

			Context("when it can't register TLS config", func() {
				BeforeEach(func() {
					mySQLAdapter.RegisterTLSConfigReturns(errors.New("bad things happened"))
				})

				It("retruns an error", func() {
					_, err := mysqlConnectionStringBuilder.Build(config)
					Expect(err).To(MatchError("registering mysql tls config: bad things happened"))
				})
			})
		})
	})

	Describe("VerifyCertificatesIgnoreHostname", func() {
		var (
			caCertPool *x509.CertPool
		)

		BeforeEach(func() {
			caCertPool = x509.NewCertPool()
			ok := caCertPool.AppendCertsFromPEM([]byte(DATABASE_CA_CERT))
			Expect(ok).To(BeTrue())
		})

		It("verifies that provided certificates are valid", func() {
			block, _ := pem.Decode([]byte(DATABASE_CLIENT_CERT))

			err := db.VerifyCertificatesIgnoreHostname([][]byte{
				block.Bytes,
			}, caCertPool)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when raw certs are not parsable", func() {
			It("returns an error", func() {
				err := db.VerifyCertificatesIgnoreHostname([][]byte{
					[]byte("foo"),
					[]byte("bar"),
				}, nil)
				Expect(err.Error()).To(ContainSubstring("tls: failed to parse certificate from server: x509: malformed certificate"))
			})
		})

		Context("when verifying an expired cert", func() {
			It("returns an error", func() {
				block, _ := pem.Decode([]byte(CERTIFICATE_FROM_ANOTHER_CA))

				err := db.VerifyCertificatesIgnoreHostname([][]byte{
					block.Bytes,
				}, caCertPool)

				Expect(err.Error()).To(ContainSubstring("x509: certificate has expired or is not yet valid"))
			})
		})
	})
})
