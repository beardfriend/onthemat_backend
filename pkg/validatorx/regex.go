package validatorx

import "regexp"

func ValidateRegex(regex, value string) bool {
	reg := regexp.MustCompile(regex)
	return reg.Match([]byte(value))
}

const (
	ImageRegex                = `^.(jpe?g|gif|svg|png)$`
	AtLeastOneCharOneNumRegex = "^(?:[0-9~!@$%^&*]+[a-zA-Z!@~!@$%^&*]|[a-zA-Z~!@$%^&*]+[0-9~!@$%^&*])[a-zA-Z0-9~!@$%^&*]*$"
	ForbiddenSpecialCharRegex = "([^\"#%'()+/:;<=>?\\[\\]^{|}~]+)$"
)
