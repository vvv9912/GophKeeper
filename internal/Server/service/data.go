package service

import (
	"GophKeeper/pkg/ShaHash"
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"GophKeeper/pkg/store"
	"GophKeeper/pkg/store/postgresql"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// CreateCredentials - Создание пары логин/пароль.
func (s *UseCase) CreateCredentials(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error) {

	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return nil, err
	}
	userData, err := s.StoreData.CreateCredentials(ctx, userId, data, name, description, hash)
	if err != nil {
		return nil, err
	}
	resp := &RespData{
		UserDataId: userData.UserDataId,
		Hash:       hash,
		CreatedAt:  userData.CreatedAt,
		UpdateAt:   userData.UpdateAt,
	}

	return resp, nil
}

// CreateCreditCard - Создание пары данные банковских карт.
func (s *UseCase) CreateCreditCard(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error) {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return nil, err
	}
	userData, err := s.StoreData.CreateCreditCard(ctx, userId, data, name, description, hash)
	if err != nil {
		return nil, err
	}
	resp := &RespData{
		UserDataId: userData.UserDataId,
		Hash:       hash,
		CreatedAt:  userData.CreatedAt,
		UpdateAt:   userData.UpdateAt,
	}
	return resp, nil
}

// createPathIfNotExists - создание пути.
func createPathIfNotExists(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			logger.Log.Error("createPathIfNotExists", zap.Error(err))
			return nil
		}
		logger.Log.Error("createPathIfNotExists", zap.Error(err))
		return err
	}
	return nil
}

// moveFile - перемещение файла.
func moveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		logger.Log.Error("moveFile", zap.Error(err))
		return err
	}
	return nil
}

// CreateFileChunks - Создание бинарных данных.
func (s *UseCase) CreateFileChunks(ctx context.Context, userId int64, tmpFile *TmpFile, name, description string, encryptedData []byte) (*RespData, error) {

	// Создаем путь
	pathStorage := path.Join("storage", strconv.Itoa(int(userId)))

	if err := createPathIfNotExists(pathStorage); err != nil {
		return nil, err
	}

	if err := moveFile(tmpFile.PathFileSave, path.Join(pathStorage, tmpFile.Uuid)); err != nil {
		return nil, err
	}

	// Структура с метаданными
	metaData := &store.MetaData{
		FileName: tmpFile.Uuid,
		PathSave: pathStorage,
		Size:     tmpFile.Size,
	}

	data, err := json.Marshal(metaData)
	if err != nil {
		return nil, err
	}
	// считаем хэш метаднныех
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return nil, err
	}

	// Сохраняем структуру с описанием файла.
	userData, err := s.StoreData.CreateFileDataChunks(ctx, userId, encryptedData, name, description, hash, metaData)
	if err != nil {
		return nil, err
	}
	resp := &RespData{
		UserDataId: userData.UserDataId,
		Hash:       hash,
		CreatedAt:  userData.CreatedAt,
		UpdateAt:   userData.UpdateAt,
	}
	return resp, nil

}

// CreateFile - Создание  данных (файл).
func (s *UseCase) CreateFile(ctx context.Context, userId int64, data []byte, name, description string) (*RespData, error) {
	hash, err := s.createData(ctx, userId, data, name, description)
	if err != nil {
		return nil, err
	}

	// Храним путь на файл.
	userData, err := s.StoreData.CreateFileData(ctx, userId, data, name, description, hash)
	if err != nil {
		return nil, err
	}
	resp := &RespData{

		UserDataId: userData.UserDataId,
		Hash:       hash,
		CreatedAt:  userData.CreatedAt,
		UpdateAt:   userData.UpdateAt,
	}
	return resp, nil

}

