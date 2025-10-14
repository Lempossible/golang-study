package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var once sync.Once

func doOnce(s string) {
	// once.Do(func() {
	fmt.Println(s)
	// })
}
func main() {

	// atomic
	var counter atomic.Int32
	counter.Add(1)
	fmt.Println(counter.Load())

	// atomic.Value 可以存储任意类型的值
	var s atomic.Value
	s.Store("hello")
	fmt.Println(s.Load())

	// 可以使用类型断言来获取存储的值, 但是需要注意, 如果存储的值不是预期的类型, 会导致 panic
	str := s.Load().(string)
	fmt.Println(str)

	// 为了避免 panic, 可以使用类型断言前先判断类型
	if v, ok := s.Load().(string); ok {
		fmt.Println(v)
	}

	// 也可以使用类型开关来判断类型
	switch v := s.Load().(type) {
	case string:
		fmt.Println(v)
	default:
		fmt.Println("unknown type")
	}

	// sync.Map 是一个并发安全的 map, 可以在多个 goroutine 中安全地读写
	var m sync.Map
	m.Store("a", 1)
	fmt.Println(m.Load("a"))

	// sync.Map 也支持遍历
	m.Range(func(key, value any) bool {
		fmt.Println(key, value)
		return true
	})

	// sync.Map 也支持删除
	m.Delete("a")
	fmt.Println(m.Load("a"))

	// sync.Map 也支持删除后返回值
	if v, ok := m.LoadAndDelete("a"); ok {
		fmt.Println(v)
	}

	// sync.Map 也支持存储后返回值
	if v, ok := m.LoadOrStore("a", 2); ok {
		fmt.Println(v)
	}

	// for 循环遍历 sync.Map 时, 不能使用 range 语句, 因为 range 语句会阻塞 sync.Map 的遍历
	// for key, value := range m {
	//	fmt.Println(key, value)
	// }

	// sync.Once 是一个并发安全的 once 操作, 可以确保某个操作只被执行一次
	// 如果没有sync.Once, 则会打印 "once i" 随机顺序
	// 如果没有sync.Once, go func 不带参数时, 会只打印 "i = 4", golang 1.22 开始, 会随机打印 "i = 9" 到 "i = 0"
	// var wg sync.WaitGroup = sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		// wg.Add(1)
		go func() {
			fmt.Println(i)
			// doOnce(fmt.Sprintf("once %d", i))
		}() // 这里最后的结果是只打印一次 "once "
	}
	// wg.Wait()
	time.Sleep(100 * time.Millisecond)

	type Student struct {
		Name string
		Age  int
	}
	studentList := []*Student{
		{
			Name: "张三",
			Age:  13,
		},
		{
			Name: "李四",
			Age:  13,
		},
		{
			Name: "王五",
			Age:  13,
		},
	}
	for idx, stu := range studentList {
		go func() {
			fmt.Printf("%v: %v\n", idx, stu)
		}()
	}
	time.Sleep(3 * time.Second)
}
