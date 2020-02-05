package i18n

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

const (
	//
	I18nEn   = "en"
	I18nZhCN = "zh-CN"
	I18nZhHK = "zh-HK"
	I18nZhTw = "zh-TW"
	//
	I18nCode = "__code__"
)

var (
	// I18nDefault 默认语言
	I18nDefault = I18nEn
	// I18n 默认翻译
	I18n = NewI18nMap()
	// 匹配正则
	regI18nKey  = regexp.MustCompile(`^([a-z]+)?(Err)?([A-Z]\w+)$`)
	regI18nCode = regexp.MustCompile(`^([a-z]+)?([0-9]+)$`)
)

// 多语言错误
type I18nErrorBase interface {
	GetCode() int32    // 错误代码
	GetPrefix() string // 模块名
	GetDetail() string // 缺省文本
}

// 多语言错误
type I18nError interface {
	I18nErrorBase
	GetVars() []interface{}              // 错误变量
	GetVarsBy(lang string) []interface{} // 错误变量
}

// 多语言翻译缓存
type I18nMap struct {
	mapTemplate map[string]map[string]*template.Template // key->lang->tpl
	mapDefault  map[string]*template.Template            // key->tpl
	mapNumVars  map[string]int                           // key->int 每个模块带的参数
	configDir   []string                                 //
	lock        sync.RWMutex                             //
}

// Parse 读取内容
func (m *I18nMap) ParseJSON(raw []byte) (ret1 int, ret2 int, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var (
		record = make(map[string]interface{}) // json内的翻译信息
		numKey int                            //
		numTxt int                            //
	)
	if err = json.Unmarshal(raw, &record); err != nil {
		return
	}

	// 解析信息
	for k, v0 := range record {
		var (
			kLis  = regI18nKey.FindStringSubmatch(k)
			model = ""
			key   = ""
			txt   = "" // 默认翻译
			code  = -1
		)
		v, ok := v0.(map[string]interface{})
		if !ok {
			continue
		}
		if len(kLis) == 0 {
			continue
		}
		model = kLis[1] // 旧版json中model的位置
		if _v, _ok := v[I18nCode]; _ok {
			if _kLis := regI18nCode.FindStringSubmatch(fmt.Sprintf("%v", _v)); len(_kLis) > 0 {
				if len(_kLis[1]) > 0 {
					model = _kLis[1] // 新版json中model的位置
				}
				if len(_kLis[2]) > 0 {
					// code
					if _v, _err := strconv.Atoi(_kLis[2]); _err == nil {
						code = _v
					}
				}
			}
			delete(v, I18nCode)
		}
		if len(model) == 0 || code < 0 {
			continue
		}
		key = fmt.Sprintf(`prefix: %s,code: %d`, model, code)

		// 多语言模版
		for _k, _v0 := range v {
			_v, _ok := _v0.(string)
			if !_ok || len(_v) == 0 {
				continue
			}
			// 转义
			switch _k {
			case "zh":
				_k = "zh-CN"
			case "cn":
				_k = "zh-CN"
			case "hk":
				_k = "zh-HK"
			case "tw":
				_k = "zh-TW"
			}

			// 默认翻译
			if _k == "en" || _k == "en-US" {
				txt = _v
			}
			if _k == "zh-CN" && len(txt) == 0 {
				txt = _v
			}

			// 存入缓存
			if m.mapTemplate[key] == nil {
				m.mapTemplate[key] = make(map[string]*template.Template)
			}
			m.mapTemplate[key][_k] = template.Must(template.New(key).Option("missingkey=zero").Parse(_v))

			// 兼容旧JSON文件
			switch _k {
			case "zh-CN":
				m.mapTemplate[key]["cn"] = m.mapTemplate[key][_k]
			case "zh-HK":
				m.mapTemplate[key]["hk"] = m.mapTemplate[key][_k]
				if m.mapTemplate[key]["tw"] == nil {
					m.mapTemplate[key]["tw"] = m.mapTemplate[key][_k]
				}
			case "zh-TW":
				m.mapTemplate[key]["tw"] = m.mapTemplate[key][_k]
				if m.mapTemplate[key]["hk"] == nil {
					m.mapTemplate[key]["hk"] = m.mapTemplate[key][_k]
				}
			case "en-US":
				m.mapTemplate[key]["en"] = m.mapTemplate[key][_k]
			}

			//
			numTxt++
		}

		// 默认模版的变量数
		m.mapNumVars[key] = strings.Index(txt, ".var")

		// 默认模版
		m.mapDefault[key] = template.Must(template.New(key).Option("missingkey=zero").Parse(txt))

		//
		numKey++
	}

	ret1 = numKey
	ret2 = numTxt
	return
}

