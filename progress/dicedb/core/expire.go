package core

import (
	"time" 
	"log"
)

func deleteExpiredSample() float32 {
	var limit int = 20
	var expiredCount int = 0

	for k, v := range store {
		if v.ExpiresAt != -1 {
			// expiry date is set
			limit--
			if v.ExpiresAt <= time.Now().UnixMicro() {
				// key has expired
				delete(store, k)
				expiredCount++

			}
			if limit == 0 {
				break
			}
		}
	}

	return float32(expiredCount) / float32(20.0)
}

func DeleteExpiredKeys() {
	for {
		frac := deleteExpiredSample()

		if frac < 0.25 {
			break
		}
	}
	log.Println("deleted the expired but undeleted keys. total keys", len(store))

}