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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tanmaydatta/food/dto"
)

type Endpoint struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
	Methods []string
}

type response struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func RegisterEndpoints(service Service, r *mux.Router) {
	for _, e := range makeHTTPEndpoints(service) {
		r.HandleFunc(e.Path, e.Handler).Methods(e.Methods...)
	}
}

func makeHTTPEndpoints(service Service) []Endpoint {
	return []Endpoint{
		makeHelloEndpoint(service),
		makePredictEndpoint(service),
		makeUploadEndpoint(),
	}
}

func makeHelloEndpoint(service Service) Endpoint {
	return Endpoint{
		Path:    "/hello",
		Methods: []string{http.MethodGet},
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

func makePredictEndpoint(service Service) Endpoint {
	return Endpoint{
		Path:    "/predict",
		Methods: []string{http.MethodPost},
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Body == nil {
				_, _ = w.Write(makeResponseObject(nil, "Empty body"))
				return
			}
			req := dto.PredictReq{}
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				_, _ = w.Write(makeResponseObject(nil, err.Error()))
				return
			}
			resp, err := service.Predict(&req)
			if err != nil {
				_, _ = w.Write(makeResponseObject(nil, err.Error()))
				return
			}
			_, _ = w.Write(makeResponseObject(resp, nil))
		},
	}
}

func makeUploadEndpoint() Endpoint {
	return Endpoint{
		Path:    "/upload",
		Methods: []string{http.MethodPost},
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			// Parse our multipart form, 10 << 20 specifies a maximum
			// upload of 10 MB files.
			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				_, _ = w.Write(makeResponseObject(nil, err.Error()))
				return
			}
			// FormFile returns the first file for the given key `myFile`
			// it also returns the FileHeader so we can get the Filename,
			// the Header and the size of the file
			file, _, err := r.FormFile("image")
			if err != nil {
				fmt.Println("Error Retrieving the File")
				fmt.Println(err)
				_, _ = w.Write(makeResponseObject(nil, err.Error()))
				return
			}
			defer file.Close()

			// Create a temporary file within our temp-images directory that follows
			// a particular naming pattern
			tempFile, err := ioutil.TempFile("/tmp/images", "upload-*.jpg")
			if err != nil {
				fmt.Printf("Error in tempfile")
				fmt.Println(err)
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a
			// byte array
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Printf("Error in readall")
				fmt.Println(err)
			}
			// write this byte array to our temporary file
			_, err = tempFile.Write(fileBytes)
			if err != nil {
				fmt.Printf("Error in write")
				_, _ = w.Write(makeResponseObject(nil, err.Error()))
				return
			}
			// return that we have successfully uploaded our file!
			fullPath := strings.Split(tempFile.Name(), "/")
			_, _ = w.Write(makeResponseObject(fullPath[len(fullPath)-1], nil))
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
