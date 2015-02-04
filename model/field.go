package model

import . "github.com/cosiner/golib/types"

type (
	// FieldSet represend a model's field set
	FieldSet LightBitSet
	// Field represend a field of model
	Field uint
)

// NewFieldSet create a new field set use given fieldcount as length,
// initial with given fields, every field will be valid by the given valid
// function
func NewFieldSet(fields ...Field) *FieldSet {
	bs := NewLightBitSet()
	for _, f := range fields {
		bs.Set(f.UNum())
	}
	return (*FieldSet)(bs)
}

// HasField check whether or not fieldset has given field
func (fs *FieldSet) HasField(field Field) bool {
	return (*LightBitSet)(fs).IsSet(field.UNum())
}

// AddField add an field to fieldset
func (fs *FieldSet) AddField(field Field) {
	(*LightBitSet)(fs).Set(field.UNum())
}

// RemoveField remove field from fieldset
func (fs *FieldSet) RemoveField(field Field) {
	(*LightBitSet)(fs).Unset(field.UNum())
}

// Uint return unit disblay of fieldset
func (fs *FieldSet) Uint() uint {
	return (*LightBitSet)(fs).Uint()
}

// ChangeField add field if has is true, alse remove field
func (fs *FieldSet) ChangeField(field Field, has bool) {
	(*LightBitSet)(fs).SetTo(field.UNum(), has)
}

func (fs *FieldSet) FlipAll() {
	(*LightBitSet)(fs).FlipAll()
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
