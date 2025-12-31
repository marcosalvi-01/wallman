package common

import "regexp"

// ImageRegex matches common image file extensions
var ImageRegex = regexp.MustCompile(`^.*\.(jpe?g|png|bmp|webp)$`)

// IsImage checks if a file name corresponds to a supported image format
func IsImage(fileName string) bool {
	return ImageRegex.Match([]byte(fileName))
}
