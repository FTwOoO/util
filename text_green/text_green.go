package text_green

import (
	"github.com/rexue2019/util/text_green/cn_text"
	"github.com/rexue2019/util/text_green/word_filter"
)

var defaultTrie *word_filter.Trie

func init() {
	defaultTrie = word_filter.NewDefaultTrie()
}

func CheckText(reqtext string, blockWithAsterisk bool) (text string, resultText string, sensitiveWords []string, err error) {
	replaceChar := rune(42)
	reqtext = cn_text.FilterChineseText(reqtext)
	has, sensitiveWorlds2, resultText2 := defaultTrie.QueryAndReplace(reqtext, replaceChar)
	text = reqtext
	sensitiveWords = sensitiveWorlds2

	if text == reqtext {
		resultText = resultText2
	} else if text != reqtext && has {
		origin := reqtext
		//因为敏感词替换返回的替换结果是先过滤掉原文符号再替换的，要用回原文把敏感词一个个字符替换才正确
		for _, sub := range sensitiveWorlds2 {
			origin = FilterSubString(origin, sub, replaceChar)
		}

		resultText = origin

	} else {
		resultText = reqtext
	}
	return
}

//将text字符串中的字串sub每一个字符替换成指定字符
//sub在text中可以不连续
func FilterSubString(text string, sub string, replaceChar rune) string {
	findRunes := []rune(sub)
	currentFindIndex := 0

	runes := []rune(text)
	resultRunes := []rune{}

	for _, r := range runes {
		if r == findRunes[currentFindIndex] {
			resultRunes = append(resultRunes, replaceChar)

			currentFindIndex += 1
			if currentFindIndex >= len(findRunes) {
				//找完所有子串的字符了
				break
			}
			continue
		} else {
			resultRunes = append(resultRunes, r)
		}
	}

	return string(resultRunes)

}
