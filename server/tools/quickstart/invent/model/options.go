package model

import (
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent/field"
)

// SchemaNameOpt schema name option
type SchemaNameOpt func(*gorm.DB) string

var defaultSchemaNameOpt = SchemaNameOpt(func(db *gorm.DB) string { return db.Migrator().CurrentDatabase() })

// Option field option
type Option interface{ OptionType() string }

const fieldType = "field"

// FieldOption ...
type FieldOption interface {
	Option
	Operator() func(*field.Field) *field.Field
}

const methodType = "method"

// MethodOption ...
type MethodOption interface {
	Option
	Methods() (methods []interface{})
}

var (
	_ Option = ModifyFieldOpt(nil)
	_ Option = FilterFieldOpt(nil)
	_ Option = CreateFieldOpt(nil)

	_ Option = AddMethodOpt(nil)
)

// ModifyFieldOpt modify field option
type ModifyFieldOpt func(*field.Field) *field.Field

// OptionType implement for interface Option
func (ModifyFieldOpt) OptionType() string { return fieldType }

// Operator implement for FieldOpt
func (o ModifyFieldOpt) Operator() func(*field.Field) *field.Field { return o }

// FilterFieldOpt filter field option
type FilterFieldOpt ModifyFieldOpt

// OptionType implement for interface Option
func (FilterFieldOpt) OptionType() string { return fieldType }

// Operator implement for FieldOpt
func (o FilterFieldOpt) Operator() func(*field.Field) *field.Field { return o }

// CreateFieldOpt create field option
type CreateFieldOpt ModifyFieldOpt

// OptionType implement for interface Option
func (CreateFieldOpt) OptionType() string { return fieldType }

// Operator implement for FieldOpt
func (o CreateFieldOpt) Operator() func(*field.Field) *field.Field { return o }

// AddMethodOpt diy method option
type AddMethodOpt func() (methods []interface{})

// OptionType implement for interface Option
func (AddMethodOpt) OptionType() string { return methodType }

// Methods ...
func (o AddMethodOpt) Methods() []interface{} { return o() }

func sortOptions(opts []Option) (modifyOpts []FieldOption, filterOpts []FieldOption, createOpts []FieldOption, methodOpt []MethodOption) {
	for _, opt := range opts {
		switch opt := opt.(type) {
		case ModifyFieldOpt:
			modifyOpts = append(modifyOpts, opt)
		case FilterFieldOpt:
			filterOpts = append(filterOpts, opt)
		case CreateFieldOpt:
			createOpts = append(createOpts, opt)
		case AddMethodOpt:
			methodOpt = append(methodOpt, opt)
		}
	}
	return
}
