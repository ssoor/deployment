package main

import (
	"fmt"
	"os"
	"reflect"
)

func buildTarget(name string, env map[string]string, targets map[string]FileTarget) Target {
	t := targets[name]

	for _, name := range t.Imports {
		targets[name] = makeTarget(name, targets)

		mergeStruct(&t, targets[name])
	}

	expandFunc := func(k string) string {
		return env[k]
	}

	ports := []Port{}
	for _, port := range t.Props {
		port.ContainerPort = os.Expand(port.ContainerPort, expandFunc)
		ports = append(ports, port)
	}

	t.Props = ports
	t.Command = os.Expand(t.Command, expandFunc)

	return t.Target
}

func mergeField(val, mergeVal reflect.Value) {
	merge := func(field, field2 reflect.Value) {
		//field.Interface() 当前持有的值
		//reflect.Zero 根据类型获取对应的 零值
		//这个必须调用 Interface 方法 否则为 reflect.Value 构造体的对比 而不是两个值的对比
		//这个地方不要用等号去对比 因为golang 切片类型是不支持 对比的

		if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) { //如果第一个构造体某个字段对应类型的默认值

			if !reflect.DeepEqual(field2.Interface(), reflect.Zero(field2.Type()).Interface()) { //如果第二个构造体 这个字段不为空

				if field.CanSet() != true { //如果不可以设置值 直接返回

					fmt.Println("not set value")
					return
				}

				field.Set(field2) //设置值
			}
		}
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)       //返回结构体的第i个字段
			field2 := mergeVal.Field(i) //返回结构体的第i个字段

			switch field.Kind() {
			case reflect.Struct:
				mergeField(field, field2)
			default:
				merge(field, field2)
			}

		}
	default:
		merge(val, mergeVal)
	}
}

func mergeStruct(val, mergeVal interface{}) {
	v1 := reflect.ValueOf(val).Elem() //初始化为c1保管的具体值的v1
	v2 := reflect.ValueOf(mergeVal)   //初始化为c2保管的具体值的v2

	for i := 0; i < v1.NumField(); i++ {
		mergeField(v1.Field(i), v2.Field(i))
	}
}
