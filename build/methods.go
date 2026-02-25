package build

import (
	"fmt"

	"github.com/pterm/pterm"
)

func (i *infoMessages) appendMsg(format string, args ...interface{}) {
	i.msgs = append(i.msgs, infoMessage{format: format, args: args})
}

func (i infoMessage) print() {
	pterm.Info.Printfln(i.format, i.args...)
}

func (e *errMsg) append(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.err != nil {
		e.err = fmt.Errorf("%w; %w", e.err, err)
		return
	}
	e.err = err
}
