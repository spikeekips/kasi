package kasi_util

import "path"

func JoinPath(base string, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Clean(path.Join(base, target))
}
