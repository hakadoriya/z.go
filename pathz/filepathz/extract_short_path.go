package filepathz

import "strings"

func ExtractShortPath(path string) string {
	// path == /path/to/directory/file
	//                           ~ <- idx
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}

	// path[:idx] == /path/to/directory
	//                       ~ <- idx
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}

	// path == /path/to/directory/file
	//                  ~~~~~~~~~~~~~~ <- filepath[idx+1:]
	return path[idx+1:]
}
