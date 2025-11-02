package main

import (
	"fmt"

	. "github.com/TimLai666/go-vdom/dom"
)

func main() {
	fmt.Println("=== PropsDef 類型使用示例 ===\n")

	// 示例 1: 使用布林值控制組件行為
	ConfigPanel := Component(
		Div(Props{
			"class": "config-panel",
			"style": "padding: 20px; border: 1px solid #ccc;",
		},
			H3("配置面板"),
			// 這些屬性不在模板的 Props 中，會保持原始類型
			// 條件渲染需要在 Go 層面處理
			Div(Props{"data-debug": "{{debugMode}}"}, P("配置信息")),
			P("最大用戶數: {{maxUsers}}"),
			P("超時時間: {{timeout}} 秒"),
		),
		nil,
		PropsDef{
			"debugMode":    false, // 布林值
			"showAdvanced": false, // 布林值
			"maxUsers":     100,   // 整數
			"timeout":      30,    // 整數
		},
	)

	// 使用預設值
	result1 := ConfigPanel(Props{})
	fmt.Println("1. 使用預設值:")
	fmt.Printf("   debugMode: %T = %v\n", result1.Props["debugMode"], result1.Props["debugMode"])
	fmt.Printf("   maxUsers: %T = %v\n", result1.Props["maxUsers"], result1.Props["maxUsers"])

	// 覆寫部分值
	result2 := ConfigPanel(Props{
		"debugMode": true,
		"maxUsers":  200,
	})
	fmt.Println("\n2. 覆寫部分值:")
	fmt.Printf("   debugMode: %T = %v\n", result2.Props["debugMode"], result2.Props["debugMode"])
	fmt.Printf("   maxUsers: %T = %v\n", result2.Props["maxUsers"], result2.Props["maxUsers"])
	fmt.Printf("   showAdvanced: %T = %v\n", result2.Props["showAdvanced"], result2.Props["showAdvanced"])

	// 示例 2: 使用數字類型進行計算
	PriceCard := Component(
		Div(Props{
			"class": "price-card",
			"style": "border: 2px solid #3b82f6; padding: 15px;",
		},
			H4("{{productName}}"),
			P("價格: ${{price}}"),
			P("數量: {{quantity}}"),
			P("折扣: {{discount}}%"),
			// 注意：這裡我們可以在 Go 層面進行計算
		),
		nil,
		PropsDef{
			"productName": "商品",
			"price":       99.99, // float64
			"quantity":    1,     // int
			"discount":    0,     // int
		},
	)

	result3 := PriceCard(Props{
		"productName": "精選商品",
		"price":       199.99,
		"quantity":    3,
		"discount":    10,
	})
	fmt.Println("\n3. 數字類型示例:")
	fmt.Printf("   price: %T = %v\n", result3.Props["price"], result3.Props["price"])
	fmt.Printf("   quantity: %T = %v\n", result3.Props["quantity"], result3.Props["quantity"])

	// 計算總價（可以在 Go 層面進行）
	price := result3.Props["price"].(float64)
	quantity := result3.Props["quantity"].(int)
	discount := result3.Props["discount"].(int)
	total := price * float64(quantity) * (1 - float64(discount)/100)
	fmt.Printf("   計算總價: $%.2f\n", total)

	// 示例 3: 使用複雜類型
	TagCloud := Component(
		Div(Props{
			"class": "tag-cloud",
			"style": "padding: 10px;",
		},
			H4("標籤雲"),
			Div(Props{"id": "tags-container"}, P("標籤將在此顯示")),
		),
		nil,
		PropsDef{
			"tags": []string{"Go", "VDOM", "Web"}, // 字串切片
		},
	)

	result4 := TagCloud(Props{
		"tags": []string{"Go", "JavaScript", "HTML", "CSS", "VDOM"},
	})
	fmt.Println("\n4. 複雜類型示例:")
	fmt.Printf("   tags: %T = %v\n", result4.Props["tags"], result4.Props["tags"])

	// 示例 4: 混合使用多種類型
	UserProfile := Component(
		Div(Props{
			"class": "user-profile",
			"style": "border: 1px solid #e5e7eb; padding: 20px; border-radius: 8px;",
		},
			H3("用戶資料"),
			P("姓名: {{name}}"),
			P("年齡: {{age}}"),
			P("郵箱: {{email}}"),
			Span(Props{"data-active": "{{isActive}}"}, "帳號狀態"),
			P("權限級別: {{level}}"),
		),
		nil,
		PropsDef{
			"name":     "訪客",
			"age":      0,
			"email":    "",
			"isActive": false,
			"level":    1,
			"tags":     []string{},
			"metadata": map[string]interface{}{},
		},
	)

	result5 := UserProfile(Props{
		"name":     "張三",
		"age":      28,
		"email":    "zhangsan@example.com",
		"isActive": true,
		"level":    5,
		"tags":     []string{"VIP", "活躍用戶"},
		"metadata": map[string]interface{}{
			"lastLogin": "2025-01-24",
			"posts":     42,
		},
	})

	fmt.Println("\n5. 混合類型示例:")
	fmt.Printf("   name: %T = %v\n", result5.Props["name"], result5.Props["name"])
	fmt.Printf("   age: %T = %v\n", result5.Props["age"], result5.Props["age"])
	fmt.Printf("   isActive: %T = %v\n", result5.Props["isActive"], result5.Props["isActive"])
	fmt.Printf("   tags: %T = %v\n", result5.Props["tags"], result5.Props["tags"])
	fmt.Printf("   metadata: %T = %v\n", result5.Props["metadata"], result5.Props["metadata"])

	// 示例 6: 實際應用 - 可配置的按鈕組件
	CustomButton := Component(
		Button(Props{
			"id":       "{{id}}",
			"class":    "custom-btn",
			"disabled": "{{disabled}}",
			"style":    "padding: {{paddingY}}px {{paddingX}}px; font-size: {{fontSize}}px; border-radius: {{borderRadius}}px; background: {{bgColor}}; color: {{textColor}}; border: none; cursor: pointer;",
		},
			"{{label}}",
		),
		nil,
		PropsDef{
			"id":           "",
			"label":        "按鈕",
			"disabled":     false,
			"paddingX":     20,
			"paddingY":     10,
			"fontSize":     14,
			"borderRadius": 4,
			"bgColor":      "#3b82f6",
			"textColor":    "#ffffff",
		},
	)

	result6 := CustomButton(Props{
		"label":        "提交",
		"paddingX":     30,
		"paddingY":     15,
		"fontSize":     16,
		"borderRadius": 8,
	})

	fmt.Println("\n6. 可配置按鈕示例:")
	fmt.Printf("   disabled: %T = %v\n", result6.Props["disabled"], result6.Props["disabled"])
	fmt.Printf("   paddingX: %T = %v\n", result6.Props["paddingX"], result6.Props["paddingX"])
	fmt.Printf("   fontSize: %T = %v\n", result6.Props["fontSize"], result6.Props["fontSize"])

	// 渲染最終 HTML
	fmt.Println("\n=== 渲染 HTML 輸出 ===\n")

	doc := Document(
		"PropsDef 類型示例",
		nil, // links
		nil, // scripts
		nil, // meta
		Div(nil,
			result2,
			Hr(),
			result3,
			Hr(),
			result4,
			Hr(),
			result5,
			Hr(),
			result6,
		),
	)

	fmt.Println(Render(doc))
}
