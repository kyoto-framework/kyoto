//go:build !KYOTO_VERBOSE

package kyoto

func logln(...any) {}

func logf(string, ...any) {}
