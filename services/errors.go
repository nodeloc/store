package services

// ServiceError 服务错误
type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
