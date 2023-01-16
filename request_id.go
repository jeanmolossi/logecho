package logecho

import (
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // Num of letter idx fitting in 63 bits
)

func generateRequestID() string {
	return generateRequestIDWithCustomLength(12)
}

// generateRequestID will generate randomic string mixed by
// letterBytes chars
func generateRequestIDWithCustomLength(length int) string {
	rand.Seed(time.Now().UnixNano())
	requestID := make([]byte, length)

	for i, cache, remain := (length - 1), rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			requestID[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(requestID)
}

func getXRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// installXRequestID will get request id incoming from headers.
// If it already set in Request just set it into response too.
//
// In case of already has request on response just ignore.
//
// In case of has not request id on incoming headers and response
// headers, generate a new one and set it on response
func installXRequestID(c echo.Context) {
	reqID := c.Request().Header.Get(echo.HeaderXRequestID)

	if reqID != "" {
		c.Response().Header().Set(echo.HeaderXRequestID, reqID)
		return
	}

	if reqID := c.Response().Header().Get(echo.HeaderXRequestID); reqID != "" {
		return
	}

	reqID = generateRequestID()
	c.Response().Header().Set(echo.HeaderXRequestID, reqID)
}
