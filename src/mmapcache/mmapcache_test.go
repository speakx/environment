package mmapcache

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var pwd string
var template []byte
var cachesize = 1024 * 1024 * 10

func TestMmapCache(t *testing.T) {
	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	template = createMMapTemplate(cachesize)

	cachefile := fmt.Sprintf("%v/0.dat", pwd)
	t.Logf("cachefile:%v", cachefile)

	err := createMMapFile(cachefile, template)
	if nil != err {
		t.Errorf("createMMapFile failed err:%v", err)
		return
	}
	t.Logf("createMMapFile is ok")

	mmapCache, err := newMMapCache(cachefile)
	if nil != err {
		t.Errorf("newMMapCache failed err:%v", err)
		return
	}
	if len(mmapCache.buf) != cachesize {
		t.Errorf("newMMapCache size:%v err(alloc:%v)", len(mmapCache.buf), cachesize)
		return
	}
	t.Logf("newMMapCache is ok")
}
