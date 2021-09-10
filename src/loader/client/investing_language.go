package client

type InvestingLanugage struct {
	Id     int32
	Code   string
	Domain string
}

var InvestingLanguagesMap = map[int]InvestingLanugage{
	1:  {1, "en", "www"},
	2:  {2, "he", "il"},
	3:  {3, "ar", "sa"},
	4:  {4, "es", "es"},
	5:  {5, "fr", "fr"},
	6:  {6, "zh-hans", "cn"},
	7:  {7, "ru", "ru"},
	8:  {8, "de", "de"},
	9:  {9, "it", "it"},
	10: {10, "tr", "tr"},
	11: {11, "ja", "jp"},
	12: {12, "pt", "pt"},
	13: {13, "sv", "se"},
	14: {14, "el", "gr"},
	15: {15, "pl", "pl"},
	16: {16, "nl", "nl"},
	17: {17, "fi", "fi"},
	18: {18, "ko", "kr"},
	52: {52, "vi", "vn"},
	53: {53, "th", "th"},
	54: {54, "id", "id"},
	55: {55, "zh-hant", "hk"},
	58: {58, "ms", "ms"},
	73: {73, "hi", "hi"},
}
