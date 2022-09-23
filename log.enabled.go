//go:build KYOTO_VERBOSE

package kyoto

import "log"

func logln(args ...any) {
	log.Println(args...)
}

func logf(row string, args ...any) {
	log.Printf(row, args...)
}
