package validatorx

import "regexp"

func ImageExtensionValidator(contentType string) (bool, error) {
	return regexp.MatchString(ImageRegex, contentType)
}
