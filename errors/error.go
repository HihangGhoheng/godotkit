package gdk_error

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type PkgErrorCode uint

type PkgErrorCommon struct {
	ClientMessage string       `json:"message"`
	SystemMessage interface{}  `json:"data"`
	ErrorCode     PkgErrorCode `json:"code"`
	Message       *string      `json:"-"`
	Trace         *string      `json:"-"`
}

type PkgErrorHttp struct {
	PkgErrorCommon
	HttpStatusNumber int    `json:"-"`
	HttpStatusName   string `json:"type"`
}

type ErrorPackages struct {
	appErrCodes  map[PkgErrorCode]*PkgErrorCommon
	httpErrCodes map[PkgErrorCode]int
}

const UNKNOWN_ERROR PkgErrorCode = 0

func (c *PkgErrorCommon) Error() string {
	return fmt.Sprintf("Error: %+v. Trace: %+v", c.Message, c.Trace)
}

func (c *PkgErrorCommon) SetClientMessage(message string) {
	c.ClientMessage = message
}

func (c *PkgErrorCommon) SetSystemMessage(err interface{}) {
	c.SystemMessage = err
}

func (c *PkgErrorCommon) GetHttpStatus() int {
	httpStatus := errPackages.httpErrCodes[c.ErrorCode]
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}
	return httpStatus
}
func (c *PkgErrorCommon) getHttpStatusText(status int) string {
	if text := http.StatusText(status); text != "" {
		upperText := strings.ToUpper(text)
		return strings.ReplaceAll(upperText, " ", "_")
	}
	return "INTERNAL_SERVER_ERROR"
}

var errPackages = &ErrorPackages{
	appErrCodes:  make(map[PkgErrorCode]*PkgErrorCommon),
	httpErrCodes: make(map[PkgErrorCode]int),
}

// Registration list of app errors

func RegisterErrorPackage(
	appErrs map[PkgErrorCode]*PkgErrorCommon,
	httpErrs map[PkgErrorCode]int,
) *ErrorPackages {
	errPackages = &ErrorPackages{
		appErrCodes:  appErrs,
		httpErrCodes: httpErrs,
	}
	errPackages.appErrCodes[UNKNOWN_ERROR] = &PkgErrorCommon{
		ClientMessage: "Unknown error!",
		SystemMessage: "Type of error is not yet registered.",
		ErrorCode:     UNKNOWN_ERROR,
	}
	errPackages.httpErrCodes[UNKNOWN_ERROR] = http.StatusInternalServerError

	return errPackages
}

// Create standard error

func MakeError(code PkgErrorCode, err error) *PkgErrorCommon {
	var (
		errMessage    *string
		errTrace      *string
		clientMessage = "Unknown error!"
		systemMessage = "Type error is not registered"
		errpackages   = errPackages.appErrCodes[code]
	)

	if err != nil {
		e := err.Error()
		errMessage = &e

		ss := fmt.Sprintf("%+v", err)
		errTrace = &ss

		if code == UNKNOWN_ERROR {
			systemMessage = ss
		}
	}

	if errpackages == nil {
		return &PkgErrorCommon{
			ClientMessage: clientMessage,
			SystemMessage: systemMessage,
			ErrorCode:     UNKNOWN_ERROR,
			Message:       errMessage,
			Trace:         errTrace,
		}
	}

	var ne *PkgErrorCommon
	if errors.As(err, &ne) {
		return ne
	}

	return &PkgErrorCommon{
		ClientMessage: errpackages.ClientMessage,
		SystemMessage: errpackages.SystemMessage,
		ErrorCode:     code,
		Message:       errMessage,
		Trace:         errTrace,
	}
}

// Create Http Error

func (c *PkgErrorCommon) ToHttpError() *PkgErrorHttp {
	status := c.GetHttpStatus()
	text := c.getHttpStatusText(status)

	return &PkgErrorHttp{
		PkgErrorCommon:   *c,
		HttpStatusNumber: status,
		HttpStatusName:   text,
	}
}
