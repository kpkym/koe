package koe

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
	"github.com/mitchellh/go-homedir"
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

func ListCode(s string) []string {
	regexp2FindAllString := func(re *regexp2.Regexp, s string) []string {
		var matches []string
		m := utils.IgnoreErr(re.FindStringMatch(s))
		for m != nil {
			matches = append(matches, m.String())
			m = utils.IgnoreErr(re.FindNextMatch(m))
		}
		return matches
	}

	codeList := hashset.New(utils.Map(regexp2FindAllString(compile, s), utils.Str2Any)...)
	return utils.Map(codeList.Values(), utils.Any2Str)
}

func ListMyCode(nodes []*others.Node) []string {
	var b strings.Builder

	for _, item := range utils.FlatTree(nodes) {
		if item.Type == utils.FolderType {
			b.WriteString(item.Title + " ")
		}
	}

	return ListCode(b.String())
}

func GetNasJson() []others.Po {
	return utils.FilePath2Struct[others.Po](GetNextNasCacheFile())
}

func GetNextNasCacheFile() string {
	dir := filepath.Dir(utils.IgnoreErr(homedir.Expand(global.GetServiceContext().Settings.NasCacheFile)))
	files := utils.IgnoreErr(ioutil.ReadDir(dir))
	sort.Slice(files, func(i, j int) bool {
		return strings.Compare(files[i].Name(), files[j].Name()) > 0
	})

	return path.Join(dir, files[0].Name())
}

// GetImgUrl 获取图片url地址
func GetImgUrl(code, typee string) string {
	code2 := code[:3] + "000"
	if utils.IgnoreErr(strconv.Atoi(code[3:])) != 0 {
		code2 = strconv.Itoa(utils.IgnoreErr(strconv.Atoi(code2)) + 1000)
	}

	config := global.GetServiceContext().Config

	url := fmt.Sprintf(config.DownloadPattern1, code2, code, typee)
	if typee == "240x240" || typee == "360x360" {
		url = fmt.Sprintf(config.DownloadPattern2, code2, code, typee)
	}

	return url
}

func GetLrcPath(name string, nodes []*others.Node, fn func(uint32) string) (string, error) {
	filter := utils.Filter[*others.Node](utils.FlatTree(nodes), func(item *others.Node) bool {
		return item.Type != "folder" && filepath.Ext(item.Title) == ".lrc"
	})

	lrcMap := make(map[int]string)

	for _, e := range utils.Map[*others.Node](filter, func(item *others.Node) string {
		return fn(item.UUID)
	}) {
		lrcMap[utils.Longest(name, e)] = e
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
