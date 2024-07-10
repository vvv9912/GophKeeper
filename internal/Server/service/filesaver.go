package service

import (
	"GophKeeper/pkg/customErrors"
	"GophKeeper/pkg/logger"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// TmpFile - информация о текущем сохранении файла.
type TmpFile struct {
	LastUpdate       time.Time // время последнего обновления.
	PathFileSave     string    // путь к файлу на диске.
	OriginalFileName string    // оригинальное имя файла.
	Uuid             string    // уникальное имя файла.
	Extension        string    // расширение
	Size             int64     // размер
}

// SaveFiles - временное хранилище.
type SaveFiles struct {
	Chunks          map[string]TmpFile // список файлов.
	tmpFileLifeTime time.Duration      // время жизни временного файла.
	muMap           sync.Mutex         // мьютекс.
	defaultPath     string             // путь по умолчанию.
	fileSave        bool               // флаг сохранения файла.
}

const defaultTmpPath = "./tmp"

// NewSaveFiles - конструктор временного хранилища.
func NewSaveFiles(tmpFileLifeTime time.Duration) (*SaveFiles, error) {

	err := os.RemoveAll(defaultTmpPath)
	if err != nil {
		return nil, err
	}

	return &SaveFiles{
		Chunks:          make(map[string]TmpFile),
		tmpFileLifeTime: tmpFileLifeTime,
		defaultPath:     defaultTmpPath,
	}, nil
}

// addNewSaveFile - добавление нового временного файла.
func (s *SaveFiles) addNewSaveFile(pathSave string, r *http.Request) (bool, *TmpFile, error) {

	file, header, err := r.FormFile("file")
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Error r.FormFile from Request.")
		return false, nil, err
	}
	defer file.Close()

	rangeMin, rangeMax, totalSize, err := ParserContentRange(r.Header.Get("Content-Range"))
	if err != nil {
		return false, nil, err
	}

	if rangeMin != 0 {
		err = customErrors.NewCustomError(nil, http.StatusBadRequest, "RangeMin must be 0 for first request.")
		return false, nil, err
	}
	if rangeMax > totalSize {
		err = customErrors.NewCustomError(nil, http.StatusBadRequest, "RangeMax more than totalSize.")
		return false, nil, err
	}
	uuId := s.generateUuid()
	pathSave = path.Join(s.defaultPath, pathSave)

	// Создание локальной папки для хранения временных файлов пользователя.
	if err := os.MkdirAll(pathSave, 0755); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error creating directory.")
		return false, nil, err
	}

	pathFileSave := path.Join(pathSave, uuId+path.Ext(header.Filename))

	s.muMap.Lock()
	tmpFile := TmpFile{
		LastUpdate:       time.Now(),
		OriginalFileName: header.Filename,
		PathFileSave:     pathFileSave,
		Uuid:             uuId,
	}

	s.Chunks[uuId] = tmpFile

	fileSize, err := s.saveFile(pathFileSave, file)
	if err != nil {
		s.muMap.Unlock()
		return false, nil, err
	}

	s.muMap.Unlock()
	fileUpload, err := s.FileUploadCompleted(fileSize, r)
	if err != nil {
		return false, nil, err
	}

	return fileUpload, &tmpFile, err
}

// UploadFile - загрузка нового временного файла.
func (s *SaveFiles) UploadFile(additionalPath string, r *http.Request) (bool, *TmpFile, error) {

	uuId := r.Header.Get("Uuid-Chunk")

	if uuId == "" {
		return s.addNewSaveFile(additionalPath, r)
	}
	s.muMap.Lock()
	defer s.muMap.Unlock()
	chunk, ok := s.Chunks[uuId]

	if !ok {
		err := customErrors.NewCustomError(nil, http.StatusBadRequest, "Error during file not found in map chunk.")
		return false, nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		err := customErrors.NewCustomError(err, http.StatusBadRequest, "Error r.FormFile from Request.")
		return false, nil, err
	}
	defer file.Close()

	// Сохраняем файл
	fileSize, err := s.saveFile(chunk.PathFileSave, file)
	if err != nil {
		return false, nil, err
	}

	// Проверяем полностью ли загружен файл
	fileUpload, err := s.FileUploadCompleted(fileSize, r)
	if err != nil {
		return false, nil, err
	}

	chunk.Size = fileSize

	chunk.LastUpdate = time.Now()
	s.Chunks[uuId] = chunk

	return fileUpload, &chunk, nil
}

// DeleteFile - удаление временного файла.
func (s *SaveFiles) DeleteFile(uuID string) error {
	f, ok := s.Chunks[uuID]
	if !ok {
		err := customErrors.NewCustomError(nil, http.StatusBadRequest, "Error during file not found in map chunk.")
		return err
	}
	// Если файл скачен то не надо удалять, тк перемещен
	if !s.fileSave {
		if err := os.Remove(f.PathFileSave); err != nil {
			err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error during delete file.")
			return err
		}
	}
	delete(s.Chunks, uuID)

	return nil
}

