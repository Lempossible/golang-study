package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	once   sync.Once
	person atomic.Value
)

type Person struct {
	Name string
	Age  int
}

func doOnce(s string) {
	// once.Do(func() {
	fmt.Println(s)
	// })
}

// 这里对person进行整体更新, 而不是分别更新name和age
// 如果person不是原子变量, 则需要使用锁来保护对person的更新, 避免竞态条件
func RightUpdatePerson(name string, age int) {
	newPerson := Person{}
	newPerson.Name = name
	newPerson.Age = age

	// atomic.Value 实现多字段原子赋值的原理不是并发操作同一块多字段内存, 而是每次 Store 时都用新内存替换旧内存. Store 和 Load 都不涉及内存拷贝, 只涉及指针操作
	person.Store(newPerson)
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

	s.CompareAndSwap("hello", "world")
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

	for i := 0; i < 5; i++ {
		// wg.Add(1)
		go func() {
			fmt.Println(i)
			// doOnce(fmt.Sprintf("once %d", i))
		}() // 这里最后的结果是只打印一次 "once "
	}
	// wg.Wait()
	time.Sleep(100 * time.Millisecond)

	// sync.Cond 是一个并发安全的条件变量, 可以在多个 goroutine 中等待和通知
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 5; i++ {
		go listen(c)
	}
	time.Sleep(1 * time.Second)
	go broadcast(c)

	// sync.Pool 是一个并发安全的对象池, 可以在多个 goroutine 中安全地读写
	syncPoolTest()

	time.Sleep(3 * time.Second)
}

var status int64

func broadcast(c *sync.Cond) {
	c.L.Lock()
	atomic.StoreInt64(&status, 1)
	c.Broadcast()
	c.L.Unlock()
}

// sync.Cond的L.Lock() 必须在调用Wait()之前调用, 否则会导致死锁,
// sync.Cond的Wait() 需要Broadcast() 或 Signal() 来唤醒等待的 goroutine
// 因为Wait()会执行两个步骤:
// 1. runtime_notifyListAdd() 将等待计数器加1并解锁.
// 2. runtime_notifyListWait() 会获取当前Goroutine并加到Goroutine通知链尾部, 等待其他Groutine的唤醒并加锁
func listen(c *sync.Cond) {
	c.L.Lock()
	for atomic.LoadInt64(&status) != 1 {
		c.Wait()
	}
	fmt.Println("listen")
	c.L.Unlock()
}

var makeBytesCount int32

func syncPoolTest() {
	pool := sync.Pool{
		New: func() any {
			atomic.AddInt32(&makeBytesCount, 1)
			return make([]byte, 0, 128)
		},
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 1024; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			buf := pool.Get().([]byte)
			str := fmt.Sprintf("协程%d的临时数据 ", id)
			buf = append(buf, str...)
			fmt.Printf("缓冲区内容 - %s（长度：%d）\n", buf, len(buf))
			buf = buf[:0] // 清空缓冲区, 避免后续使用时受到之前数据的影响
			pool.Put(buf) // 注意: 这里必须要 Put, 否则会导致内存泄漏
		}(i)
	}
	wg.Wait()

	fmt.Printf("create new bytes %d times\n", atomic.LoadInt32(&makeBytesCount))
}
