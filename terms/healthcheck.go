package terms

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

func HealthCheckHandler(db *gorm.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if db != nil {
			if err := db.Raw("SELECT 1").Error; err != nil {
				returnError(w, IngestError{100, err.Error()}, 500)
				return
			}
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\": \"ok\"}"))
	}
}
