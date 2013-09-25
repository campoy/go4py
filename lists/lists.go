package lists

import (
	"fmt"
	"reflect"
)

func Filter(q interface{}, l interface{}) interface{} {
	if _, ok := l.(error); ok {
		return fmt.Errorf("filter: %v", l)
	}

	lt := reflect.TypeOf(l)
	if lt.Kind() != reflect.Slice {
		return fmt.Errorf("filter: expected slice in 2nd arg, got %T", l)
	}
	if reflect.ValueOf(l).Len() == 0 {
		return l
	}
	et := lt.Elem()

	qt := reflect.TypeOf(q)
	if qt.Kind() != reflect.Func {
		return fmt.Errorf("filter: expected func in 1st arg, got %T", q)
	}
	if qt.NumIn() != 1 {
		return fmt.Errorf("filter: expected func with one parameter, got %v", qt.NumIn())
	}
	if it := qt.In(0); it != et {
		return fmt.Errorf("filter: non macthing types of the slice and the function: %v vs %v", et, qt)
	}
	if qt.NumOut() != 1 {
		return fmt.Errorf("filter: expected func with one result, got %v", qt.NumOut())
	}
	if ot := qt.Out(0); ot.Kind() != reflect.Bool {
		return fmt.Errorf("filter: expected func returning bool, it returns %v", ot)
	}

	lv, qv := reflect.ValueOf(l), reflect.ValueOf(q)

	v := reflect.MakeSlice(lt, 0, 0)
	for i := 0; i < lv.Len(); i++ {
		x := lv.Index(i)
		if qv.Call([]reflect.Value{x.Convert(et)})[0].Interface().(bool) {
			v = reflect.Append(v, x)
		}
	}

	return v.Interface()
}

func Map(f interface{}, l interface{}) interface{} {
	if _, ok := l.(error); ok {
		return fmt.Errorf("map: %v", l)
	}

	lt := reflect.TypeOf(l)
	if lt.Kind() != reflect.Slice {
		return fmt.Errorf("map: expected slice in 2nd arg, got %T", l)
	}
	et := lt.Elem()

	ft := reflect.TypeOf(f)
	if ft.Kind() != reflect.Func {
		return fmt.Errorf("map: expected func in 1st arg, got %T", f)
	}
	if ft.NumIn() != 1 {
		return fmt.Errorf("map: expected func with one parameter, got %v", ft.NumIn())
	}
	if it := ft.In(0); it != et {
		return fmt.Errorf("map: non macthing types of the slice and the function: %v vs %v", et, ft)
	}
	if ft.NumOut() != 1 {
		return fmt.Errorf("map: expected func with one result, got %v", ft.NumOut())
	}
	ot := ft.Out(0)

	lv, fv := reflect.ValueOf(l), reflect.ValueOf(f)

	v := reflect.MakeSlice(reflect.SliceOf(ot), 0, lv.Len())
	for i := 0; i < lv.Len(); i++ {
		x := lv.Index(i)
		y := fv.Call([]reflect.Value{x.Convert(et)})[0]
		v = reflect.Append(v, y)
	}

	return v.Interface()
}
