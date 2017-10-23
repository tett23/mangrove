package storage_balancer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"syscall"

	"github.com/pkg/errors"
	"github.com/tett23/mangrove/assets"
	"github.com/tett23/mangrove/lib/mangrove_environment"
	yaml "gopkg.in/yaml.v2"
)

// Storage 個々の格納容器
type Storage struct {
	Name       string
	Path       string
	DiskStatus DiskStatus
}
type Storages []Storage

// DiskStatus ディスク状態
type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

// StorageConfig 設定ファイル
type StorageConfig struct {
	Storages StorageConfigItems `yaml:"storages"`
}

type StorageConfigItem struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
type StorageConfigItems []StorageConfigItem

const configFile = "config/storage.yml"

// LoadStorages ストレージ読みこみ
func LoadStorages() (Storages, error) {
	var ret Storages

	conf, err := loadStorageConfig(mangrove_environment.Get())
	if err != nil {
		return nil, errors.Wrap(err, "loadStorageConfig initialize error")
	}

	for _, item := range conf.Storages {
		s := Storage{
			Name: item.Name,
			Path: item.Path,
		}

		if err := s.UpdateDiskStatus(); err != nil {
			return ret, err
		}

		ret = append(ret, s)
	}

	return ret, nil
}

func loadStorageConfig(env string) (*StorageConfig, error) {
	bytes, err := assets.Asset(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "loadStorageConfig load asset")
	}

	var enviroments map[string]StorageConfig

	err = yaml.Unmarshal(bytes, &enviroments)
	if err != nil {
		return nil, errors.Wrap(err, "loadStorageConfig yaml.Unmarshal")
	}

	ret, ok := enviroments[env]
	if !ok {
		errors.Errorf("storage_barancer.loadStorageConfig not found %s", env)
	}

	return &ret, nil
}

// UpdateDiskStatus ディスクの状態取得
func (s *Storage) UpdateDiskStatus() error {
	status, err := diskFree(s.Path)
	if err != nil {
		return nil
	}
	s.DiskStatus = *status

	return nil
}

// Write 書きこむ
func (ss *Storages) Write(path string, data []byte) (*Storage, error) {
	if len(*ss) == 0 {
		return nil, errors.Errorf("storage_balancer.Storages.Write ss(len) == 0")
	}

	sort.Slice(ss, func(i, j int) bool {
		return (*ss)[i].DiskStatus.Free < (*ss)[j].DiskStatus.Free
	})

	s, err := ss.writableStorage(data)
	if err != nil {
		return s, err
	}

	if err = s.Write(path, data); err != nil {
		return s, err
	}
	if err = s.UpdateDiskStatus(); err != nil {
		return s, err
	}

	return nil, nil
}

func (ss *Storages) writableStorage(data []byte) (*Storage, error) {
	s := &(*ss)[0]
	if !s.haveEnoghDiskSpace(uint64(len(data))) {
		return nil, errors.Errorf("storage_balancer.Storages.writableStorage DiskStatus: %+v data: %+v", s, data)

	}

	return s, nil
}

func (s *Storage) haveEnoghDiskSpace(size uint64) bool {
	return s.DiskStatus.Free > size
}

const fileCreateMode = 0644
const directoryCreateMode = 0755

// Write ストレージに書きこみ
func (s *Storage) Write(path string, data []byte) error {
	abs := filepath.Join(s.Path, path)
	dir := filepath.Dir(abs)

	if err := os.MkdirAll(dir, directoryCreateMode); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, data, fileCreateMode); err != nil {
		return err
	}

	return nil
}

// Move ファイルを移動
func (s *Storage) Move(srcStorage *Storage, path string) error {
	srcPath := filepath.Join(srcStorage.Path, path)
	fileInfo, err := os.Stat(srcPath)
	if err != nil {
		return errors.Errorf("storage_barancer.Storage.Move s=%+v srcStorege=%+v path=%s", s, srcStorage, path)
	}
	if !s.haveEnoghDiskSpace(uint64(fileInfo.Size())) {
		return errors.Errorf("storage_barancer.Storage.Move s=%+v srcStorege=%+v path=%s", s, srcStorage, path)
	}

	destPath := filepath.Join(s.Path, path)
	if err := os.Rename(srcPath, destPath); err != nil {
		return err
	}

	return nil
}

func diskFree(path string) (*DiskStatus, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return nil, err
	}

	ret := DiskStatus{
		All:  fs.Blocks * uint64(fs.Bsize),
		Free: fs.Bfree * uint64(fs.Bsize),
	}
	ret.Used = ret.All - ret.Free

	return &ret, nil
}
