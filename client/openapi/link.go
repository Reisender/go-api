package openapi

type Link struct {
	Rel string `json:"rel"`
	URI string `json:"uri"`
}

type Links []Link

func (ls Links) Next() (string, bool) {
	for _, l := range ls {
		if l.Rel == "next" {
			return l.URI, true
		}
	}

	return "", false
}
