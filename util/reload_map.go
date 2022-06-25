package util

import (
	"encoding/json"
	"errors"
)

const (
	// AsyncMapActionOverwrite ...
	AsyncMapActionOverwrite = 1
	// AsyncMapActionDelete ...
	AsyncMapActionDelete = 2
)

// AsyncMapOperand ...
type AsyncMapOperand struct {
	Key    string
	Value  interface{}
	Action int
	Done   chan bool
}

// AsyncMapInit ...
func AsyncMapInit() (chan *AsyncMapOperand, map[string]interface{}) {
	queue := make(chan *AsyncMapOperand, 0)
	container := make(map[string]interface{}, 0)
	return queue, container
}

// AsyncMapStart ...
func AsyncMapStart(queue chan *AsyncMapOperand, container *map[string]interface{}) error {
	if queue == nil {
		return errors.New("operand queue chan is nil")
	}

	if container == nil {
		return errors.New("map container is empty")
	}

	go func() {
		for {
			operand := <-queue
			if operand == nil {
				continue
			}

			for i := 0; i == 0; i++ {

				bytesApps, err := json.Marshal(*container)
				if err != nil {
					break
				}

				rawContainer := make(map[string]interface{}, 0)
				err = json.Unmarshal([]byte(bytesApps), &rawContainer)
				if err != nil {
					break
				}

				switch operand.Action {
				case AsyncMapActionOverwrite:
					rawContainer[operand.Key] = operand.Value
				case AsyncMapActionDelete:
					delete(rawContainer, operand.Key)
				default:
					rawContainer[operand.Key] = operand.Value
				}

				*container = rawContainer
			}
			operand.Done <- true
		}
	}()

	return nil
}

// AsyncMapUpdate Update map k-v safely
func AsyncMapUpdate(queue chan *AsyncMapOperand, operand *AsyncMapOperand) {
	if operand == nil {
		return
	}

	queue <- operand

	<-operand.Done
}

// NewAsyncMapOperand ...
func NewAsyncMapOperand(action int, key string, value interface{}) *AsyncMapOperand {
	return &AsyncMapOperand{
		Key:    key,
		Value:  value,
		Action: action,
		Done:   make(chan bool),
	}
}
