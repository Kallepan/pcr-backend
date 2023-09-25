package repository

type SystemRepository interface {
	GetVersion() (map[string]string, error)
}

type SystemRepositoryImpl struct {
}

func (s SystemRepositoryImpl) GetVersion() (map[string]string, error) {
	return map[string]string{
		"version": "1.0",
	}, nil
}

func SystemRepositoryInit() *SystemRepositoryImpl {
	return &SystemRepositoryImpl{}
}
