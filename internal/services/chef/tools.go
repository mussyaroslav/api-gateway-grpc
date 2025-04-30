package chef

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) decideByErrorGRPC(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch status.Code(err) {
	case codes.NotFound:
		s.ErrorBadRequest(c, errResNotFound)
	case codes.InvalidArgument:
		s.ErrorBadRequest(c, errInvalidArgument)
	case codes.Unauthenticated:
		s.ErrorBadRequest(c, errUnauthenticated)
	case codes.PermissionDenied:
		s.ErrorBadRequest(c, errPermissionDenied)
	case codes.DeadlineExceeded:
		s.ErrorBadRequest(c, errDeadlineExceeded)
	default:
		s.ErrorServerError(c, errSrvBadRequest)
	}
}

func loadCerts(crtBase64, keyBase64 string) (tls.Certificate, error) {
	crtBytes, err := base64.StdEncoding.DecodeString(crtBase64)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return tls.Certificate{}, err
	}

	cert, err := tls.X509KeyPair(crtBytes, keyBytes)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}
