package errorkit

var (
	SensitiveWordsCodeError  = NewStructuredError().SetCode(506).SetMessage("您输入的内容包括敏感词语")
	KeepSilentCodeErrorError = NewStructuredError().SetCode(1001).SetMessage("您已被禁言")

	MissingAccessTokenError   = NewHttpError(401).SetCode(2001).SetMessage("missing accesstoken")
	UserHasNoPermissionForBot = NewHttpError(401).SetCode(2002).SetMessage("没有机器人调用权限")
)
