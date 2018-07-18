package main

import (
	"io/ioutil"
	"os"
	"time"

	// internal
	"github.com/sniperkit/iris"
)

// a test listenTLS for our localhost
func listenTLS(app *iris.Application) {

	const (
		testTLSCert = `-----BEGIN CERTIFICATE-----
MIIDBTCCAe2gAwIBAgIJAOYzROngkH6NMA0GCSqGSIb3DQEBBQUAMBkxFzAVBgNV
BAMMDmxvY2FsaG9zdDo4MDgwMB4XDTE3MDIxNzAzNDM1NFoXDTI3MDIxNTAzNDM1
NFowGTEXMBUGA1UEAwwObG9jYWxob3N0OjgwODAwggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQCfsiVHO14FpKsi0pvBv68oApQm2MO+dCvq87sDU4E0QJhG
KV1RCUmQVypChEqdLlUQsopcXSyKwbWoyg1/KNHYO3DHMfePb4bC1UD2HENq7Ph2
8QJTEi/CJvUB9hqke/YCoWYdjFiI3h3Hw8q5whGO5XR3R23z69vr5XxoNlcF2R+O
TdkzArd0CWTZS27vbgdnyi9v3Waydh/rl+QRtPUgEoCEqOOkMSMldXO6Z9GlUk9b
FQHwIuEnlSoVFB5ot5cqebEjJnWMLLP83KOCQekJeHZOyjeTe8W0Fy1DGu5fvFNh
xde9e/7XlFE//00vT7nBmJAUV/2CXC8U5lsjLEqdAgMBAAGjUDBOMB0GA1UdDgQW
BBQOfENuLn/t0Z4ZY1+RPWaz7RBH+TAfBgNVHSMEGDAWgBQOfENuLn/t0Z4ZY1+R
PWaz7RBH+TAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBBQUAA4IBAQBG7AEEuIq6
rWCE5I2t4IXz0jN7MilqEhUWDbUajl1paYf6Ikx5QhMsFx21p6WEWYIYcnWAKZe2
chAgnnGojuxdx0qjiaH4N4xWGHsWhaesnIF1xJepLlX3kJZQURvRxM4wlljlQPIb
9tqzKP131K1HDqplAtp7nWQ72m3J0ZfzH0mYIUxuaS/uQIVtgKqdilwy/VE5dRZ9
QFIb4G9TnNThXMqgTLjfNr33jVbTuv6fzKHYNbCkP3L10ydEs/ddlREmtsn9nE8Q
XCTIYXzA2kr5kWk7d3LkUiSvu3g2S1Ol1YaIKaOQyRveseCGwR4xohLT+dPUW9dL
3hDVLlwE3mB3
-----END CERTIFICATE-----
`
		testTLSKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAn7IlRzteBaSrItKbwb+vKAKUJtjDvnQr6vO7A1OBNECYRild
UQlJkFcqQoRKnS5VELKKXF0sisG1qMoNfyjR2DtwxzH3j2+GwtVA9hxDauz4dvEC
UxIvwib1AfYapHv2AqFmHYxYiN4dx8PKucIRjuV0d0dt8+vb6+V8aDZXBdkfjk3Z
MwK3dAlk2Utu724HZ8ovb91msnYf65fkEbT1IBKAhKjjpDEjJXVzumfRpVJPWxUB
8CLhJ5UqFRQeaLeXKnmxIyZ1jCyz/NyjgkHpCXh2Tso3k3vFtBctQxruX7xTYcXX
vXv+15RRP/9NL0+5wZiQFFf9glwvFOZbIyxKnQIDAQABAoIBAEzBx4ExW8PCni8i
o5LAm2PTuXniflMwa1uGwsCahmOjGI3AnAWzPRSPkNRf2a0q8+AOsMosTphy+umi
FFKmQBZ6m35i2earaE6FSbABbbYbKGGi/ccH2sSrDOBgdfXRTzF8eiSBrJw8hnvZ
87rNOLtCNnSOdJ7lItODfgRo+fLo4uQenJ8VONYwtwm1ejn8qLXq8O5zF66IYUD6
gAzqOiAWumgZL0tEmndeQ+noe4STpJZlOjiCsA12NiJaKDDeDIn5A/pXce+bYNfJ
k4yoroyq/JXBkhyuZDvX9vYp5AA+Q68h8/KmsKkifUgSGSHun5/80lYyT/f60TLX
PxT9GYECgYEA0s8qck7L29nBBTQ6IPF3GHGmqiRdfH+qhP/Jn4NtoW3XuVe4A15i
REq1L8WAiOUIBnBaD8HzbeioqJJYx1pu7x9h/GCNDhdBfwhTjnBe+JjfLqvJKnc0
HUT5wj4DVqattxKzUW8kTRBSWtVremzeffDo+EL6dnR7Bc02Ibs4WpUCgYEAwe34
Uqhie+/EFr4HjYRUNZSNgYNAJkKHVxk4qGzG5VhvjPafnHUbo+Kk/0QW7eIB+kvR
FDO8oKh9wTBrWZEcLJP4jDIKh4y8hZTo9B8EjxFONXVxZlOSYuGjheL8AiLzE7L9
C1spaKMM/MyxAXDRHpG/NeEgXM7Kn6kUGwJdNekCgYAshLNiEGHcu8+XWcAs1NFh
yB56L9PORuerzpi1pvuv65JzAaNKktQNt/krbXoHbtaTBYb/bOYLf+aeMsmsz9w9
g1MeCQXAxAiA2zFKE1D7Ds2S/ZQt8559z+MusgnicrCcyMY1nFL+M0QxCoD4CaWy
0v1f8EUUXuTcBMo5tV/hQQKBgDoBBW8jsiFDu7DZscSgOde00QZVzZAkAfsJLisi
LfNXGjZdZawUUuoX1iYLpZgNK25D0wtp1hdvjf2Ej/dAMd8bexHjvcaBT7ncqjiq
NmDcWjofIIXspTIyLwjStXGmJnJT7N/CqoYDjtTmHGND7Shpi3mAFn/r0isjFUJm
2J5RAoGALuGXxzmSRWmkIp11F/Qr3PBFWBWkrRWaH2TRLMhrU/wO8kCsSyo4PmAZ
ltOfD7InpDiCu43hcDPQ/29FUbDnmAhvMnmIQuHXGgPF/LhqEhbKPA/o/eZdQVCK
QG+tmveBBIYMed5YbWstZu/95lIHF+u8Hl+Z6xgveozfE5yqiUA=
-----END RSA PRIVATE KEY-----
    `
	)

	// create the key and cert files on the fly, and delete them when this test finished
	certFile, ferr := ioutil.TempFile("", "cert")

	if ferr != nil {
		panic(ferr)
	}

	keyFile, ferr := ioutil.TempFile("", "key")
	if ferr != nil {
		panic(ferr)
	}

	certFile.WriteString(testTLSCert)
	keyFile.WriteString(testTLSKey)

	// https://localhost
	app.Run(iris.TLS("localhost:443", certFile.Name(), keyFile.Name()))

	certFile.Close()
	time.Sleep(50 * time.Millisecond)
	os.Remove(certFile.Name())

	keyFile.Close()
	time.Sleep(50 * time.Millisecond)
	os.Remove(keyFile.Name())
}
