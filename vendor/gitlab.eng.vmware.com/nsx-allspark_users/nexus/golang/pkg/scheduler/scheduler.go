package scheduler

import (
	"sync"
	"sync/atomic"
)

type Queue struct {
	channels     [](chan struct{})
	mutex        sync.Mutex
	waitingTasks uint32
}

func (q *Queue) Wait() chan struct{} {
	q.mutex.Lock()
	self := make(chan struct{}, 1)
	q.channels = append(q.channels, self)
	atomic.AddUint32(&q.waitingTasks, 1)
	if len(q.channels) == 1 {
		q.mutex.Unlock()
		// don't we need to remove from channels
		return self
	}
	var channel chan struct{}
	// queue vs array pop
	channel, q.channels = q.channels[0], q.channels[1:]
	q.mutex.Unlock()
	<-channel
	return self
}
func (q *Queue) Empty() bool {
	return atomic.LoadUint32(&q.waitingTasks) == 0
}
func (q *Queue) Done(self chan struct{}) {
	self <- struct{}{}
	q.mutex.Lock()
	atomic.AddUint32(&q.waitingTasks, ^uint32(0))
	q.mutex.Unlock()
}

// ----------------------
type SchedulerHandler struct {
	QDt chan struct{}
	Q   *Queue
	Key string
}
type Scheduler struct {
	ds    map[string]*Queue
	mutex sync.Mutex
}

func NewScheduler() *Scheduler {
	e := &Scheduler{
		ds: make(map[string]*Queue)}
	return e
}
func (s *Scheduler) Wait(key string) SchedulerHandler {
	s.mutex.Lock()
	q, ok := s.ds[key]
	if !ok {
		q = &Queue{}
		s.ds[key] = q
	}
	s.mutex.Unlock()
	return SchedulerHandler{
		QDt: q.Wait(),
		Q:   q,
		Key: key}
}

func (s *Scheduler) Done(self SchedulerHandler) {
	key := self.Key
	self.Q.Done(self.QDt)
	s.mutex.Lock()
	if self.Q.Empty() {
		delete(s.ds, key)
	}
	s.mutex.Unlock()
}

/*
var queue = &Queue{}

func print(i int) {
	self := queue.Wait(i)
	defer queue.Done(i, self)
	defer fmt.Println("finishing i:", i)
	time.Sleep(time.Second)
	fmt.Println("executing i:", i)
	if rand.Int()%2 == 0 {
		panic(fmt.Sprintf("panicking i: %d", i))
	}
}
func main() {
	rand.Seed(time.Now().Unix())
	waitGroup := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		num := i
		waitGroup.Add(1)
		go func() {
			print(num)
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}
*/
