package mdruntime

import (
	"errors"
	"runtime"

	"github.com/sentadmedia/elf/fw"
)

var _ fw.ProgramRuntime = (*BuildIn)(nil)

type BuildIn struct {
}

func (b BuildIn) LockOSThread() {
	runtime.LockOSThread()
}

func (b BuildIn) Caller(numLevelsUp int) (fw.Caller, error) {
	_, file, line, ok := runtime.Caller(numLevelsUp + 1)
	if !ok {
		return fw.Caller{}, errors.New("fail to obtain caller info")
	}
	return fw.Caller{
		FullFilename: file,
		LineNumber:   line,
	}, nil
}

func NewBuildIn() BuildIn {
	return BuildIn{}
}
