package buffer

import (
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
		timeout: 1000 * time.Millisecond,
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
	if m.pending == cmd {
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

		case *DeleteTextCommand:
			pendingcmd, ok := m.pending.(*DeleteTextCommand)
			if ok {
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
			pendingcmd.Text = curr.Text + pendingcmd.Text
			pendingcmd.End = curr.End

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
