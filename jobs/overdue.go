package jobs

import (
	"log"
	"time"

	"debitask/store"
)

func StartOverdueChecker(interval time.Duration) {
	go func() {
		runCheck()
		ticker := time.NewTicker(interval)
		for range ticker.C {
			runCheck()
		}
	}()
}

func runCheck() {
	count, err := store.MarkOverdueTasks()
	if err != nil {
		log.Printf("overdue checker error: %v\n", err)
		return
	}
	log.Printf("overdue checker: marked %d task(s) as overdue\n", count)
}
