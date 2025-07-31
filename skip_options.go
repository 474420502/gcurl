package gcurl

type SkipType int

const (
	ST_NotSkipType SkipType = 0
	ST_OnlyOption  SkipType = 1
	ST_WithValue   SkipType = 2
)

var skipList = map[string]SkipType{
	// 已实现的选项已从跳过列表中移除：
	// "-O": ST_OnlyOption,           // 现在支持 --remote-name
	// "--remote-name": ST_OnlyOption, // 现在支持 --remote-name
	// "-I": ST_OnlyOption,           // 现在支持 --include
}

func checkInSkipList(optstr string) SkipType {
	if v, ok := skipList[optstr]; ok {
		return v
	}
	return ST_NotSkipType
}
