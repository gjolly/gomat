package Tasks

import (
	"sync"
	"github.com/matei13/gomat/Daemon/gomatcore"
	"net"
)

type TaskMap struct {
	Tasks map[string][]Task
	Lock  *sync.RWMutex
}

type Task struct {
	ID     uint32
	Origin net.UDPAddr
	Mat1   gomatcore.SubMatrix
	Mat2   gomatcore.SubMatrix
}

func (tm *TaskMap) GetTasks(p string) []Task {
	tm.Lock.Lock()
	defer tm.Lock.Unlock()
	return tm.Tasks[p]
}

func (tm *TaskMap) AddTask(peer string, t Task) {
	tm.Lock.Lock()
	defer tm.Lock.Unlock()
	if _, ok := tm.Tasks[peer]; !ok {
		tm.Tasks[peer] = make([]Task, 1)
	}
	tm.Tasks[peer] = append(tm.Tasks[peer], t)
}

func (tm *TaskMap) RemoveTask(peer string, id uint32) {
	tm.Lock.Lock()
	defer tm.Lock.Unlock()
	if l, ok := tm.Tasks[peer]; ok {
		for n, t := range l {
			if t.ID == id {
				l = append(l[:n], l[n+1:]...)
				return
			}
		}
	}
}

func (t Task) Size() int {
	if t.Mat1.MaxDim() > t.Mat2.MaxDim() {
		return t.Mat1.MaxDim()
	}
	return t.Mat2.MaxDim()
}
