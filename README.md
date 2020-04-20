## Catcher

Catcher can record the position of error automatically. All apis are nestable

### 错误的简单序列化与错误捕捉

一共有三个简易函数用于包裹错误
+ `Wrap(code int, err error) error`
+ `WrapString(code int, err string) error`
+ `WrapCode(code int) error`

最简使用方法：

```go
package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

func TestExampleWrap(t *testing.T) {
	frame := Wrap(233, errors.New("QAQ"))
	fmt.Println(Describer{
		Pack: "github.com/Myriad-Dreamin",
		Rel:  handle(filepath.Abs("")),
	}.Describe(frame))
}
```

```plain
1 <- <pos:<catcher.TestExampleWrap,catcher_example_test.go:11>,code:233,err:QAQ>
```

该方法可以嵌套使用用于分析复杂的远程函数调用：

```go
package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

func innerLogic() error {
	return Wrap(233, errors.New("QAQ"))
}

func outerLogic() error {
	err := innerLogic()
	if err != nil {
		return Wrap(666, err)
	}
	return nil
}

func TestExampleNestWrap(t *testing.T) {
	err := outerLogic()
	fmt.Println(Describer{
		Pack: "github.com/Myriad-Dreamin",
		Rel:  handle(filepath.Abs("")),
	}.Describe(err))
}
```

```plain
2 <- <pos:<catcher.outerLogic,catcher_example_test.go:17>,code:666,err:>
1 <- <pos:<catcher.innerLogic,catcher_example_test.go:11>,code:233,err:QAQ>
```

其他两种方法类似，可以避免创建`error`对象的消耗。

### 错误反序列化

一共有6个反序列化函数

+ `FromBytes(s []byte) (f Frame, ok bool)`
+ `FromString(s string) (f Frame, ok bool)`
+ `FromError(err error) (f Frame, ok bool)`
+ `StackFromBytes(s []byte) (fs Frames, ok bool)`
+ `StackFromString(s string) (fs Frames, ok bool)`
+ `StackFromError(err error) (fs Frames, ok bool)`

前面三个函数可以将错误展开一层，后面三个函数可以帮助你把整个错误栈打开。
以`FromString`为例，其他几个函数使用方法相同：

```go
package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

func TestExampleFromString(t *testing.T) {
	frameToTransfer := Wrap(233, errors.New("QAQ"))
	if frameToTransfer == nil {
		t.Fatal("empty")
	}
	transferObject := frameToTransfer.Error()
	rawFrames, ok := FromString(transferObject)
	if !ok {
		t.Fatal("deserialize error, or not a frame")
	}
	fmt.Println(rawFrames.GetPos().Func.Name, rawFrames.GetCode(), rawFrames.GetErr())
}
```

```plain
github.com/Myriad-Dreamin/catcher.TestExampleFromString 233 QAQ
C:/work/go/src/github.com/Myriad-Dreamin/catcher/catcher_example_test.go 11
```

### 使用Describe辅助函数帮助阅读错误

最简单的办法是直接使用全局函数`Describe`

```go
package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

func TestExampleSimpleNestWrap(t *testing.T) {
	err := outerLogic()
	fmt.Println(Describe(err))
}
```

```plain
2 <- <pos:<<github.com/Myriad-Dreamin/catcher.outerLogic,C:/work/go/src/github.com/Myriad-Dreamin/catcher/catcher_example_test.go:33>,C:/work/go/src/github.com/Myriad-Dreamin/catcher/catcher_example_test.go:36>,code:666,err:>
1 <- <pos:<<github.com/Myriad-Dreamin/catcher.innerLogic,C:/work/go/src/github.com/Myriad-Dreamin/catcher/catcher_example_test.go:30>,C:/work/go/src/github.com/Myriad-Dreamin/catcher/catcher_example_test.go:30>,code:233,err:QAQ>
```

虽然我们得到了所有信息，但不好读，为此可以使用`Describer.Describe`简化信息，这里说明一下`Describer`的两个参数。具体使用方法可以参考第一节的示例代码。

