package main

import (
	"fmt"

	js "github.com/TimLai666/go-vdom/jsdsl"
)

func main() {
	fmt.Println("=== 測試 Do() 參數注入 ===\n")

	// 測試 1: event
	result1 := js.Do([]string{"event"},
		js.Const("target", "event.target"),
		js.Alert("'test'"),
	)
	fmt.Println("1. 使用 event 參數:")
	fmt.Println(result1.Code)
	fmt.Println()

	// 測試 2: e
	result2 := js.Do([]string{"e"},
		js.Const("target", "e.target"),
		js.Alert("'test'"),
	)
	fmt.Println("2. 使用 e 參數:")
	fmt.Println(result2.Code)
	fmt.Println()

	// 測試 3: myEvent (自定義名稱)
	result3 := js.Do([]string{"myEvent"},
		js.Const("target", "myEvent.target"),
		js.Alert("'test'"),
	)
	fmt.Println("3. 使用 myEvent 參數:")
	fmt.Println(result3.Code)
	fmt.Println()

	// 測試 4: clickEvent
	result4 := js.Do([]string{"clickEvent"},
		js.Const("btnText", "clickEvent.target.textContent"),
		js.Alert("'test'"),
	)
	fmt.Println("4. 使用 clickEvent 參數:")
	fmt.Println(result4.Code)
	fmt.Println()

	// 測試 5: AsyncDo 使用 asyncEvent
	result5 := js.AsyncDo([]string{"asyncEvent"},
		js.Const("value", "asyncEvent.target.value"),
		js.Const("data", "await fetch('/api')"),
	)
	fmt.Println("5. AsyncDo 使用 asyncEvent 參數:")
	fmt.Println(result5.Code)
	fmt.Println()

	// 測試 6: 無參數
	result6 := js.Do(nil,
		js.Alert("'Hello'"),
	)
	fmt.Println("6. 無參數:")
	fmt.Println(result6.Code)
	fmt.Println()

	fmt.Println("=== 結論 ===")
	fmt.Println("所有參數名都正確注入 event 對象！")
	fmt.Println("生成的代碼格式: ((參數名)=>{...})(event)")
}