// RunCronDeleteFiles - cron job удаления временных файлов.
func (s *SaveFiles) RunCronDeleteFiles(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(s.tmpFileLifeTime)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cronDeleteFiles(ctx)
			}
		}
	}()
	return nil
}

// cronDeleteFiles - cron job удаления временных файлов.
func (s *SaveFiles) cronDeleteFiles(ctx context.Context) {

	s.muMap.Lock()
	defer s.muMap.Unlock()
	for k, v := range s.Chunks {

		if time.Since(v.LastUpdate) >= s.tmpFileLifeTime {

			err := s.DeleteFile(k)
			if err != nil {
				return
			}

		}
	}

}

// generateUuid - генерация уникального имени для временного файла.
func (s *SaveFiles) generateUuid() string {
	uuid := uuid.New()

	s.muMap.Lock()
	defer s.muMap.Unlock()
	if _, ok := s.Chunks[uuid.String()]; ok {
		return s.generateUuid()
	}

	return uuid.String()
}

// saveFile - сохранение файла.
func (s *SaveFiles) saveFile(pathFileSave string, file io.Reader) (int64, error) {

	f, err := os.OpenFile(pathFileSave, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error during create tmp file.")
		logger.Log.Error("Error during create tmp file.", zap.Error(err), zap.String("pathFileSave", pathFileSave))
		os.Remove(pathFileSave)
		return 0, err
	}
	defer f.Close()

	// Сохранение чанка в файл.
	if _, err := io.Copy(f, file); err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error during copy tmp file.")

		// Удаление остаточных файлов.
		os.Remove(pathFileSave)
		return 0, err
	}
	fileInfo, err := f.Stat()
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusInternalServerError, "Error during stat tmp file.")
		logger.Log.Error("Error during stat tmp file.", zap.Error(err), zap.String("pathFileSave", pathFileSave))
		// Удаление остаточных файлов.
		os.Remove(pathFileSave)
		return 0, err
	}

	if fileInfo.Size() > 52_428_800 {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Too large file size. File more than 50MB")
		return 0, err
	}
	filesize := fileInfo.Size()

	return filesize, nil
}

// ParserContentRange - парсер Content-Range.
func ParserContentRange(contentRangeHeader string) (int, int, int, error) {

	r := regexp.MustCompile(`bytes (\d+)-(\d+)/(\d+)`)

	match := r.FindStringSubmatch(contentRangeHeader)
	if len(match) != 4 {
		err := customErrors.NewCustomError(nil, http.StatusBadRequest, "Error parser Content-Range.Incorrect Content-Range header.")
		logger.Log.Error("Error parser Content-Range.Incorrect Content-Range header.", zap.Error(err), zap.String("contentRangeHeader", contentRangeHeader))

		return 0, 0, 0, err
	}
	rangeMin, err := strconv.Atoi(match[1])
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Could not parse range min from header.")
		logger.Log.Error("Could not parse range min from header.", zap.Error(err), zap.String("contentRangeHeader", contentRangeHeader), zap.String("rangeMin", match[1]))
		return 0, 0, 0, err
	}
	rangeMax, err := strconv.Atoi(match[2])
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Could not parse range max from header.")
		logger.Log.Error("Could not parse range max from header.", zap.Error(err), zap.String("contentRangeHeader", contentRangeHeader), zap.String("rangeMax", match[2]))
		return 0, 0, 0, err
	}
	totalFileSize, err := strconv.Atoi(match[3])
	if err != nil {
		err = customErrors.NewCustomError(err, http.StatusBadRequest, "Could not parse file size from header.")
		logger.Log.Error("Could not parse file size from header.", zap.Error(err), zap.String("contentRangeHeader", contentRangeHeader), zap.String("totalFileSize", match[3]))
		return 0, 0, 0, err
	}

	return rangeMin, rangeMax, totalFileSize, nil
}

// FileUploadCompleted - Проверка полной загрузки файла  по запросу.
func (s *SaveFiles) FileUploadCompleted(realFileSize int64, r *http.Request) (bool, error) {

	contentRangeHeader := r.Header.Get("Content-Range")

	_, rangeMax, fileSize, err := ParserContentRange(contentRangeHeader)
	if err != nil {
		return false, err
	}

	if fileSize == rangeMax && realFileSize == int64(rangeMax) {
		s.fileSave = true
		return true, nil
	}

	if realFileSize > int64(rangeMax) {
		err = customErrors.NewCustomError(nil, http.StatusBadRequest, "Incorect fileSize.")
		logger.Log.Error("Incorect fileSize.", zap.Error(err), zap.Int64("realFileSize", realFileSize), zap.Int("rangeMax", rangeMax))

		return false, err
	}

	return false, nil
}
