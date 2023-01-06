package logerr

type Logerr struct {
	FilePath string
}

func NewLogerr(FilePath string) *Logerr {
	return &Logerr{
		FilePath: FilePath,
	}
}
