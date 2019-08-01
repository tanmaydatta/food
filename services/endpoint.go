// Copyright (c) 2012-2019 Grabtaxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package services

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tanmaydatta/food/dto"
)

type Endpoint struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

type response struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func RegisterEndpoints(service Service, r *mux.Router) {
	for _, e := range makeHTTPEndpoints(service) {
		r.HandleFunc(e.Path, e.Handler)
	}
}

func makeHTTPEndpoints(service Service) []Endpoint {
	return []Endpoint{
		makeHelloEndpoint(service),
	}
}

func makeHelloEndpoint(service Service) Endpoint {
	return Endpoint{
		Path: "/hello",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			keys, ok := r.URL.Query()["name"]

			if !ok {
				_, _ = w.Write(makeResponseObject(nil, "error getting request data"))
				return
			}

			resp, err := service.Hello(&dto.HelloReq{Name: keys[0]})
			if err != nil {
				_, _ = w.Write(makeResponseObject(nil, err.Error()))
				return
			}
			_, _ = w.Write(makeResponseObject(resp, nil))
		},
	}
}

func makeResponseObject(res interface{}, err interface{}) []byte {
	resp := response{
		Result: res,
		Error:  err,
	}

	marshaled, _ := json.Marshal(resp)
	return marshaled
}
