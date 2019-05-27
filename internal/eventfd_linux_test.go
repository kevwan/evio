package internal

import (
	"testing"
)

func TestNew(t *testing.T) {
	efd, err := newEventFd()
	if err != nil {
		t.Error("Could not create EventFd")
		return
	}
	defer efd.Close()

	if efd.Fd() < 0 {
		t.Errorf("invalid FD %d", efd.Fd())
		return
	}
}

func TestReadWriteEvent(t *testing.T) {
	efd, err := newEventFd()
	if err != nil {
		t.Error(err)
	}
	defer efd.Close()

	var good uint64 = 0x78
	if err := efd.WriteEvent(good); err != nil {
		t.Error(err)
	}

	if actual, err := efd.ReadEvent(); err != nil {
		t.Error(err)
	} else if actual != good {
		t.Errorf("error reading from eventfd, expected: %q, actual: %q", good, actual)
	}
}

func BenchmarkReadWriteEvent(b *testing.B) {
	const event = 15
	efd, err := newEventFd()
	if err != nil {
		b.Fatal(err)
	}
	defer efd.Close()

	for i := 0; i < b.N; i++ {
		if err := efd.WriteEvent(event); err != nil {
			b.Fatal(err)
		}
		val, err := efd.ReadEvent()
		if err != nil {
			b.Fatal(err)
		} else if val != event {
			b.Fail()
		}
	}
}
