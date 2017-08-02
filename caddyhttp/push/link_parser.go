package push

import (
	"regexp"
	"strings"
)

const (
	commaSeparator     = ","
	semicolonSeparator = ";"
	equalSeparator     = "="
)

var (
	resourceRegexp = regexp.MustCompile("^<(.*?)>;?(.*)$")
)

type linkResource struct {
	uri    string
	params map[string]string
}

// parseLinkHeader is responsible for parsing Link header and returning list of found resources.
//
// Accepted formats are:
// Link: </resource>; as=script
// Link: </resource>; as=script,</resource2>; as=style
// Link: </resource>;</resource2>
func parseLinkHeader(header string) []linkResource {

	if header == "" {
		return make([]linkResource, 0)
	}

	resources := make([]linkResource, 0)

	for _, link := range strings.Split(header, commaSeparator) {
		l := linkResource{params: make(map[string]string)}

		groups := resourceRegexp.FindAllStringSubmatch(link, -1)

		if len(groups) == 0 {
			continue
		}

		l.uri = strings.TrimSpace(groups[0][1])

		for _, param := range strings.Split(strings.TrimSpace(groups[0][2]), semicolonSeparator) {
			parts := strings.SplitN(strings.TrimSpace(param), equalSeparator, 2)

			key := strings.TrimSpace(parts[0])

			if key == "" {
				continue
			}

			if len(parts) == 1 {
				l.params[key] = key
			}

			if len(parts) == 2 {
				l.params[key] = strings.TrimSpace(parts[1])
			}
		}

		resources = append(resources, l)
	}

	return resources
}
