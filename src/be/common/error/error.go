package error

type SFError struct {
	msg string
}

func (e *SFError) Error() string {
	return e.msg
}

func New(msg string) *SFError {
	return &SFError{msg: msg}
}

func DBError() *SFError {
	return New("数据库错误")
}

func AuthError() *SFError {
	return New("认证失败")
}

func RestError() *SFError {
	return New("REST交互失败")
}

func HandleRequestError() *SFError {
	return New("服务异常")
}
