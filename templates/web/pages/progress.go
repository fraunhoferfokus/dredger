// Edit this file, as it is a specific handler function for your service
package rest

import (
	"fmt"
	"time"

	"github.com/r3labs/sse/v2"
)

// learn medium duration
var durationSum int = 100 // 10 Seconds as initial duration
var durationNb int = 1

const picoNull = `<progress value=\"0\"max=\"100\" />`
const picoEndless = `<progress />`
const picoProgress = `<progress value=\"%d\"max=\"100\" />`
const picoEmpty = `<div/>`

func ProgressPico(f func()) {
	duration := durationSum / durationNb
	c1 := make(chan bool)
	go func() {
		f()
		c1 <- true
	}()

	v := 0
	SseServer.Publish("progress", &sse.Event{
		Event: []byte("Progress"),
		Data:  []byte(picoNull),
	})
	for ready := false; !ready; {
		select {
		case <-c1:
			ready = true
		case <-time.After(500 * time.Millisecond):
			v = v + 5
			if v > duration && v <= duration+5 {
				SseServer.Publish("progress", &sse.Event{
					Event: []byte("Progress"),
					Data:  []byte(picoEndless),
				})
			} else if v < duration {
				SseServer.Publish("progress", &sse.Event{
					Event: []byte("Progress"),
					Data:  []byte(fmt.Sprintf(picoProgress, v)),
				})
			}
		}
	}
	SseServer.Publish("progress", &sse.Event{
		Event: []byte("Progress"),
		Data:  []byte(picoEmpty),
	})

	// update medium duration
	durationSum += v
	durationNb += 1
	// avoid overflow of durationNb
	if durationNb > 10000 {
		durationSum = durationSum / durationNb
		durationNb = 1
	}
}

const bsNull = `<div class="progress" role="progressbar" aria-label="Fortschritt " aria-valuenow="0" aria-valuemin="0" aria-valuemax="100">
  <div class="progress-bar" style="width: 0%"></div>
</div>`
const bsEndless = `<div class="progress" role="progressbar" aria-label="Fortschritt" aria-valuenow="100" aria-valuemin="0" aria-valuemax="100">
<div class="progress-bar" style="width: 100%"></div>
</div>`
const bsProgress = `<div class="progress" role="progressbar" aria-label="Fortschritt" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">
<div class="progress-bar" style="width: %d%%"></div>
</div>`
const bsEmpty = `<div/>`

func ProgressBootstrap(f func()) {
	duration := durationSum / durationNb
	c1 := make(chan bool)
	go func() {
		f()
		c1 <- true
	}()

	v := 0
	SseServer.Publish("progress", &sse.Event{
		Event: []byte("Progress"),
		Data:  []byte(bsNull),
	})
	for ready := false; !ready; {
		select {
		case <-c1:
			ready = true
		case <-time.After(500 * time.Millisecond):
			v = v + 5
			if v > duration {
				SseServer.Publish("progress", &sse.Event{
					Event: []byte("Progress"),
					Data:  []byte(bsEndless),
				})
			} else {
				SseServer.Publish("progress", &sse.Event{
					Event: []byte("Progress"),
					Data:  []byte(fmt.Sprintf(bsProgress, v, v)),
				})
			}
		}
	}
	SseServer.Publish("progress", &sse.Event{
		Event: []byte("Progress"),
		Data:  []byte(bsEmpty),
	})

	// update medium duration
	durationSum += v
	durationNb += 1
	// avoid overflow of durationNb
	if durationNb > 10000 {
		durationSum = durationSum / durationNb
		durationNb = 1
	}
}
