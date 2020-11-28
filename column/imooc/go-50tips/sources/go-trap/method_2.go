package main

type foo struct{}

func (foo) methodWithValueReceiver() {
	println("methodWithValueReceiver invoke ok")
}

func (*foo) methodWithPointerReceiver() {
	println("methodWithPointerReceiver invoke ok")
}

type fooer interface {
	methodWithPointerReceiver()
}

func main() {
	f := foo{}
	pf := &f

	f.methodWithPointerReceiver() // 值类型实例调用采用指针类型receiver的方法 ok
	pf.methodWithValueReceiver()  // 指针类型实例调用采用值类型receiver的方法 ok

	// var i fooer = f  // 错误：f并未实现methodWithPointerReceiver
	var i fooer = pf // ok
	i.methodWithPointerReceiver()
}
