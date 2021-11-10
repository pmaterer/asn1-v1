package asn1

type options struct {
	private    bool
	enumerated bool
	stringType Tag
	timeType   Tag
}

func NewDefaultOptions() options {
	return options{
		private:    false,
		enumerated: false,
		stringType: TagUTF8String,
		timeType:   TagGeneralizedTime,
	}
}

func parseOptions(optionsStr string) options {
	var option string
	opts := NewDefaultOptions()
	for len(optionsStr) > 0 {
		option, optionsStr, _ = cut(optionsStr, ",")
		switch {
		case option == "private":
			opts.private = true
		case option == "ia5":
			opts.stringType = TagIA5String
		case option == "printable":
			opts.stringType = TagPrintableString
		case option == "utc":
			opts.timeType = TagUTCTime
		case option == "enumerated":
			opts.enumerated = true
		case option == "bitstring":
			opts.stringType = TagBitString
		}
	}

	return opts
}
