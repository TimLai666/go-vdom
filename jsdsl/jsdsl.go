// jsdsl.go
package jsdsl

import (
	"fmt"
	"strings"

	. "github.com/TimLai666/go-vdom/vdom"
)

type Elem struct {
	Selector string
	VarName  string
}

type ElemList struct {
	Selector string
}

func El(selector string) Elem {
	return Elem{Selector: selector}
}

func Els(selector string) ElemList {
	return ElemList{Selector: selector}
}

func (e Elem) Ref() string {
	if e.VarName != "" {
		return e.VarName
	}
	return fmt.Sprintf("document.querySelector('%s')", e.Selector)
}

func (e ElemList) Ref() string {
	return fmt.Sprintf("document.querySelectorAll('%s')", e.Selector)
}

func Do(actions ...JSAction) JSAction {
	var sb strings.Builder
	for _, a := range actions {
		line := strings.TrimSpace(a.Code)
		if !strings.HasSuffix(line, ";") {
			line += ";"
		}
		sb.WriteString(line + "\n")
	}
	return JSAction{Code: sb.String()}
}

func DomReady(actions ...JSAction) JSAction {
	return JSAction{Code: fmt.Sprintf(`document.addEventListener("DOMContentLoaded", function() {
%s
});`, indent(Do(actions...).Code, "  "))}
}

func SetText(el Elem, text string) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.innerText = '%s'`, el.Ref(), escapeJS(text))}
}

func AddClass(el Elem, class string) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.classList.add('%s')`, el.Ref(), class)}
}

func RemoveClass(el Elem, class string) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.classList.remove('%s')`, el.Ref(), class)}
}

func Log(msg string) JSAction {
	return JSAction{Code: fmt.Sprintf(`console.log('%s')`, escapeJS(msg))}
}

func Redirect(url string) JSAction {
	return JSAction{Code: fmt.Sprintf(`location.href = '%s'`, url)}
}

func OnClick(el Elem, action JSAction) JSAction {
	return JSAction{Code: fmt.Sprintf(`%s.addEventListener('click', function() {
%s
});`, el.Ref(), indent(action.Code, "  "))}
}

func Alert(jsExpr string) JSAction {
	return JSAction{Code: fmt.Sprintf(`alert(%s)`, jsExpr)}
}

func InnerText(el Elem) string {
	return fmt.Sprintf("%s.innerText", el.Ref())
}

func ForEach(arrayExpr string, fn func(el Elem) JSAction) JSAction {
	el := "el"
	return JSAction{
		Code: fmt.Sprintf(`%s.forEach(function(%s) {
%s
});`, arrayExpr, el, indent(fn(Elem{VarName: el}).Code, "  ")),
	}
}

func QueryEach(list ElemList, fn func(el Elem) JSAction) JSAction {
	return ForEach(list.Ref(), fn)
}

func indent(code string, prefix string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

func escapeJS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
