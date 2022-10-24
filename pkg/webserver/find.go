package webserver

import "github.com/mattfenwick/collections/pkg/slice"
import "github.com/mattfenwick/collections/pkg/function"

func FindKeyInJson(json any, path []any, key string) [][]any {
	var out [][]any
	pathCopy := copy(path)
	switch val := json.(type) {
	case map[string]interface{}:
		for k, v := range val {
			if k == key {
				out = append(out, pathCopy)
			}
			out = append(out, FindKeyInJson(v, slice.Append(pathCopy, []any{k}), key)...)
		}
	case []interface{}:
		for i, x := range val {
			out = append(out, FindKeyInJson(x, slice.Append(pathCopy, []any{i}), key)...)
		}
	default:
	}
	return out
}

func copy[A any](xs []A) []A {
	return slice.Map(function.Id[A], xs)
}
