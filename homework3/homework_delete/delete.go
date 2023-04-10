package homework_delete

import (
	"errors"
	"reflect"
	"strings"
)

type Deleter[T any] struct {
	table     string
	sb        *strings.Builder
	args      []any
	predicate []Predicate
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.sb = &strings.Builder{}
	//写入DELETE * FROM
	d.sb.WriteString("DELETE FROM ")
	var t T
	//写入表名
	if d.table == "" {
		d.sb.WriteByte('`')
		d.sb.WriteString(reflect.TypeOf(t).Name())
		d.sb.WriteByte('`')
	} else {
		d.sb.WriteString(d.table)
	}
	if len(d.predicate) > 0 {
		d.sb.WriteString(" WHERE ")
		p := d.predicate[0]
		for i := 1; i < len(d.predicate); i++ {
			p = p.And(d.predicate[i])
		}
		if err := d.BuildExpression(p); err != nil {
			return nil, err
		}
	}
	d.sb.WriteByte(';')
	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

//func BuildWhere(expr Expression, sb *strings.Builder, args []any) error {
//	switch e := expr.(type) {
//	case nil:
//		return nil
//	case Column:
//		sb.WriteByte('`')
//		sb.WriteString(e.name)
//		sb.WriteByte('`')
//	case value:
//		sb.WriteByte('?')
//		args = append(args, e.val)
//	case Predicate:
//		_, l := e.left.(Predicate)
//		if l {
//			sb.WriteByte('(')
//			BuildWhere(e.left, sb, args)
//			sb.WriteByte(')')
//		} else {
//			BuildWhere(e.left, sb, args)
//		}
//		sb.WriteByte(' ')
//		sb.WriteString(e.op.String())
//		sb.WriteByte(' ')
//		_, r := e.right.(Predicate)
//		if r {
//			sb.WriteByte('(')
//			BuildWhere(e.right, sb, args)
//			sb.WriteByte(')')
//		} else {
//			BuildWhere(e.right, sb, args)
//		}
//	default:
//		return errors.New("orm:不支持的Predicate")
//	}
//	return nil
//}

func (d *Deleter[T]) BuildExpression(expr Expression) error {
	switch e := expr.(type) {
	case nil:
		return nil
	case Column:
		d.sb.WriteByte('`')
		d.sb.WriteString(e.name)
		d.sb.WriteByte('`')
	case value:
		d.sb.WriteByte('?')
		d.args = append(d.args, e.val)
	case Predicate:
		_, l := e.left.(Predicate)
		if l {
			d.sb.WriteByte('(')
			d.BuildExpression(e.left)
			d.sb.WriteByte(')')
		} else {
			d.BuildExpression(e.left)
		}
		d.sb.WriteByte(' ')
		d.sb.WriteString(e.op.String())
		d.sb.WriteByte(' ')
		_, r := e.right.(Predicate)
		if r {
			d.sb.WriteByte('(')
			d.BuildExpression(e.right)
			d.sb.WriteByte(')')
		} else {
			d.BuildExpression(e.right)
		}
	default:
		return errors.New("orm:不支持的Predicate")
	}
	return nil
}

// From accepts model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.predicate = predicates
	return d
}
