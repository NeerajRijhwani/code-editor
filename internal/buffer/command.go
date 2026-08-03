package buffer

import "log"

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
	Start  Position
	length int
}

type SplitLineCommand struct {
	Start Position
}

func Init_InsertTextCommand(x1, y1, x2, y2 int, text string) *InsertTextCommand {
	return &InsertTextCommand{
		Start: Position{
			Row: x1,
			Col: y1,
		},
		End: Position{
			Row: x2,
			Col: y2,
		},
		Text: text,
	}
}
func Init_DeleteTextCommand(x1, y1, x2, y2 int, text string) *DeleteTextCommand {
	return &DeleteTextCommand{
		Start: Position{
			Row: x1,
			Col: y1,
		},
		End: Position{
			Row: x2,
			Col: y2,
		},
		Text: text,
	}
}

func Init_MergeLineCommand(x, y, l int) *MergeLineCommand {
	return &MergeLineCommand{
		Start: Position{
			Row: x,
			Col: y,
		},
		length: l,
	}
}
func Init_SplitLineCommand(x, y int) *SplitLineCommand {
	return &SplitLineCommand{
		Start: Position{
			Row: x,
			Col: y,
		},
	}
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
	log.Printf("Undo fo DeleteTextCommand called x = %d, y = %d", cmd.Start.Row, cmd.Start.Col)
	err := b.InsertText(cmd.Text, cmd.Start.Row, cmd.Start.Col)
	if err != nil {
		log.Printf(" Error : %v", err)
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
	err := b.SplitLine(cmd.Start.Row-1, cmd.length)
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
	err := b.MergeLine(cmd.Start.Row + 1)
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
