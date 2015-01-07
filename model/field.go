package model

import "github.com/cosiner/golib/types"

type (
	// FieldSet represend a model's field set
	FieldSet struct {
		bs *types.BitSet
	}
	// Field represend a field of model
	Field uint
)

// NewFieldSet create a new field set use given fieldcount as length,
// initial with given fields
func NewFieldSet(fieldCount uint, fields ...Field) FieldSet {
	bs := types.NewBitSet(fieldCount)
	for _, f := range fields {
		bs.Set(uint(f))
	}
	return FieldSet{bs}
}

// HasField check whether or not fieldset has given field
func (fs FieldSet) HasField(field Field) bool {
	return fs.bs.IsSet(uint(field))
}

// AddField add an field to fieldset
func (fs FieldSet) AddField(field Field) {
	fs.bs.Set(uint(field))
}

// RemoveField remove field from fieldset
func (fs FieldSet) RemoveField(field Field) {
	fs.bs.UnSet(uint(field))
}

// ChangeField add field if has is true, alse remove field
func (fs FieldSet) ChangeField(field Field, has bool) {
	fs.bs.SetTo(uint(field), has)
}

// FieldCount return all fields' count exist in fieldset
// it differente to 'Len' method
func (fs FieldSet) FieldCount() uint {
	return fs.bs.BitCount()
}

// Len return fieldset's length
func (fs FieldSet) Len() uint {
	return fs.bs.Len()
}

// Num display field as number
func (f Field) Num() int {
	return int(f)
}

// UNum display field as unsigned number
func (f Field) UNum() uint {
	return uint(f)
}

// Equal check whether two field is equal by field's display number
func (f Field) Equal(field Field) bool {
	return f.UNum() == field.UNum()
}
