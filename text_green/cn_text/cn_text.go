package cn_text

type UnicodeRange [2]rune
type TypeUnicodeRanges []UnicodeRange

func (this *TypeUnicodeRanges) Contains(r rune) bool {
	for _, ur := range *this {
		if ur[0] <= r && ur[1] >= r {
			return true
		}
	}

	return false
}

func (this *TypeUnicodeRanges) Add(r ...UnicodeRange) {
	for _, x := range r {
		*this = append(*this, x)
	}
}

//reference : http://unicode.org/charts/

var CJKUnifiedIdeographs UnicodeRange = [2]rune{0x4E00, 0x9FEF}
var CJKUnifiedIdeographsExtensionA UnicodeRange = [2]rune{0x3400, 0x4DB5}
var CJKUnifiedIdeographsExtensionB UnicodeRange = [2]rune{0x20000, 0x2A6D6}
var CJKUnifiedIdeographsExtensionC UnicodeRange = [2]rune{0x2A700, 0x2B734}
var CJKUnifiedIdeographsExtensionD UnicodeRange = [2]rune{0x2B740, 0x2B81D}
var CJKUnifiedIdeographsExtensionE UnicodeRange = [2]rune{0x2B820, 0x2CEA1}
var CJKUnifiedIdeographsExtensionF UnicodeRange = [2]rune{0x2CEB0, 0x2EBE0}

var CJKCompatibilityIdeographs UnicodeRange = [2]rune{0xF900, 0xFAFF}
var CJKCompatibilityIdeographsSupplement UnicodeRange = [2]rune{0x2F800, 0x2FA1F}

var KangxiRadicals UnicodeRange = [2]rune{0x2F00, 0x2FDF}
var CJKRadicalsSupplement UnicodeRange = [2]rune{0x2E80, 0x2EFF}
var CJKStrokes UnicodeRange = [2]rune{0x31C0, 0x31EF}

var CjkUnicodeRanges TypeUnicodeRanges

func init() {

	CjkUnicodeRanges.Add(
		CJKUnifiedIdeographs,
		CJKUnifiedIdeographsExtensionA,
		CJKUnifiedIdeographsExtensionB,
		CJKUnifiedIdeographsExtensionC,
		CJKUnifiedIdeographsExtensionD,
		CJKUnifiedIdeographsExtensionE,
		CJKUnifiedIdeographsExtensionF,
		CJKCompatibilityIdeographs,
		CJKCompatibilityIdeographsSupplement,
		KangxiRadicals,
		CJKRadicalsSupplement,
		CJKStrokes,
	)

}

//过滤出中文unicode
func FilterChineseText(text string) string {
	runes := []rune(text)
	result := make([]rune, 0)
	for _, r := range runes {
		if CjkUnicodeRanges.Contains(r) {
			result = append(result, r)
		}
	}

	return string(result)
}
