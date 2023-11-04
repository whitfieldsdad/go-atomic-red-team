package atomic_red_team

type File struct {
	Path string `json:"path"`
}

func GetFile(path string) (*File, error) {
	return &File{
		Path: path,
	}, nil
}
