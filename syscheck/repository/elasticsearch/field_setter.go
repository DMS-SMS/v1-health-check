// Create file in v.1.0.0
// field_setter.go is file that gather function that satisfy signature of FieldSetter in syscheck.go file.
// Functions of FieldSetter type can be used by handing over to repository constructor

// You can show about package description in syscheck.go file

package elasticsearch

import (
	"log"
	"reflect"
)

// Return closure function set field named 'IndexName' with reflect
func IndexName(name string) FieldSetter {
	fieldName := "IndexName"
	return func(i interface{}) {
		if fieldValue := reflect.ValueOf(i).Elem().FieldByName(fieldName); fieldValue.IsValid() {
			fieldValue.Set(reflect.ValueOf(name))
			return
		}
		log.Fatalf("%s dosn't have %s field", reflect.TypeOf(i).String(), fieldName)
	}
}

// Return closure function set field named 'IndexShardNum' with reflect
func IndexShardNum(shardNum int) FieldSetter {
	fieldName := "IndexShardNum"
	return func(i interface{}) {
		if fieldValue := reflect.ValueOf(i).Elem().FieldByName(fieldName); fieldValue.IsValid() {
			fieldValue.Set(reflect.ValueOf(shardNum))
			return
		}
		log.Fatalf("%s dosn't have %s field", reflect.TypeOf(i).String(), fieldName)
	}
}

// Return closure function set field named 'IndexReplicaNum' with reflect
func IndexReplicaNum(replicaNum int) FieldSetter {
	fieldName := "IndexReplicaNum"
	return func(i interface{}) {
		if fieldValue := reflect.ValueOf(i).Elem().FieldByName(fieldName); fieldValue.IsValid() {
			fieldValue.Set(reflect.ValueOf(replicaNum))
			return
		}
		log.Fatalf("%s dosn't have %s field", reflect.TypeOf(i).String(), fieldName)
	}
}
