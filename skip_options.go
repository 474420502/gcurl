package gcurl

type SkipType int

const (
	ST_NotSkipType SkipType = 0
	ST_OnlyOption  SkipType = 1
	ST_WithValue   SkipType = 2
)

var skipList = map[string]SkipType{
	"-O":            ST_OnlyOption,
	"-I":            ST_OnlyOption,
	"--remote-name": ST_OnlyOption,
}

func checkInSkipList(optstr string) SkipType {
	if v, ok := skipList[optstr]; ok {
		return v
	}
	return ST_NotSkipType
}
