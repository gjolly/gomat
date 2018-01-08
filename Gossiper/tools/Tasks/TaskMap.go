package Tasks

import (
	"sync"
	"github.com/matei13/gomat/Daemon/gomatcore"
)

type TaskMap struct {
	Tasks map[string][]Task
	Lock  *sync.RWMutex
}

type Task struct {
	ID     uint32
	Origin string
	Mat1   gomatcore.Matrix
	Mat2   gomatcore.Matrix
}

func (tm TaskMap) getTasks(p string) []Task {
	tm.Lock.Lock()
	defer tm.Lock.Unlock()
	return tm.Tasks[p]
}

func (tm TaskMap) addTask(peer string, origin string, mat1 gomatcore.Matrix, mat2 gomatcore.Matrix) (id uint32) {
	tm.Lock.Lock()
	defer tm.Lock.Unlock()
	if _, ok := tm.Tasks[peer]; !ok {
		tm.Tasks[peer] = make([]Task, 1)
	}
	id = uint32(0)

	for _, t := range tm.Tasks[peer] {
		if t.ID > id {
			id = t.ID
		}
	}
	if id > 0 {
		id++
	}
	tm.Tasks[peer] = append(tm.Tasks[peer],
		Task{ID: id,
			Origin: origin,
			Mat1: mat1,
			Mat2: mat2,})
	return
}

func (tm TaskMap) removeTask(peer string, id uint32) {
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
