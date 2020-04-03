package gcurl

import (
	"strconv"
	"testing"
)

func TestPQueue(t *testing.T) {
	PQExec := newPQueueExecute()
	PQExec.Push(&parseFunction{Priority: 5})
	PQExec.Push(&parseFunction{Priority: 10})
	PQExec.Push(&parseFunction{Priority: 4})
	PQExec.Push(&parseFunction{Priority: 4})
	PQExec.Push(&parseFunction{Priority: 20})
	PQExec.Push(&parseFunction{Priority: 10})
	PQExec.Push(&parseFunction{Priority: 15})

	content := ""
	for PQExec.Len() > 0 {
		content += strconv.Itoa(PQExec.Pop().Priority)
		content += " "
	}
	if content != "4 4 5 10 10 15 20 " {
		t.Error(content)
	}
}
