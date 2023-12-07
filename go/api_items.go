/*
 * Сервис по обмену вещами Steam
 *
 * API for exchanging virtual items
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

/*
	Наполнить все методы (связь с БД и тд) - то есть настроить логику

Изначально захардкодить вывод json
*/
func ItemsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// fmt.Fprintf(w, string())
	w.WriteHeader(http.StatusOK)
}

func ItemsItemIdDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ItemsItemIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	item_id, ok := vars["item_id"]
	if !ok {
		fmt.Println("id is missing in parameters")
	}
	fmt.Fprintf(w, string("<h1>"+item_id+"</h1>"))
	w.WriteHeader(http.StatusOK)
}

func ItemsPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
