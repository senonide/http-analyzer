package ja3

import (
	"crypto/md5" //nolint:gosec
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"strings"
)

func GenerateJA3(hello *tls.ClientHelloInfo) string {
	tlsVersion := fmt.Sprintf("%d", hello.SupportedVersions[0])

	cipherSuites := []string{}
	for _, suite := range hello.CipherSuites {
		cipherSuites = append(cipherSuites, fmt.Sprintf("%d", suite))
	}

	extensions := []string{}
	if hello.SupportedCurves != nil {
		extensions = append(extensions, "10")
	}
	if hello.SupportedPoints != nil {
		extensions = append(extensions, "11")
	}

	supportedCurves := []string{}
	for _, curve := range hello.SupportedCurves {
		supportedCurves = append(supportedCurves, fmt.Sprintf("%d", curve))
	}

	supportedPoints := []string{}
	for _, pointFormat := range hello.SupportedPoints {
		supportedPoints = append(supportedPoints, fmt.Sprintf("%d", pointFormat))
	}

	ja3String := strings.Join([]string{
		tlsVersion,
		strings.Join(cipherSuites, "-"),
		strings.Join(extensions, "-"),
		strings.Join(supportedCurves, "-"),
		strings.Join(supportedPoints, "-"),
	}, ",")

	ja3Hash := md5.Sum([]byte(ja3String)) //nolint:gosec
	return hex.EncodeToString(ja3Hash[:])
}
