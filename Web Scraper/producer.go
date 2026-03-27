package main

func produceIDs(idChan chan int, start int, end int) {
	for id := start; id <= end; id++ {
		idChan <- id
	}
	close(idChan)
}
