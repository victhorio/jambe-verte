package internal

import (
	"fmt"
	"net/http"
)

func WriteInternalError(w http.ResponseWriter, errorCode string) {
	http.Error(w, fmt.Sprintf("Internal Error. Code: %s", errorCode), http.StatusInternalServerError)
}
