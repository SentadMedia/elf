package mdtracer

import (
	"fmt"
	"time"

	"github.com/sentadmedia/elf/fw"

	uuid "github.com/satori/go.uuid"
)

var _ fw.Tracer = (*Local)(nil)

type Local struct {
	logger fw.Logger
}

type LocalTrace struct {
	id     string
	name   string
	start  time.Time
	logger fw.Logger
}

func (t LocalTrace) Next(name string) fw.Segment {
	start := time.Now()
	t.logger.Trace(fmt.Sprintf("[Trace Start id=%s name=%s startAt=%v]", t.id, name, start))
	// log.Printf("[Trace Start id=%s name=%s startAt=%v]", t.id, name, start)
	return LocalTrace{
		id:     t.id,
		name:   name,
		start:  start,
		logger: t.logger,
	}
}

func (t LocalTrace) End() {
	end := time.Now()
	diff := end.Sub(t.start)
	t.logger.Trace(fmt.Sprintf("[Trace End   id=%s name=%s endAt=%v duration=%v]", t.id, t.name, end, diff))
	// log.Printf()
}

func (l Local) BeginTrace(name string) fw.Segment {
	id := uuid.NewV4().String()
	start := time.Now()
	l.logger.Trace(fmt.Sprintf("[Trace Start id=%s name=%s startAt=%v]", id, name, start))
	// log.Printf()
	return LocalTrace{
		id:     id,
		name:   name,
		start:  start,
		logger: l.logger,
	}
}

func NewLocal(logger fw.Logger) fw.Tracer {
	return Local{logger: logger}
}