// createData - проверка правильности данных и расчет хэша.
func (s *UseCase) createData(ctx context.Context, userId int64, data []byte, name, description string) (string, error) {
	var err error
	if data == nil || len(data) == 0 {
		logger.Log.Error("data is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "data is empty"))
	}
	if name == "" {
		logger.Log.Error("name is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "name is empty"))
	}
	if description == "" {
		logger.Log.Error("description is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "description is empty"))
	}
	if userId == 0 {
		logger.Log.Error("userId is empty")
		err = errors.Join(err, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty"))
	}
	//todo в вывод добавить unwrap
	if err != nil {
		return "", err
	}
	// Считаем хэш полученных данных
	hash := ShaHash.Sha256Hash(data)
	return hash, err
}

// ChangeData - проверка изменения данных.
func (s *UseCase) ChangeData(ctx context.Context, userId int64, userDataId int64, lastTimeUpdate time.Time) ([]byte, error) {
	ok, _ := s.StoreData.ChangeData(ctx, userId, userDataId, lastTimeUpdate)
	resp := struct {
		Status bool `json:"status"`
	}{
		Status: ok,
	}

	response, err := json.Marshal(resp)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		return nil, err
	}
	return []byte(response), nil
}

// ChangeAllData - список изменненых данных.
func (s *UseCase) ChangeAllData(ctx context.Context, userId int64, lastTimeUpdate time.Time) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	data, err := s.StoreData.ChangeAllData(ctx, userId, lastTimeUpdate)
	if err != nil {
		return nil, err
	}
	resp, err := json.Marshal(data)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		return nil, err
	}
	return resp, nil
}

// GetFileSize - получение размера бинарного файла.
func (s *UseCase) GetFileSize(ctx context.Context, userId int64, userDataId int64) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	size, err := s.StoreData.GetFileSize(ctx, userId, userDataId)
	if err != nil {
		return nil, err
	}
	resp := struct {
		FileSize int64 `json:"fileSize"`
	}{
		FileSize: size,
	}

	response, err := json.Marshal(resp)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		return nil, err
	}
	return response, nil
}

// GetFileChunks - получение чанков бинарного файла.
func (s *UseCase) GetFileChunks(ctx context.Context, userId int64, userDataId int64, r *http.Request) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	rangeMin, rangeMax, totalSize, err := ParserContentRange(r.Header.Get("Content-Range"))
	if err != nil {
		return nil, err
	}

	metaData, err := s.StoreData.GetMetaData(ctx, userId, userDataId)
	if err != nil {
		return nil, err
	}
	if metaData.Size != int64(totalSize) {
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "Wrong total size")
	}

	if rangeMin > rangeMax || rangeMin < 0 {
		err = customErrors.NewCustomError(nil, http.StatusBadRequest, "RangeMin must be 0 for first request.")
		return nil, err
	}
	if rangeMax > totalSize {
		err = customErrors.NewCustomError(nil, http.StatusBadRequest, "RangeMax more than totalSize.")
		return nil, err
	}
	data, err := s.getFile(ctx, path.Join(metaData.PathSave, metaData.FileName), rangeMin, rangeMax)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (s *UseCase) getFile(ctx context.Context, pathFile string, byteStart int, byteEnd int) ([]byte, error) {
	f, err := os.OpenFile(pathFile, os.O_RDONLY, 0644)
	if err != nil {
		logger.Log.Error("error open file", zap.String("path", pathFile), zap.Error(err))
		return nil, err
	}

	_, err = f.Seek(int64(byteStart), 0)
	if err != nil {
		logger.Log.Error("error seek file", zap.String("path", pathFile), zap.Error(err))
		return nil, err
	}

	data := make([]byte, byteEnd-byteStart)
	_, err = f.Read(data)
	if err != nil {
		logger.Log.Error("error read file", zap.String("path", pathFile), zap.Error(err))
		return nil, err
	}

	err = f.Close()
	if err != nil {
		logger.Log.Error("error close file", zap.String("path", pathFile), zap.Error(err))
		return nil, err
	}

	return data, nil
}

// GetData - получение данных.
func (s *UseCase) GetData(ctx context.Context, userId int64, userDataId int64) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	// Получаем ингформацию из бд о файле
	usersData, data, err := s.StoreData.GetData(ctx, userId, userDataId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		InfoUsersData *store.UsersData `json:"infoUsersData"`
		EncryptData   *store.DataFile  `json:"encryptData"`
	}

	resp := Data{
		InfoUsersData: usersData,
		EncryptData:   data,
	}
	response, err := json.Marshal(resp)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		return nil, err
	}
	return response, nil
}

