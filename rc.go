package routineContainers

import "sync"

// RoutineContainers 协程任务池
type RoutineContainers struct {
	Contains []func()
	Wg       sync.WaitGroup
}

// Put 协程任务
func (r *RoutineContainers) Put(f func()) {
	tf := func() {
		f()
		r.Wg.Done()
	}
	r.Contains = append(r.Contains, tf)
	r.Wg.Add(1)
}

// Run 协程任务
func (r *RoutineContainers) Run() {
	for idx := range r.Contains {
		go r.Contains[idx]()
	}
	r.Wg.Wait()
	r.Contains = r.Contains[:0]
}

// 我也新加了一行
