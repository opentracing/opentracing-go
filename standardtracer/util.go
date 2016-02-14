package standardtracer

import (
	"math/rand"
	"sync"
	"time"
)

var (
	seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	// The golang rand generators are *not* intrinsically thread-safe.
	seededIDLock sync.Mutex
)

func randomID() int64 {
	seededIDLock.Lock()
	defer seededIDLock.Unlock()
	return seededIDGen.Int63()
}

func randomID2() (int64, int64) {
	seededIDLock.Lock()
	defer seededIDLock.Unlock()
	return seededIDGen.Int63(), seededIDGen.Int63()
}
