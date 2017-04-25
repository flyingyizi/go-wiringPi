package gpio

import "container/heap"

//ref https://github.com/brian-armstrong/gpio

/*

//ref https://github.com/hugozhu/rpi/blob/master/rpi.go


*/

type watcherAction int
type fdHeap []uintptr

type watcherCmd struct {
	pin    Pin
	action watcherAction
}
type watcherNotify struct {
	pin   Pin
	value uint
}

// Janitor provides asynchronous notifications on input changes
// The user should supply it pins to watch with AddPin and then
//wait for changes
type Janitor struct {
	pins       map[uintptr]Pin
	fds        fdHeap
	cmdChan    chan watcherCmd
	notifyChan chan watcherNotify
}

// NewJanitor creates a new Watcher instance for asynchronous inputs
func NewJanitor() *Janitor {
	w := &Watcher{
		pins:       make(map[uintptr]Pin),
		fds:        fdHeap{},
		cmdChan:    make(chan watcherCmd, watcherCmdChanLen),
		notifyChan: make(chan watcherNotify, notifyChanLen),
	}
	heap.Init(&w.fds)
	go w.watch()
	return w
}

// AddPin adds a new pin to be watched for changes
// The pin provided should be the pin known by the kernel
func (w *Janitor) AddPin(pin Pin) {
	setEdgeTrigger(pin, edgeBoth)
	w.cmdChan <- watcherCmd{
		pin:    pin,
		action: watcherAdd,
	}
}
