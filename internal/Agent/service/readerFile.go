package service

import (
	"GophKeeper/pkg/logger"
	"fmt"
	"go.uber.org/zap"
	"math"
	"os"
)

var ErrSize = fmt.Errorf("Error size must be > 0 and < 50mb") // Ошибка если размер больше 50mb.
var ErrOpenFile = fmt.Errorf("Error open file")               // Ошибка если не получилось открыть файл.
// Reader - читатель файла.
type Reader struct {
	SizeChunk int      // Размер чанка в байтах.
	Path      string   // Путь к файлу.
	f         *os.File // Указатель на файл.
	size      int64    // Размер файла в байтах.
	maxChunk  int      // Количество чанков в файле.
	NameFile  string   // Имя файла.
}

// NewReader - конструктор читателя файла.
func NewReader(path string) *Reader {
	return &Reader{Path: path, SizeChunk: 1024 * 1024}
}

// NumChunk - определение количества чанков в файле.
func (r *Reader) NumChunk() (int, error) {
	f, err := os.OpenFile(r.Path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logger.Log.Error("Error open file", zap.String("file", r.Path), zap.Error(ErrOpenFile))
		return 0, ErrOpenFile
	}
	r.f = f

	fileInfo, err := f.Stat()
	if err != nil {
		logger.Log.Error("Error get file info", zap.String("file", r.Path), zap.Error(err))
		return 0, err
	}
	if fileInfo.Size() > 52_428_800 || fileInfo.Size() == 0 {
		return 0, ErrSize
	}
	r.size = fileInfo.Size()
	// Определяем количество чанков
	n := math.Ceil(float64(fileInfo.Size()) / float64(r.SizeChunk))

	r.maxChunk = int(n)
	r.NameFile = fileInfo.Name()
	return r.maxChunk, nil
}

// ReadFile - считывание файла.
func (r *Reader) ReadFile(numChunk int) ([]byte, error) {
	var data []byte

	if numChunk == r.maxChunk {
		data = make([]byte, r.size%int64(r.SizeChunk))
	} else {
		data = make([]byte, r.SizeChunk)
	}

	_, err := r.f.ReadAt(data, int64(r.SizeChunk*(numChunk-1)))
	if err != nil {
		logger.Log.Error("Error read file", zap.String("file", r.Path), zap.Int("numChunk", numChunk), zap.Int("sizeChunk", r.SizeChunk), zap.Error(err))
		return nil, err
	}

	return data, nil
}

// Size - размер файла.
func (r *Reader) Size() int64 {
	return r.size
}
