package converter

func DtoVerify(verified bool, imageUrl string) map[string]interface{} {
	data := make(map[string]interface{})

	data["verified"] = verified

	if verified {
		data["passport"] = imageUrl
	} else {
		data["passport"] = nil
	}

	return data
}
