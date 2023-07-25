package http

import (
	"encoding/json"
	"github.com/ahlixinjie/mongoose/transport/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

type Codec interface {
	EncodeResponse(r *http.Request, w http.ResponseWriter, respValue interface{}, err error) error
	Decode(r *http.Request, requestValue interface{}) error
}

type baseRequest struct {
	Body json.RawMessage `json:"body"`
}
type baseResponse struct {
	RequestID string          `json:"request_id"`
	Code      codes.Code      `json:"code,omitempty"`
	Message   string          `json:"msg,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

type defaultCodec struct {
	code2Http map[codes.Code]int
}

func (d defaultCodec) EncodeResponse(
	r *http.Request, w http.ResponseWriter, respValue interface{}, respErr error) (err error) {
	resp := baseResponse{
		RequestID: metadata.ValueFromIncomingContext(r.Context(), common.HeaderRequestID)[0],
		Code:      0,
		Message:   "",
		Data:      nil,
	}
	w.Header().Set("Content-type", "application/json")

	respByte, err := json.Marshal(respValue)
	resp.Data = respByte

	s, _ := status.FromError(respErr)
	resp.Code = s.Code()
	resp.Message = s.Message()
	w.WriteHeader(d.code2Http[resp.Code])
	return json.NewEncoder(w).Encode(&resp)
}

func (d defaultCodec) Decode(r *http.Request, requestValue interface{}) error {
	req := baseRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	if len(req.Body) == 0 {
		return nil
	}

	return json.Unmarshal(req.Body, requestValue)
}
