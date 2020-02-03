package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/graph-gophers/graphql-go"
	"github.com/sentadmedia/elf/fw"
)

func NewAuthMiddleWare(schema *graphql.Schema, store sessions.Store, cookieName string, logger fw.Logger) fw.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var params struct {
				Query         string                 `json:"query"`
				OperationName string                 `json:"operationName"`
				Variables     map[string]interface{} `json:"variables"`
			}

			buf, _ := ioutil.ReadAll(r.Body)

			if err := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf))).Decode(&params); err != nil {
				logger.Error(fmt.Errorf("Unable to decode request body (%s)", err))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if strings.HasPrefix(params.Query, "mutation") && strings.Contains(params.Query, "signIn") {
				session, _ := store.Get(r, cookieName)
				logger.Debug(fmt.Sprintf("Session Before=%+v", session))
				ctx := r.Context()
				session.Values["authenticated"] = true
				if err := session.Save(r, w); err != nil {
					logger.Error(fmt.Errorf("Unable to save session. (%s)", err))
				} else {
					logger.Debug(fmt.Sprintf("Session After=%+v", session))
				}
				response := schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
				responseJSON, err := json.Marshal(response)
				if err != nil {
					logger.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(responseJSON)
				return
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			handler.ServeHTTP(w, r)
		})
	}
}