// UpdateData - обновление данных.
func (s *UseCase) UpdateData(ctx context.Context, userId int64, userDataId int64, data []byte) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	hash := ShaHash.Sha256Hash(data)

	// todo проверка, если данные уже обновлены с другого устр-ва
	userData, err := s.StoreData.UpdateData(ctx, userId, userDataId, data, hash)
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(userData)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		return nil, err
	}

	return resp, nil

}

// RemoveData - удаление данных (выставление флага в бд).
func (s *UseCase) RemoveData(ctx context.Context, userId, userDataId int64) error {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}
	if userDataId == 0 {
		logger.Log.Error("userDataId is empty")
		return customErrors.NewCustomError(nil, http.StatusBadRequest, "userDataId is empty")
	}

	err := s.StoreData.RemoveData(ctx, userId, userDataId)
	if err != nil {
		return err
	}
	return nil

}

// GetListData - получение списка данных для пользователя.
func (s *UseCase) GetListData(ctx context.Context, userId int64) ([]byte, error) {
	if userId == 0 {
		logger.Log.Error("userId is empty")
		return nil, customErrors.NewCustomError(nil, http.StatusBadRequest, "userId is empty")
	}

	data, err := s.StoreData.GetListData(ctx, userId)
	if err != nil {
		return nil, err
	}
	var resp []struct {
		Id          int64  `json:"id"`
		DataType    string `json:"dataType,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description"`
	}

	for _, val := range data {
		dataType, ok := postgresql.DataType[val.DataType]
		if !ok {
			dataType = "Unknown"
		}

		resp = append(resp, struct {
			Id          int64  `json:"id"`
			DataType    string `json:"dataType,omitempty"`
			Name        string `json:"name,omitempty"`
			Description string `json:"description"`
		}{
			Id:          val.DataId,
			DataType:    dataType,
			Name:        val.Name,
			Description: val.Description,
		})
	}

	response, err := json.Marshal(resp)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Marshal data error")
		logger.Log.Error("Marshal data error", zap.Error(err))
		return nil, err
	}
	return response, nil
}

// UploadFile - загрузка файла.
func (s *UseCase) UploadFile(additionalPath string, r *http.Request) (bool, *TmpFile, error) {
	return s.FileSaver.UploadFile(additionalPath, r)
}

// UpdateBinaryFile - Обновление бинарных данных.
func (s *UseCase) UpdateBinaryFile(ctx context.Context, userId int64, userDataId int64, tmpFile *TmpFile, encryptedData []byte) (*RespData, error) {

	// Создаем путь
	pathStorage := path.Join("storage", strconv.Itoa(int(userId)))

	if err := createPathIfNotExists(pathStorage); err != nil {
		return nil, err
	}

	if err := moveFile(tmpFile.PathFileSave, path.Join(pathStorage, tmpFile.Uuid)); err != nil {
		return nil, err
	}

	// Структура с метаданными
	metaData := &store.MetaData{
		FileName: tmpFile.Uuid,
		PathSave: pathStorage,
		Size:     tmpFile.Size,
	}

	data, err := json.Marshal(metaData)
	if err != nil {
		return nil, err
	}
	// считаем хэш метаднныех
	hash := ShaHash.Sha256Hash(data)
	if err != nil {
		err := customErrors.NewCustomError(err, http.StatusInternalServerError, "Hash error")
		return nil, err
	}

	// Сохраняем структуру с описанием файла.
	userData, err := s.StoreData.UpdateBinaryFile(ctx, userId, userDataId, encryptedData, hash, data)
	if err != nil {
		return nil, err
	}

	resp := &RespData{
		UserDataId: userData.UserDataId,
		Hash:       hash,
		CreatedAt:  userData.CreatedAt,
		UpdateAt:   userData.UpdateAt,
	}
	return resp, nil

}
