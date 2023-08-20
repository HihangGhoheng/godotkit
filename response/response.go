package gdk_response

import (
	"encoding/json"
	"errors"
	gdk_error "github.com/HihangGhoheng/godotkit/errors"
	gdk_helpers "github.com/HihangGhoheng/godotkit/helpers"
	"net/http"
)

type MetaResponse struct {
	Page      int `json:"page"`
	PerPage   int `json:"perPage"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}

type HtmlResponse struct {
	Message string        `json:"message"`
	Meta    *MetaResponse `json:"meta,omitempty"`
	Data    interface{}   `json:"data"`
}

type ResponseHttpImpl interface {
	HttpSuccess(
		w http.ResponseWriter,
		message string,
		data interface{},
		meta *MetaResponse,
	) error
	HttpError(w http.ResponseWriter, err error) error
	MakeMeta(page int, perpage int, total int64) *MetaResponse
}

type responseHttp struct{}

func NewResponseHttp() ResponseHttpImpl {
	return &responseHttp{}
}

func (res *responseHttp) HttpSuccess(
	w http.ResponseWriter,
	message string,
	data interface{},
	meta *MetaResponse,
) error {
	if m, err := json.Marshal(
		HtmlResponse{
			Message: message,
			Meta:    meta,
			Data:    data,
		},
	); err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(m)
	} else {
		gdk_helpers.FailOnError(err, "Failed to marshal the json response")
	}
	return nil
}

func (res *responseHttp) HttpError(
	w http.ResponseWriter,
	err error,
) error {
	var httpError gdk_error.PkgErrorHttp
	var ce *gdk_error.PkgErrorCommon
	if errors.As(err, &ce) {
		httpError = *ce.ToHttpError()
	} else {
		httpError = *gdk_error.MakeError(gdk_error.UNKNOWN_ERROR, err).ToHttpError()
	}

	if m, e := json.Marshal(&httpError); e == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpError.GetHttpStatus())
		_, _ = w.Write(m)
	} else {
		gdk_helpers.FailOnError(err, "Failed to marshal the json response")
	}
	return nil
}

func (res *responseHttp) MakeMeta(
	page int,
	perPage int,
	total int64,
) *MetaResponse {
	totalPage := float64(total) / float64(perPage)
	return &MetaResponse{
		Page:      page,
		PerPage:   perPage,
		Total:     int(total),
		TotalPage: int(totalPage),
	}
}
