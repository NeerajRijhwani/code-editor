package buffer

import (
	"log"
	"reflect"
	"time"

	"github.com/NeerajRijhwani/code-editor/internal/utils"
)

type Manager struct {
	undo *utils.Stack[Command]
	redo *utils.Stack[Command]

	pending Command

	last time.Time

	timeout time.Duration
}

func Init_Manager() *Manager {
	return &Manager{
		undo:    utils.Init_Stack[Command](),
		redo:    utils.Init_Stack[Command](),
		pending: nil,
		last:    time.Now(),
		timeout: 5000 * time.Millisecond,
	}
}

func (m *Manager) Execute(cmd Command, b *Buffer, ch rune) {
	switch ch {
	case 'u':
		m.executeUndo(b)
	case 'r':
		m.executeRedo(b)

	default:
		if !m.canMerge(cmd) {
			log.Printf("Merging failed pending %v and current %v  ", m.pending, cmd)
			m.commitPending()
			m.pending = cmd
			m.redo.Clear()
		}
	}
	m.last = time.Now()
}

func (m *Manager) executeUndo(b *Buffer) error {
	m.commitPending()
	cmd, err := m.undo.Pop()
	log.Printf("Undo Command %v", cmd)
	if err != nil {
		return err
	}

	if err := cmd.Undo(b); err != nil {
		return err
	}

	m.redo.Push(cmd)
	return nil
}

func (m *Manager) executeRedo(b *Buffer) error {
	m.commitPending()
	cmd, err := m.redo.Pop()
	if err != nil {
		return err
	}

	if err := cmd.Redo(b); err != nil {
		return err
	}

	m.undo.Push(cmd)
	return nil
}

func (m *Manager) commitPending() {
	if m.pending == nil {
		return
	}
	m.undo.Push(m.pending)
	m.pending = nil
}

func (m *Manager) canMerge(cmd Command) bool {
	p := reflect.TypeOf(m.pending)
	c := reflect.TypeOf(cmd)
	if p == c {
		switch curr := cmd.(type) {
		case *InsertTextCommand:
			pendingcmd, ok := m.pending.(*InsertTextCommand)
			if !ok {
				return false
			}
			x1 := pendingcmd.Start.Row
			y1 := pendingcmd.Start.Col
			// adjacent request check
			if x1 != curr.Start.Row && y1+1 != curr.Start.Col {
				return false
			}
			t := time.Now()
			// concurrent request check based on time
			if t.Sub(m.last) > m.timeout {
				return false
			}
			pendingcmd.Text += curr.Text
			pendingcmd.End = curr.End
			return true

		case *DeleteTextCommand:
			pendingcmd, ok := m.pending.(*DeleteTextCommand)
			if !ok {
				return false
			}
			x1 := pendingcmd.Start.Row
			y1 := pendingcmd.Start.Col
			// adjacent request check
			log.Printf("x1 %d and y1 %d and currx %d and curry %d", x1, y1, curr.End.Row, curr.End.Col)
			if x1 != curr.End.Row && y1 != curr.End.Col {
				return false
			}
			t := time.Now()
			// concurrent request check based on time
			if t.Sub(m.last) > m.timeout {
				log.Printf("Timeout for merging")
				return false
			}
			pendingcmd.Text = curr.Text + pendingcmd.Text
			pendingcmd.Start = curr.Start
			return true

		case *SplitLineCommand:
			return false
		case *MergeLineCommand:
			return false
		default:
			return false
		}
	}
	return false
}
