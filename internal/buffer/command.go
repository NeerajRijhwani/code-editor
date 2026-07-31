package buffer

type Command interface {
	Undo(*Buffer) error
	Redo(*Buffer) error
}

type Position struct {
	Row int
	Col int
}

type InsertTextCommand struct {
	Start Position
	End   Position
	Text  string
}

type DeleteTextCommand struct {
	Start Position
	End   Position
	Text  string
}

type MergeLineCommand struct {
	Start Position
}

type SplitLineCommand struct {
	Start Position
}

func (cmd *InsertTextCommand) Undo(b *Buffer) error {

	x1 := cmd.Start.Row
	y1 := cmd.Start.Col
	x2 := cmd.End.Row
	y2 := cmd.End.Col

	b.DeleteText(x1, y1, x2, y2)

	return nil
}

func (cmd *InsertTextCommand) Redo(b *Buffer) error {
	err := b.InsertText(cmd.Text, cmd.Start.Row, cmd.Start.Col)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *DeleteTextCommand) Undo(b *Buffer) error {
	err := b.InsertText(cmd.Text, cmd.Start.Row, cmd.Start.Col)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *DeleteTextCommand) Redo(b *Buffer) error {
	x1 := cmd.Start.Row
	y1 := cmd.Start.Col
	x2 := cmd.End.Row
	y2 := cmd.End.Col

	b.DeleteText(x1, y1, x2, y2)

	return nil
}

func (cmd *MergeLineCommand) Undo(b *Buffer) error {
	err := b.SplitLine(cmd.Start.Row, cmd.Start.Col)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *MergeLineCommand) Redo(b *Buffer) error {
	err := b.MergeLine(cmd.Start.Row)
	if err != nil {
		return err
	}
	return nil
}
func (cmd *SplitLineCommand) Undo(b *Buffer) error {
	err := b.MergeLine(cmd.Start.Row)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *SplitLineCommand) Redo(b *Buffer) error {
	err := b.SplitLine(cmd.Start.Row, cmd.Start.Col)
	if err != nil {
		return err
	}
	return nil
}
