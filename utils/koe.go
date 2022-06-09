package utils

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/others"
	"github.com/mitchellh/go-homedir"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	compile = regexp2.MustCompile(`(?<!\d)\d{6}(?!\d)`, regexp2.None)
)

func FilePath2Struct[V interface{}](filePath string) []V {
	list := make([]V, 0)
	value := gjson.ParseBytes(IgnoreErr(ioutil.ReadFile(IgnoreErr(homedir.Expand(filePath)))))

	value.ForEach(func(_, v gjson.Result) bool {
		var po V
		Unmarshal(v.Raw, &po)
		list = append(list, po)
		return true
	})

	return list
}

func ListCode(s string) []string {
	regexp2FindAllString := func(re *regexp2.Regexp, s string) []string {
		var matches []string
		m := IgnoreErr(re.FindStringMatch(s))
		for m != nil {
			matches = append(matches, m.String())
			m = IgnoreErr(re.FindNextMatch(m))
		}
		return matches
	}

	codeList := hashset.New(Map(regexp2FindAllString(compile, s), Str2Any)...)
	return Map(codeList.Values(), Any2Str)
}

func ListMyCode(nodes []*others.Node) []string {
	var b strings.Builder

	for _, item := range FlatTree(nodes) {
		if item.Type == FolderType {
			b.Write([]byte(item.Title))
		}
	}

	return ListCode(b.String())
}

func GetNasJson() []others.Po {
	return FilePath2Struct[others.Po](GetNextNasCacheFile())
}

func GetNextNasCacheFile() string {
	dir := filepath.Dir(IgnoreErr(homedir.Expand(global.GetServiceContext().Config.FlagConfig.NasCacheFile)))
	files := IgnoreErr(ioutil.ReadDir(dir))
	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i].Name(), files[j].Name()) > 0
	})

	return path.Join(dir, files[0].Name())
}

// GetImgUrl 获取图片url地址
func GetImgUrl(code, typee string) string {
	code2 := code[:3] + "000"
	if IgnoreErr(strconv.Atoi(code[3:])) != 0 {
		code2 = strconv.Itoa(IgnoreErr(strconv.Atoi(code2)) + 1000)
	}

	config := global.GetServiceContext().Config

	url := fmt.Sprintf(config.DownloadPattern1, code2, code, typee)
	if typee == "240x240" || typee == "360x360" {
		url = fmt.Sprintf(config.DownloadPattern2, code2, code, typee)
	}

	return url
}

func GetLrcPath(name string, nodes []*others.Node, fn func(string) string) (string, error) {
	filter := Filter[*others.Node](FlatTree(nodes), func(item *others.Node) bool {
		return item.Type != "folder" && filepath.Ext(item.Title) == ".lrc"
	})

	lrcMap := make(map[int]string)

	for _, e := range Map[*others.Node](filter, func(item *others.Node) string {
		return fn(item.UUID)
	}) {
		lrcMap[Longest(name, e)] = e
	}

	keys := make([]int, 0, len(lrcMap))
	for k := range lrcMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	if len(keys) == 0 {
		return "", fmt.Errorf("没有找到lrc文件")
	}

	return lrcMap[keys[len(keys)-1]], nil
}
