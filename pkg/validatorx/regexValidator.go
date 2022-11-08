package validatorx

import "regexp"

func ImageContentTypeValidator(contentType string) (bool, error) {
	return regexp.MatchString(ImageRegex, contentType)
}
