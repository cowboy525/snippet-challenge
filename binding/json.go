package binding

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/model"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) (map[string]interface{}, *model.AppError) {
	if req == nil || req.Body == nil {
		return nil, model.InvalidRequestBodyError()
	}
	body, _ := ioutil.ReadAll(req.Body)

	var request map[string]interface{}
	json.NewDecoder(bytes.NewReader(body)).Decode(&request)

	if err := decodeJSON(bytes.NewReader(body), obj); err != nil {
		return nil, err
	}
	return request, validate(obj, request)
}

func (jsonBinding) BindArray(req *http.Request, obj interface{}) ([]map[string]interface{}, *model.AppError) {
	if req == nil || req.Body == nil {
		return nil, model.ParseError("decodeJSON", "api.request.malformed", nil, "")
	}
	body, _ := ioutil.ReadAll(req.Body)

	var request []map[string]interface{}
	json.NewDecoder(bytes.NewReader(body)).Decode(&request)

	if err := decodeJSON(bytes.NewReader(body), obj); err != nil {
		return nil, err
	}
	return request, validate(obj, request)
}

func decodeJSON(r io.Reader, obj interface{}) *model.AppError {
	decoder := json.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(&obj); err != nil {
		mlog.Error(err.Error())
		return model.ParseError("decodeJSON", "api.request.malformed", nil, "")
	}
	return nil
}
