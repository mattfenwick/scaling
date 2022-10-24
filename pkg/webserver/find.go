package webserver

import (
	"github.com/mattfenwick/collections/pkg/slice"
	"github.com/sirupsen/logrus"
)
import "github.com/mattfenwick/collections/pkg/function"

func FindKeyInJson(json any, path []any, key string) [][]any {
	var out [][]any
	pathCopy := copy(path)
	logrus.Debugf("looking for key %s at path %+v, type %T", key, pathCopy, json)
	switch val := json.(type) {
	case string:
		if val == key {
			logrus.Debugf("found key %s at path %+v", key, pathCopy)
			out = append(out, pathCopy)
		}
	case map[string]interface{}:
		for k, v := range val {
			extendedPath := slice.Append(pathCopy, []any{k})
			if k == key {
				logrus.Debugf("found key %s at path %+v", key, extendedPath)
				out = append(out, extendedPath)
			}
			out = append(out, FindKeyInJson(v, extendedPath, key)...)
		}
	case []interface{}:
		for i, x := range val {
			out = append(out, FindKeyInJson(x, slice.Append(pathCopy, []any{i}), key)...)
		}
	default:
		logrus.Debugf("skipping at %+v, wrong type (%T)", pathCopy, val)
	}
	return out
}

func copy[A any](xs []A) []A {
	return slice.Map(function.Id[A], xs)
}
