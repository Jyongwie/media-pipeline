package worker

import (
	"context"
	"fmt"
	"time"

	// REPLACE YOUR_MODULE_NAME with your actual go.mod name!
	"github.com/Jyongwie/media-pipeline/backend/internal/infrastructure"
)

// StartRenderPool spins up a background worker that checks for jobs
func StartRenderPool(repo *infrastructure.Repository, broadcast func()) {
	// The 'go' keyword spins this anonymous function off into its own background thread
	go func() {
		fmt.Println("⚙️  Background Render Worker spinning up...")
		
		// A Ticker fires an event every 3 seconds
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			// This blocks the loop until the 3-second ticker ticks
			<-ticker.C 

			ctx := context.Background()
			
			// 1. Attempt to grab a locked job
			job, err := repo.GetNextQueuedJob(ctx)
			if err != nil {
				fmt.Printf("Worker error fetching job: %v\n", err)
				continue
			}
			
			// If queue is empty, silently wait for the next tick
			if job == nil {
				continue 
			}

			fmt.Printf("\n👷 Worker picked up Job [%s] - Type: %s\n", job.ID, job.AssetType)
			fmt.Println("⏳ Simulating intense GPU rendering...")

			// 2. Simulate the heavy lifting (Rendering a 4K video takes time!)
			time.Sleep(8 * time.Second)

			// 3. Mark as finished
			err = repo.MarkJobCompleted(ctx, job.ID)
			if err != nil {
				fmt.Printf("Failed to mark job completed: %v\n", err)
				continue
			}

			fmt.Printf("✅ Job [%s] successfully completed!\n\n", job.ID)

			broadcast()
		}
	}()
}