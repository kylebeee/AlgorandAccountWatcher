package server

import (
	"errors"
	"fmt"
	"sync"
)

// generalize our watchlist so in the future we can use this
// for addresses, applications and assets
type ListType interface {
	string | uint64
}

// generalized function that will take in whatever kind of list we setup
type Processor[L ListType] func(item L) error

type WatchList[L ListType] struct {
	Subs       *sync.Map
	Processors []Processor[L]
	ListChan   chan []L
}

type QueueFailures[L ListType] struct {
	*sync.Mutex
	List map[L]error
}

func (wl *WatchList[L]) AddToQueue(list []L) error {
	select {
	case wl.ListChan <- list:
		return nil
	default:
		return errors.New("failed to add list to queue channel")
	}
}

func (wl *WatchList[L]) Start() {
	for list := range wl.ListChan {
		for _, item := range list {
			for _, processor := range wl.Processors {
				err := processor(item)
				if err != nil {
					fmt.Println("[ERROR] processor threw an error: ", err)
					continue
				}
			}
		}
	}
}
