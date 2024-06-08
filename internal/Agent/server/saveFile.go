package server

import (
	"GophKeeper/pkg/customErrors"
	"net/http"
	"os"
	"path"
)

type SaveFile struct {
	FileName string `json:"fileName"`
	pathFile string
	f        *os.File
}

func NewSaveFile(fileName string) (*SaveFile, error) {
	// Создание локальной папки для хранения временных файлов пользователя.
	pathSave := path.Join("./tmp", "agent")
	if err := os.MkdirAll(pathSave, 0755); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error creating directory.")
		return nil, err
	}

	PathFile := path.Join(pathSave, fileName)
	f, err := os.OpenFile(PathFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error during create tmp file.")
		return nil, err
	}

	return &SaveFile{FileName: fileName, f: f, pathFile: PathFile}, nil
}

func (s *SaveFile) GetPathFile() string {
	return s.pathFile
}

func (s *SaveFile) Write(data []byte) (int, error) {
	return s.f.Write(data)
}

func (s *SaveFile) CloseFile() error {
	return s.f.Close()
}
