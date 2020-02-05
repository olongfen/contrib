package i18n

import (
	"github.com/olongfen/contrib"
	"github.com/olongfen/contrib/log"
	"github.com/stretchr/testify/require"

	"testing"
)

// 测试翻译基本功能
func Test_i18n(t *testing.T) {
	as := require.New(t)
	e := contrib.ErrUndefined
	ev := e.SetVars("haha")
	log.Warnln(I18n.GetError(e, I18nZhCN))
	as.Equal("contrib26|detail: undefined", I18n.GetError(e, I18nZhCN), "取翻译后的值")
	as.Equal("contrib26|detail: undefined haha", I18n.GetError(ev, I18nZhCN), "取翻译后的值")
	as.Equal(ev.Error(), e.Error(), "取翻译后的值")
}

// 载入翻译
//func Test_i18nFile(t *testing.T) {
//	as := require.New(t)
//	num, err := I18n.LoadDir("./lang")
//	as.Nil(err)
//	log.Println(num, "aaaaaaaaaaaaaa")
//	as.Equal(true, num >= 51, "rest存量翻译")
//
//	//
//	as.Equal("rest129", contrib.ErrParamInvalid.Error())
//	as.Equal("rest129|参数非法 some val", I18n.GetError(contrib.ErrParamInvalid.SetVars("some val"), I18nZhCN))
//}

// 自定义翻译
func Test_i18nFileSetVar(t *testing.T) {
	as := require.New(t)
	_, err := I18n.LoadDir("./lang")
	as.Nil(err)
	er := contrib.ErrParamInvalid.SetVarsBy(I18nZhCN, "第一个参数非法").SetVarsBy(I18nEn, "first params")
	//t.Logf(`"%s" "%s"`, I18n.GetError(er, I18nZhCN), I18n.GetError(er, I18nEn))
	as.Equal("contrib1|detail: param invalid", I18n.GetError(er, I18nZhCN))
	as.Equal("contrib1|detail: param invalid", I18n.GetError(er, I18nEn))
	as.Equal("contrib1|detail: param invalid", I18n.GetError(er, ""))
}
