package bgo

import (
	"encoding/json"
	"net/http"
	"strings"

	graphql "github.com/graph-gophers/graphql-go"
)

// fork from github.com/graph-gophers/graphql-go/relay

type relayHandler struct {
	Schema *graphql.Schema
}

// ServeHTTP func
func (h *relayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.Schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)

	// https://github.com/graph-gophers/graphql-go/pull/207
	if response.Errors != nil {
		// mask panic error
		panicMsg := "graphql: panic occurred"
		for _, rErr := range response.Errors {
			if isPanic := strings.Contains(rErr.Message, panicMsg); isPanic {
				rErr.Message = panicMsg
			}
		}
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}