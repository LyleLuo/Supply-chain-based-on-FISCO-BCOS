# 后端使用指南
打开后端服务
```
go run server.go
```
共有五个API
```go
	//根据分组设置路由
	{
		v1.POST("/login", login)
		v1.POST("/register", register)
		v1.GET("/self", self)
		v1.POST("/logout", logout)
		v1.POST("/trans", trans)
    }
```

 - `POST http://localhost:8080/user/register`
    
    接收字段
    ```json
    {
        "username":"lp",
        "password":"123"
    }
    ```
    注册成功则返回`200`和以下字段
    ```json
    {
    "status": "success"
    }
    ```
 - `POST http://localhost:8080/user/login`
    
    同上

- `GET http://localhost:8080/user/self`

    如果已登陆则返回`200`并获得当前信息（用户名，地址）
    ```json
    {
    "address": "0x7740b658bb71c99db4792dd21dc7b10852e0563a",
    "username": "lp"
    }
    ```

 - `POST http://localhost:8080/user/logout`
    
    登出

- `POST http://localhost:8080/user/trans`

    调用链上方法，并返回结果。例如：
    
    传入以下字段（函数名，参数列表）以调用链中 `function register_company(string memory name, uint asset)`
    ```json
    {
        "funcName":"register_company",
        "funcParam":["Car",1000]
    }
    ```
    结果（该函数没有返回值）
    ```json
    {
    "result": []
    }
    ```

    又如，调用链中 `function getAsset() public returns(uint asset)`
    ```json
    {
    "funcName":"getAsset",
    "funcParam":[]
    }
    ```
    返回结果asset
    ```json
    {
    "result": [
        {
        "data": 1000,
        "name": "asset",
        "type": "uint256"
        }
    ]
    }
    ```
