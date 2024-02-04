package component

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.kyoto.codes/zen/v3/errorsx"
	"go.kyoto.codes/zen/v3/logic"
)

// Server is a default server component state implementation.
// It uses local temporary files and JSON encoding
// to store, marshal and unmarshal the state.
// Please, make sure this strategy actually fits to your environment.
type Server struct {
	Name

	Path    string        // Path to store component state (default "/tmp/")
	Timeout time.Duration // State timeout (default 24 hours, clean up running on each unmarshal)
}

// path wraps `Path` and resolves with default option.
func (s *Server) path() string {
	return logic.Or(s.Path, "/tmp/")
}

// timeout wraps `Timeout` and resolves with default option.
func (s *Server) timeout() time.Duration {
	return logic.Or(s.Timeout, 24*time.Hour)
}

// cleanup removes outdated state files.
func (s *Server) cleanup() {
	for _, file := range errorsx.Must(os.ReadDir(s.path())) {
		// Pass if file is not a component state
		if !strings.HasSuffix(file.Name(), ".component") {
			continue
		}
		// If creation/modification date is out of timeout bounds,
		// remove that file.
		if time.Since(errorsx.Must(file.Info()).ModTime()) > s.timeout() {
			errorsx.Must(0, os.Remove(path.Join(s.path(), file.Name())))
		}
	}
}

// Marshal encodes state with json into temporary file in `Path` directory.
func (s *Server) Marshal(src any) string {
	// Create new tmp file
	tmp := errorsx.Must(os.CreateTemp(s.path(), "*.component"))
	// Encode state
	errorsx.Must(0, json.NewEncoder(tmp).Encode(src))
	// Return filename as a marshaled state
	return filepath.Base(tmp.Name())
}

// Unmarshal decodes state with json from temporary file in `Path` directory.
// Fires up a cleanup goroutine in the end.
func (s *Server) Unmarshal(dst any, str string) {
	// Open tmp file
	tmp := errorsx.Must(os.OpenFile(path.Join(s.path(), str), os.O_RDONLY, 0777))
	// Decode
	errorsx.Must(0, json.NewDecoder(tmp).Decode(dst))
	// Fire up cleanup
	go s.cleanup()
}