```go
type Describer struct {
    // Pack: 以Pack为基准解释包名
	Pack string
    // Rel: 以Rel为基准解释路径
    Rel string
}
```

### 使用WrapN函数自定义Wrapper

你可以使用下面三个函数自定义包装错误方法
+ `WrapN(skip, code int, err error) error`
+ `WrapStringN(skip, code int, err string) error`
+ `WrapCodeN(skip, code int) error`

实际上`Wrap`相当于`WrapN(skip=BaseSkip)`

这里提供一个自定义包装器的例子：

```go
package catcher

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

type DatabaseCode int

type DuplicateError struct { Data string }

func (d DuplicateError) Error() string {	
	return fmt.Sprintf(`duplicate key at field "%v"`, d.Data)
}

type UserDefinedError struct { Data []interface{} }

func (d UserDefinedError) Error() string {
	return fmt.Sprintf("user defined: %v", d.Data)
}

func DatabaseError(code DatabaseCode, errDesc ...interface{}) error {
	switch code {
	case 1:
		return WrapN(BaseSkip+1, int(code), DuplicateError{Data: errDesc[0].(string)})
	default:
		return WrapN(BaseSkip+1, int(code), UserDefinedError{Data: errDesc})
	}
}

func TestExampleDatabaseError(t *testing.T) {
	err := DatabaseError(1, "id")
	fmt.Println(d.Describe(err))
}

```

```plain
1 <- <pos:<catcher.TestExampleDatabaseError,catcher_example_test.go:34>,code:1,err:duplicate key at field "id">
```

不错，程序正确捕捉到了错误被包装的位置

### 全局设定

##### `CodeDeserializeError`

当发生反序列化错误时，code的值，该code将会成为反序列化结果`Frame.GetCode`的值。

##### `SetErrorFlag`

你可以关闭错误位置的记录，这将节省一些时间。但我建议所有的时候都保持记录，因为错误记录的代价一般每级错误栈的损耗只有一微秒

ErrorFlag的可选值为`Prod`或`Debug`，默认为`Debug`

##### `SetCodeDescriptor`

`_codeDescriptor`用于解释error code，默认值为`strconv.Itoa`

##### `SetReportBad`

如果为`true`，会自动将反序列化错误输出到日志。默认为`true`

##### `SetMagic`

是一个字符串，用于分隔数据，类似于multiform的魔数。

### 重要数据类型

```go
package main

type Frames []Frame
type Frame interface {
	GetPos() Caller
	GetCode() int
	GetErr() string

	Error() string
    String() string
    Bytes() []byte

	Dump() string
	RelDump(pack, rel string) (string, error)
	
	ReleaseError()
}
```

### Benchmark

```plain
goos: windows
goarch: amd64
pkg: github.com/Myriad-Dreamin/catcher
BenchmarkWrap
BenchmarkWrap-8                               	 1239074	      1000 ns/op
BenchmarkWrapWithOutCollectInfo
BenchmarkWrapWithOutCollectInfo-8             	21654830	        52.3 ns/op
BenchmarkWrapChain
BenchmarkWrapChain-8                          	  287155	      4423 ns/op
BenchmarkWrapChainWithOutCollectInfo
BenchmarkWrapChainWithOutCollectInfo-8        	 2344903	       534 ns/op
BenchmarkStackFromError
BenchmarkStackFromError-8                     	  868796	      1494 ns/op
BenchmarkStackFromErrorWithOutCollectInfo
BenchmarkStackFromErrorWithOutCollectInfo-8   	 1405251	       865 ns/op
PASS

Process finished with exit code 0
```

### Future Work

1. `Describer`在多级微服务中表现不理想，因为这个对象只针对单repo，单路径简化解释
2. `frameImpl`不具有同进程的错误对象完全恢复能力
   
   例如 `myerror.SomeError`被该结构体handle的时候会被序列化（调用`SomeError.Error()`）
   
   将来优先考虑增加一个`rawErrFrameImpl.RecoverRaw() (err error)`
   
   初步规划：
   
   + 当元信息为错误对象时，`Recover`将返回同一个对象
   + 当元信息不是错误对象时，返回一个`errors.errorString`类型的字符串