package kit

import (
	"log"
	"sync"
)

func Run(router *Router, transports ...Transport) {
	wg := &sync.WaitGroup{}

	for i := range transports {
		s := transports[i]
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := s.Serve(router); err != nil {
				log.Println(err)
			}
		}()
	}

	wg.Wait()
}
