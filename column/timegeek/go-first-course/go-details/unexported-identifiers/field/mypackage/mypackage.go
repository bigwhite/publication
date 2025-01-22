package mypackage

type myStruct struct {
	Field string // 导出的字段
}

func NewMyStruct(value string) *myStruct {
	return &myStruct{Field: value}
}

func (m *myStruct) M1() {
	// ...
}
