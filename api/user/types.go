package user
//这里面是请求响应的一些类型
//发个请求比如说用户登录，请求的时候传什么参数，相应的时候返回什么参数等定义成结构体
//创建用户请求结构体
type CreateUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
}

//创建用户响应
type CreateUserResponse struct {
	Username string `json:"username"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 登录响应
type LoginResponse struct {
	Username string `json:"username"`
	UserId   uint   `json:"userId"`
	Token    string `json:"token"`
}
