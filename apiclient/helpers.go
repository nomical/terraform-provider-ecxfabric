package apiclient

func (ads *ActionDetails) ExtractValueFromActionRequiredDataItem(key string) string {
	for _, ad := range *ads {
		for _, ard := range ad.ActionRequiredData {
			if ard.Key == key {
				return ard.Value
			}
		}
	}

	return ""
}
