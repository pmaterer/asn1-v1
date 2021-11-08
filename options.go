package asn1

type options struct {
	private    bool
	enumerated bool
	stringType Tag
	timeType   Tag
}

func parseOptions(optionsStr string) (opts options) {
	var option string
	for len(optionsStr) > 0 {
		option, optionsStr, _ = cut(optionsStr, ",")
		switch {
		case option == "private":
			opts.private = true
		case option == "ia5string":
			opts.stringType = TagIA5String
		case option == "utc":
			opts.timeType = TagUTCTime
		case option == "enumerated":
			opts.enumerated = true
		}
	}

	return
}
