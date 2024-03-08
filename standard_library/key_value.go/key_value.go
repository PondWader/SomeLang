package keyvalue

import (
	"database/sql"
	"fmt"
	"main/interop"
	"main/interpreter"
)

var OpenDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args: []interpreter.TypeDef{
		interpreter.GenericTypeDef{
			Type: interpreter.TypeString,
		},
	},
	ReturnType: interpreter.StructDef{
		GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeStruct},
		Properties: map[string]int{
			"get": 0,
			"set": 1,
		},
		PropertyDefs: []interpreter.TypeDef{
			GetDef,
			SetDef,
		},
	},
}

type KeyValueDb struct {
	db *sql.DB
}

func Open(file string) []any {
	db, err := sql.Open("sqlite3", file)
	fmt.Println(db, err)
	return interop.CreateRuntimeStruct(&KeyValueDb{db})
}

var GetDef = interpreter.FuncDef{}

func (kv *KeyValueDb) Get(self []any, key string) {
	result, err := self[0].(*sql.DB).Query("SELECT * FROM key_value WHERE key = ?", key)
	fmt.Println(result, err)
}

var SetDef = interpreter.FuncDef{}

func (kv *KeyValueDb) Set(self []any, key string, value string) {

}
