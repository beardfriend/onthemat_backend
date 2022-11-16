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
	PhoneRegex                = `^01([0|1|6|7|8|9])([0-9]{3,4})([0-9]{4})$`
	UrlRegex                  = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
)
