package model

import "github.com/cosiner/golib/types"

type (
	FieldSet struct {
		*types.BitSet
	}
	Field uint
)

func NewFieldSet(fieldCount uint, fields ...Field) FieldSet {
	bs := types.NewBitSet(fieldCount)
	for _, f := range fields {
		bs.Set(uint(f))
	}
	return FieldSet{bs}
}

func (fs FieldSet) HasField(field Field) bool {
	return fs.IsSet(uint(field))
}

func (fs FieldSet) AddField(field Field) {
	fs.Set(uint(field))
}

func (fs FieldSet) ChangeField(field Field, has bool) {
	fs.SetTo(uint(field), has)
}

func (fs FieldSet) RemoveField(field Field) {
	fs.UnSet(uint(field))
}

func (fs FieldSet) FieldCount() uint {
	return fs.BitCount()
}

func (f Field) Num() int {
	return int(f)
}

func (f Field) UNum() uint {
	return uint(f)
}

func (f Field) Equal(field Field) bool {
	return f.UNum() == field.UNum()
}
