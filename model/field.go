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
// initial with given fields, every field will be valid by the given valid
// function
func NewFieldSet(fieldCount uint, valid func(field Field), fields ...Field) FieldSet {
	bs := types.NewBitSet(fieldCount)
	if valid == nil {
		valid = func(Field) {}
	}
	for _, f := range fields {
		valid(f)
		bs.Set(f.UNum())
	}
	return FieldSet{bs}
}

// HasField check whether or not fieldset has given field
func (fs FieldSet) HasField(field Field) bool {
	return fs.bs.IsSet(field.UNum())
}

// AddField add an field to fieldset
func (fs FieldSet) AddField(field Field) {
	fs.bs.Set(field.UNum())
}

// RemoveField remove field from fieldset
func (fs FieldSet) RemoveField(field Field) {
	fs.bs.UnSet(field.UNum())
}

// ChangeField add field if has is true, alse remove field
func (fs FieldSet) ChangeField(field Field, has bool) {
	fs.bs.SetTo(field.UNum(), has)
}

func (fs FieldSet) FlipAll() {
	fs.bs.FlipAll()
}

// Fields return all fields in this field set
func (fs FieldSet) Fields() (res []Field) {
	bs := fs.bs
	res = make([]Field, 0, bs.BitCount())
	for i := types.Uint0; i < fs.Len(); i++ {
		if bs.IsSet(i) {
			res = append(res, NewField(uint(i)))
		}
	}
	return
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

// NewField create a new field from a number
func NewField(num uint) Field {
	return Field(num)
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
