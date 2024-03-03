package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rinha2024q1crebito"
)

type ErrorResponse struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
}

func parseJsonRequest(r *http.Request, v interface{}) error {
	// // Se for debug, loga o corpo da requisição
	// // https://stackoverflow.com/questions/49745252/reverseproxy-depending-on-the-request-body-in-golang
	// // read all bytes from content body and create new stream using it.
	// bodyBytes, errD := io.ReadAll(r.Body)
	// if errD != nil {
	// 	log.Println(r.Method, r.URL, "Erro ao ler r.Body", errD)
	// 	return fmt.Errorf(
	// 		"error reading request body: %w (%w)",
	// 		errD,
	// 		rinha2024q1crebito.ErrInvalidParameter,
	// 	)
	// }
	// r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return fmt.Errorf(
			"error parsing json request: %w (%w)",
			err,
			rinha2024q1crebito.ErrInvalidParameter)
	}
	return nil
}

func sendJsonResponse(httpStatus int, v interface{}, w http.ResponseWriter) error {
	header := w.Header()
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "application/json")
	}

	w.WriteHeader(httpStatus)

	// reposta sem corpo
	if v == nil {
		return nil
	}

	err := json.NewEncoder(w).Encode(v)
	return err
}

func sendOkJsonResponse(v interface{}, w http.ResponseWriter) error {
	return sendJsonResponse(http.StatusOK, v, w)
}

// func sendCreatedJsonResponse(v interface{}, w http.ResponseWriter) error {
// 	return sendJsonResponse(http.StatusCreated, v, w)
// }

// func sendNoContentResponse(w http.ResponseWriter) error {
// 	w.WriteHeader(http.StatusNoContent)
// 	return nil
// }

func sendErrorResponseWithStatusCode(httpStatus int, vErr error, w http.ResponseWriter) error {
	e := ErrorResponse{
		Title:  vErr.Error(),
		Status: httpStatus,
	}
	return sendJsonResponse(httpStatus, e, w)
}

func sendErrorResponse(vErr error, w http.ResponseWriter) error {
	httpStatus := getErrorHTTPStatus(vErr)

	if httpStatus >= http.StatusInternalServerError {
		log.Println(
			"http request ended with error response",
			"error",
			vErr.Error(),
			"http_status",
			httpStatus,
		)
	}

	return sendErrorResponseWithStatusCode(httpStatus, vErr, w)
}

// getErrorHTTPStatus retorna o código de status HTTP apropriado com base no tipo de erro.
func getErrorHTTPStatus(err error) int {
	switch {
	case errors.Is(err, rinha2024q1crebito.ErrInvalidParameter):
		return http.StatusBadRequest
	case errors.Is(err, rinha2024q1crebito.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, rinha2024q1crebito.ErrUnprocessable):
		return http.StatusUnprocessableEntity
	default:
		// Demais erros são considerados erros internos do servidor
		return http.StatusInternalServerError
	}
}
