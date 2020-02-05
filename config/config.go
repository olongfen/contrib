package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/olongfen/contrib/log"

	"gopkg.in/yaml.v2"
)

// Config
type Config struct {
	sync.RWMutex
	pathYaml   *string                 // yaml配置文件保存地址
	savePoint  interface{}             //
	hookChange func(interface{}) error //
	comments   map[string]string       // 配置文件的备注信息
	silent     bool
}

// LoadConfiguration load config
func LoadConfiguration(configPath string, targetConfig, defaultConfig interface{}) (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile(configPath); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		if defaultConfig == nil {
			err = fmt.Errorf(`[Config] dedaultConfig undefined, "%s" error: %v`, configPath, err)
			return
		}
		// 自动创建配置文件
		if d, _err := yaml.Marshal(defaultConfig); _err != nil {
			err = _err
			return
		} else if err = ioutil.WriteFile(configPath, d, 0666); err != nil {
			return
		}
		err = fmt.Errorf(`[Config] please modify "%s" and run again`, configPath)
		return
	}

	if err = yaml.Unmarshal(data, targetConfig); err != nil {
		return
	}

	// 设置保存地址对对象指针
	type configInterface interface {
		SetSavePath(savePath string) (err error)
		SetSavePoint(saveTarget interface{}) (err error)
	}
	if _c, _ok := targetConfig.(configInterface); _ok == true {
		if err = _c.SetSavePath(configPath); err != nil {
			return
		}
		if err = _c.SetSavePoint(targetConfig); err != nil {
			return
		}
	}

	return
}

// Save save config
func (c *Config) Save(newConfig interface{}) (err error) {
	c.Lock()
	defer c.Unlock()
	var (
		savePath    string
		readContent []byte
		saveContent []byte
	)
	if savePath, err = c.GetSavePath(); err != nil {
		return
	}
	// 读旧记录
	readContent, _ = ioutil.ReadFile(savePath)

	if newConfig == nil {
		newConfig = c.savePoint
	}
	if saveContent, err = yaml.Marshal(newConfig); err != nil {
		return
	} else if bytes.Equal(readContent, saveContent) == true {
		// 不重复保存
		return
	}

	// 写入记录
	if err = ioutil.WriteFile(savePath, saveContent, 0666); err != nil {
		return
	}
	if c.silent == false {
		log.Println(fmt.Sprintf("[Config] save Config to %s bytes:%d->%d",
			savePath, len(readContent), len(saveContent)))
	}

	// hook
	if c.hookChange != nil {
		if _err := c.hookChange(newConfig); _err != nil {
			log.Println(fmt.Sprintf(`[Config] hookChange error: %v`, _err))
		}
	}
	return
}

// GetSavePath get path of save config
func (c *Config) GetSavePath() (ret string, err error) {
	if c.pathYaml == nil {
		err = errors.New("param invalid")
		return
	} else {
		ret = *c.pathYaml
	}
	return
}

// SetSavePath set save path
func (c *Config) SetSavePath(savePath string) (err error) {
	if len(savePath) == 0 {
		c.pathYaml = nil
		err = errors.New("param invalid")
		return
	} else {
		c.pathYaml = &savePath
	}
	return
}

// SetSavePoint set save object
func (c *Config) SetSavePoint(saveTarget interface{}) (err error) {
	c.savePoint = saveTarget
	return
}
