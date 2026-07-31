package buffer

import (
	"github.com/NeerajRijhwani/code-editor/internal/utils"
	"time"
)

type Manager struct {
	undo utils.Stack[Command]
	redo utils.Stack[Command]

	pending Command

	last time.Time
}

func (m *Manager) Undo(b *Buffer) error {
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

func (m *Manager) Redo(b *Buffer) error {
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
