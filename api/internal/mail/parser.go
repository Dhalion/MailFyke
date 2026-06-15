package mail

type ParsedEmail struct {
	From          string
	To            []string
	Subject       string
	BodyHTML      string
	BodyText      string
	Headers       map[string]string
	RawEML        string
	HasAttachments bool
	SizeBytes     int
}

func Parse(rawEML string) (*ParsedEmail, error) {
	// TODO: implement enmime parsing
	return &ParsedEmail{RawEML: rawEML}, nil
}
