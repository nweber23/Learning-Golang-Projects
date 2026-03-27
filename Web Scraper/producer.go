package main

// produceIDs generates IDs from start to end (inclusive) and sends them down the channel.
// This runs in its own goroutine.
//
// When done, it closes idChan to signal to workers: "no more IDs coming"
func produceIDs(idChan chan int, start int, end int) {
	// Loop from start to end (inclusive)
	for id := start; id <= end; id++ {
		// Send the ID down the channel
		// This blocks if the channel buffer is full, creating backpressure
		idChan <- id
	}

	// Signal to workers: no more IDs coming
	// Workers read from a closed channel and exit
	close(idChan)
}