// LoadDir 读取翻译文件
func (m *I18nMap) LoadDir(dirs ...string) (numKey int, err error) {
	if len(dirs) == 0 {
		dirs = m.configDir
	}
	m.configDir = dirs

	//
	for _, dir := range dirs {
		// walk
		if err = filepath.Walk(dir, func(fp string, info os.FileInfo, fileErr error) (err error) {
			// 忽律非".json文件"
			if path.Ext(fp) != ".json" {
				return
			}
			// 忽律文件夹
			if fileErr == nil && info.IsDir() == true {
				return
			}

			// 文件错误
			if err = fileErr; err != nil {
				return
			}

			defer func() {
				// 捕获错误
				if _err := recover(); _err != nil {
					err = fmt.Errorf(`[rest-i18n] load "%s", err: %v`, fp, _err)
				}
			}()
			// 文件错误
			if err = fileErr; err != nil {
				return
			}

			// 解析内容
			if _con, _err := ioutil.ReadFile(fp); _err == nil {
				if _ret, _, _err := m.ParseJSON(_con); _err == nil {
					numKey += _ret
				} else {
					err = _err
					return
				}
			} else {
				err = _err
				return
			}

			return
		}); err != nil {
			return
		}
	}

	return
}

// GetError 翻译错误
func (m *I18nMap) GetError(err I18nError, lang string) (ret string) {
	var (
		key  = fmt.Sprintf(`prefix: %s,code: %d`, err.GetPrefix(), err.GetCode())
		vars = err.GetVarsBy(lang) // 优先取多语言变量
		txt  = err.GetDetail()
		tpl  = m.getTemplate(key, m.ParseLang(lang), len(vars), txt)
	)
	if len(vars) == 0 {
		vars = err.GetVars()
	}
	if tpl != nil {
		// 填充变量
		var (
			b  bytes.Buffer
			f  = bufio.NewWriter(&b)
			ps = make(map[string]string)
		)
		for _idx, _v := range vars {
			ps[fmt.Sprintf("var%d", _idx+1)] = fmt.Sprintf("%v", _v)
		}
		if _err := tpl.Execute(f, ps); _err == nil {
			_ = f.Flush()
			txt = b.String()
		}
	}

	ret = strings.TrimSpace(key + "," + txt)
	return
}

// GetErrorBase 翻译错误
func (m *I18nMap) GetErrorBase(err I18nErrorBase, lang string) (ret string) {
	return m.GetError(&i18nError{err}, lang)
}

// I18nLangParse 解析语言格式
func (m *I18nMap) ParseLang(lang string) (ret string) {
	if len(lang) == 0 {
		ret = I18nDefault
		return
	}

	//
	switch lang {
	case "zh":
		return I18nZhCN
	case "cn":
		return I18nZhCN
	case "hk":
		return I18nZhHK
	case "tw":
		return I18nZhTw
	case I18nEn, I18nZhCN, I18nZhHK:
		return lang
	default:
		break
	}

	// TODO: 优化语言识别
	if strings.Index(lang, "zh-CN") > -1 {
		ret = I18nZhCN
	} else if strings.Index(lang, "zh-HK") > -1 {
		ret = I18nZhHK
	} else if strings.Index(lang, "zh-TW") > -1 {
		ret = I18nZhTw
	} else if strings.Index(lang, "zh") > -1 {
		ret = I18nZhCN
	} else if strings.Index(lang, "en") > -1 {
		ret = I18nEn
	} else {
		ret = lang
	}
	return
}

//
func (m *I18nMap) getTemplate(key string, lang string, pNum int, txt string) (ret *template.Template) {
	lang = m.ParseLang(lang)
	if _v, _ok := m.mapTemplate[key]; _ok {
		ret = _v[lang]
	}
	if ret == nil {
		ret = m.mapDefault[key]
	}
	if ret == nil && len(txt) > 0 {
		m.lock.Lock()
		defer m.lock.Unlock()
		m.mapDefault[key] = template.Must(template.New(key).Option("missingkey=zero").Parse(txt))
		ret = m.mapDefault[key]
	}
	return
}

// NewI18nMap
func NewI18nMap() (ret *I18nMap) {
	ret = new(I18nMap)
	ret.mapTemplate = make(map[string]map[string]*template.Template)
	ret.mapDefault = make(map[string]*template.Template)
	ret.mapNumVars = make(map[string]int)
	return
}

// PubError 错误转换
func PubError(in error) (out error) {
	if in == nil || I18n == nil {
		return in
	}
	switch err := in.(type) {
	case I18nError:
		out = fmt.Errorf(`%s`, I18n.GetError(err, I18nDefault))
	case I18nErrorBase:
		out = fmt.Errorf(`%s`, I18n.GetErrorBase(err, I18nDefault))
	default:
		out = in
	}
	return
}

// 判断例子
type i18nError struct {
	I18nErrorBase
}

//
func (e *i18nError) GetVars() []interface{} {
	return nil
}

//
func (e *i18nError) GetVarsBy(lang string) []interface{} {
	return nil
}
