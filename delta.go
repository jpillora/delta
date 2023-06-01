package delta

import (
	"github.com/r3labs/diff/v3"
)

// differ is a JSON differ
var differ, _ = diff.NewDiffer(
	diff.TagName("json"),
	diff.SliceOrdering(true),
)

// CopyPatch returns a JSON Patch that was used to transform src into dst
func CopyPatch(dst, src any, opts ...Option) ([]Operation, error) {
	// build the copy context
	options := &options{
		apply: true,
	}
	for _, option := range opts {
		option(options)
	}
	// compare the two objects
	changes, err := differ.Diff(src, dst)
	if err != nil {
		return nil, err
	}
	// optionally apply the changes as well
	if options.apply {
		// log.Printf("PRE: SRC %s", j(src))
		// log.Printf("PRE: DST %s", j(dst))
		// log.Printf("CHANGES: %d", len(changes))
		// for _, change := range changes {
		// 	log.Printf("CHANGE: %s", changeToOperation(change))
		// }
		// patches :=
		differ.Patch(changes, src)
		// log.Printf("PATCHES: %d", len(patches))
		// for _, p := range patches {
		// 	log.Printf("PATCH: %s", j(p))
		// 	if !patches.Applied() && p.Errors != nil {
		// 		return nil, fmt.Errorf("patch err: %s", p.Errors)
		// 	}
		// }
		// log.Printf("POST: %s", j(dst))
	}
	return changesToOperations(changes), nil
}
