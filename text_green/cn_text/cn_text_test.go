package cn_text

import "testing"

func TestFilterCnText(t *testing.T) {
	cases := [][]string{
		{"how are you?", ""},   //英文
		{"大傻逼", "大傻逼"},         //中文
		{"you are 大傻逼", "大傻逼"}, //中英文混合
		{"じ☆ひě綪勼ポ象緉嗰拉ポ象皮筋děん，受饬dě樬遈вμ愿魴手dě那嗰", "綪勼象緉嗰拉象皮筋受饬樬遈愿魴手那嗰"}, //英文字符，日文都去掉
		{"≥▽≤哈(◕‿◕✿)☪ ☣ ☢◐ ◑ 哈————————", "哈哈"},                       //emoji和标点符号
		{"≥▽≤哈(◕‿◕✿)☪ ☣ ☢◐ ◑ 哈————————", "哈哈"},                       //emoji和标点符号
		{"習近平", "習近平"}, //繁体
	}

	for _, cc := range cases {
		input := cc[0]
		expect := cc[1]
		out := FilterChineseText(input)

		if expect != out {
			t.Fatalf("expect:%v, got:%v", expect, out)
		}

	}

}
