package common

import "regexp"

// ImageRegex matches JPEG and PNG image file extensions
var ImageRegex = regexp.MustCompile(`^.*\.(jpeg|png)$`)

// IsImage checks if a file name corresponds to a supported image format
func IsImage(fileName string) bool {
	return ImageRegex.Match([]byte(fileName))
}
