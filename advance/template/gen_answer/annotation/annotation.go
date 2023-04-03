package annotation

import (
	"go/ast"
	"strings"
)

type Annotations[N ast.Node] struct {
	Node N
	Ans  []Annotation
}

func (a Annotations[N]) Get(key string) (Annotation, bool) {
	for _, an := range a.Ans {
		if an.Key == key {
			return an, true
		}
	}
	return Annotation{}, false
}

type Annotation struct {
	Key   string
	Value string
}

func newAnnotations[N ast.Node](n N, cg *ast.CommentGroup) Annotations[N] {
	if cg == nil || len(cg.List) == 0 {
		return Annotations[N]{Node: n}
	}
	ans := make([]Annotation, 0, len(cg.List))
	for _, c := range cg.List {
		text, ok := extractContent(c)
		if !ok {
			continue
		}
		if strings.HasPrefix(text, "@") {
			segs := strings.SplitN(text, " ", 2)
			key := segs[0][1:]
			val := ""
			if len(segs) == 2 {
				val = segs[1]
			}
			ans = append(ans, Annotation{
				Key:   key,
				Value: val,
			})
		}
	}
	return Annotations[N]{
		Node: n,
		Ans:  ans,
	}
}

func extractContent(c *ast.Comment) (string, bool) {
	text := c.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	} else if strings.HasPrefix(text, "/* ") {
		length := len(text)
		return text[3 : length-2], true
	}
	return "", false
}
