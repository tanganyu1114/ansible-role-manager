package inventory

import (
	"fmt"
	"github.com/tanganyu1114/ansible-role-manager/config"
	"os"
	"path/filepath"
	"sync"
)

var (
	singletonInventoryStorageIns InventoryStorage
	onceForInventoryStorageIns   = sync.Once{}
)

type InventoryStorage interface {
	Load() (Inventory, error)
	Save(inv Inventory) error
}

type inventoryFileStorage struct {
	dir    string
	parser InventoryFileParser
}

func newInventoryFileStorage(dirPath string, parser InventoryFileParser) InventoryStorage {
	storage := &inventoryFileStorage{
		dir:    dirPath,
		parser: parser,
	}
	return InventoryStorage(storage)
}

func GetSingletonInventoryStorageIns() InventoryStorage {
	onceForInventoryStorageIns.Do(func() {
		if singletonInventoryStorageIns == nil {
			dirPath := filepath.Join(config.ExtConfig.Ansible.BaseDir, config.ExtConfig.Ansible.InventoryDir)
			parser := NewInventoryFileParser()
			singletonInventoryStorageIns = newInventoryFileStorage(dirPath, parser)
		}
	})
	return singletonInventoryStorageIns
}

func (i inventoryFileStorage) Load() (Inventory, error) {
	err := i.checkDir()
	if err != nil {
		return nil, err
	}

	filePaths, err := i.findFiles()
	if err != nil {
		return nil, err
	}
	groups := make(map[string]Group)
	for _, filePath := range filePaths {
		b, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("read inventory file %s failed, cased by: %s\n", filePath, err)
		}
		g, err := i.parser.Parse(b)
		if err != nil {
			fmt.Printf("parse inventory file %s failed, cased by: %s\n", filePath, err)
			continue
		}
		if _, has := groups[g.GetName()]; !has {
			groups[g.GetName()] = g
			continue
		}
		_ = groups[g.GetName()].addHost(g.GetHosts()...)
	}

	inv := newInventory(groups)
	return inv, nil
}

func (i inventoryFileStorage) Save(inv Inventory) error {
	err := i.checkDir()
	if err != nil {
		return err
	}
	for gName, g := range inv.GetGroups() {
		b, err := i.parser.Dump(g)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(i.dir, gName), b, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i inventoryFileStorage) findFiles() ([]string, error) {
	filePaths, err := filepath.Glob(filepath.Join(i.dir, "*"))
	if err != nil {
		return nil, err
	}

	for j := 0; j < len(filePaths); j++ {
		fStat, err := os.Stat(filePaths[j])
		if err == nil && !fStat.IsDir() {
			continue

		}
		if err != nil {
			fmt.Printf("check %s error, %s\n", filePaths[j], err)
		}
		filePaths = append(filePaths[:j], filePaths[j+1:]...)
	}

	//if len(filePaths) == 0 {
	//	err = errors.New("no matching to inventory file")
	//}

	return filePaths, err
}

func (i inventoryFileStorage) checkDir() error {
	dirStat, err := os.Stat(i.dir)
	if err != nil {
		return err
	}
	if !dirStat.IsDir() {
		return fmt.Errorf("%s is not a directory", i.dir)
	}
	return nil
}
