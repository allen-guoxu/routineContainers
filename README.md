## routineContainers
简单的协程任务池

### 描述

​	为了解决类似以下问题，如果wg.Add数量错误，比如填少了，可能没有等待全部任务执行完成，程序就退出了，可能导致ctx意外被销毁，而导致正常的任务没被执行完成，或者填多了导致主协程一直阻塞。

​	另外在需求变更时，也有忘记修改Add数量的风险

```go
	wg := sync.WaitGroup{}
	wg.Add(3)  // 这里可能会写错
	
	go func() {
		defer wg.Done()
		doSomeThing(ctx, ...)
	}()
	go func() {
		defer wg.Done()
		doSomeThing(ctx, ...)
	}()
	go func() {
		defer wg.Done()
		doSomeThing(ctx, ...)
	}()
	wg.Wait()

```

### 方案

实现一个简单的自动Add、Done的协程任务池，避免每次修改Add参数

```go
// RoutineContainers 协程任务池
type RoutineContainers struct {
	Contains []func()
	Wg       sync.WaitGroup
}

// Put 协程任务
func (r *RoutineContainers) Put(f func()) {
	tf := func() {
		f()
		r.Wg.Done()
	}
	r.Contains = append(r.Contains, tf)
	r.Wg.Add(1)
}

// Run 协程任务
func (r *RoutineContainers) Run() {
	for idx := range r.Contains {
		go r.Contains[idx]()
	}
	r.Wg.Wait()
    r.Contains = r.Contains[:0]
}
```

### 使用如下

```go
func main() {
	rc := RoutineContainers{}
	fmt.Println("begin")
	rc.Put(func() {
		fmt.Println("do something 1")
		time.Sleep(time.Millisecond * 300)
	})
	rc.Put(func() {
		fmt.Println("do something 2")
		time.Sleep(time.Millisecond * 400)
	})
	rc.Put(func() {
		fmt.Println("do something 3")
		time.Sleep(time.Millisecond * 500)
	})
	rc.Put(func() {
		fmt.Println("do something 4")
		time.Sleep(time.Millisecond * 600)
	})
	rc.Run()
	fmt.Println("done")
}

#####################输出如下#####################
[allenxguo@VM_66_25_centos routinecontainer]$ go run rc.go 
begin
do something 4
do something 1
do something 3
do something 2
done
```
