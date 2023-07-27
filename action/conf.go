package action

var (
	// terminator sequence for splitting chunks.
	// Required for chunks workaround solution.
	// terminator used in at least 2 places:
	// action.Flush
	// and action.FuncMap["client"] (passed to the client part).
	terminator = []byte("=!EOC!=")
)
