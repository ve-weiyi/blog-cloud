package copyutil

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"

	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

/**
 * @StructComment: 利用gob进行深拷贝
 */
func DeepCopyByGob(src, dst interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buffer).Decode(dst)
}

/**
 * @StructComment: 利用json进行深拷贝    obj,&objTo
 */
func DeepCopyByJson(src, dst any) error {
	if tmp, err := jsoniter.Marshal(&src); err != nil {
		return err
	} else {
		err = jsoniter.Unmarshal(tmp, dst)
		return err
	}
}

/*
*
  - @StructComment: 利用反射进行深拷贝    obj,&objTo

参数传递时，第src使用指针还是实例请自行斟酌，dst必须是指针，涉及的字段必须是对外的
*/
func DeepCopyByReflect(src, dst interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针，.Elem()类似于*ptr的操作返回指针指向的地址反射类型
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

// 驼峰式json
func DeepCopyByJsonCamelOrCase(src, dst any, useCamel bool) error {
	if useCamel {
		if tmp, err := jsoniter.Marshal(&src); err != nil {
			return err
		} else {
			tmp = []byte(jsonconv.Case2Camel(string(tmp)))
			err = jsoniter.Unmarshal(tmp, dst)
			return err
		}

	} else {
		if tmp, err := jsoniter.Marshal(&src); err != nil {
			return err
		} else {
			tmp = []byte(jsonconv.Camel2Case(string(tmp)))
			err = jsoniter.Unmarshal(tmp, dst)
			return err
		}
	}
}
