package jobs

import (
	"log"
	"time"

	"debitask/store"
)

func StartHabitMissChecker(interval time.Duration) {
	go func() {
		runHabitCheck()
		ticker := time.NewTicker(interval)
		for range ticker.C {
			runHabitCheck()
		}
	}()
}

func runHabitCheck() {
	yesterday := time.Now().UTC().Truncate(24 * time.Hour).Add(-24 * time.Hour)
	count, err := store.MarkMissedHabits(yesterday)
	if err != nil {
		log.Printf("habit miss checker error: %v\n", err)
		return
	}
	log.Printf("habit miss checker: marked %d missed habit(s) for %s\n", count, yesterday.Format("2006-01-02"))
}
