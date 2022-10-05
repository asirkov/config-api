package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/santhosh-tekuri/jsonconfig/v5"

	"github.com/savsgio/atreugo/v11"
)

type MsTimestamp time.Time

func HandleError(ctx *atreugo.RequestCtx, err error) error {
	response := NewErrorResponse(err)

	switch err.(type) {
	case *InfoNotFound:
		AddResponse(ctx, response, http.StatusNotFound)
	case *ErrorNotModified:
		AddResponse(ctx, response, http.StatusNotModified)
	case *ErrorDuplicateEntries:
		AddResponse(ctx, response, http.StatusMultipleChoices)
	case *strconv.NumError, *json.SyntaxError, *json.InvalidUnmarshalError, *json.UnmarshalTypeError, *WarningValidation:
		AddResponse(ctx, response, http.StatusBadRequest)
	case *jsonconfig.ValidationError:
		AddResponse(ctx, err, http.StatusBadRequest)
	default:
		AddResponse(ctx, response, http.StatusInternalServerError)
	}
	return err
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	errorResponse := ErrorResponse{
		Message: err.Error(),
	}
	return errorResponse
}

func AddResponse(ctx *atreugo.RequestCtx, body interface{}, code int) {
	ctx.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	ctx.Response.Header.Set("X-Content-Type-Options", "nosniff")

	ctx.JSONResponse(body, code)
}

type Pagination struct {
	limit  int
	offset int
}

func (r *Pagination) GetLimit() int {
	if r.limit <= 0 {
		return paginationDefaultLimit
	}

	if r.limit > paginationMaxLimit {
		return paginationMaxLimit
	}
	return r.limit
}

func (r *Pagination) GetOffset() int {
	if r.offset <= 0 {
		return 0
	}
	return r.offset
}

const (
	paginationDefaultLimit = 100
	paginationMaxLimit     = 1000
)

func GetPagination(ctx *atreugo.RequestCtx) Pagination {
	limit := string(ctx.FormValue("limit"))
	offset := string(ctx.FormValue("offset"))

	pagination := Pagination{}
	if limit != "" {
		limit, err := strconv.Atoi(limit)
		if err == nil {
			pagination.limit = limit
		}
	}
	if offset != "" {
		offset, err := strconv.Atoi(offset)
		if err == nil {
			pagination.offset = offset
		}
	}

	return pagination
}

type CollectionMeta struct {
	Total int `json:"total"`
}
