package resp

import (
	"encoding/json"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/dbielecki97/url-shortener/pkg/logger"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	switch statusCode {
	case http.StatusNoContent:
	default:
		encodeErr := json.NewEncoder(w).Encode(body)
		if encodeErr != nil {
			logger.Error("could not encode", encodeErr)
		}
	}
}

func Error(w http.ResponseWriter, err errs.RestErr) {
	JSON(w, err.Code(), err)
}
