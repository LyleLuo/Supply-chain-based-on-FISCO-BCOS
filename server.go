package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func main() {
	dbInit()
	// Init()
	r := gin.Default()
	// r.LoadHTMLGlob("templates/*")

	//允许跨域访问
	r.Use(CrosHandler())

	//设置分组路由
	v1 := r.Group("/user")

	//根据分组设置路由
	{
		v1.POST("/login", login)
		v1.POST("/register", register)
		v1.GET("/self", self)
		v1.POST("/logout", logout)
		v1.POST("/trans", trans)
	}

	// r.GET("/", Home)
	//启动
	r.Run() // listen and serve on 0.0.0.0:8080

}

func dbInit() {
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists([]byte("account"))
		if err1 != nil {
			log.Fatal(err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		_, err2 := tx.CreateBucketIfNotExists([]byte("address"))
		if err2 != nil {
			log.Fatal(err)
		}
		return nil
	})
	db.Close()
	fmt.Println(dbSearch("account", "bank"))
	if dbSearch("account", "bank") == "" {
		dbInsert("account", "bank", "bank")
		dbInsert("address", "bank", "0xa49a7036e0eeb1190918798b446c8be159b0b8bc")
	}
}

func dbInsert(bucket string, key string, value string) {
	db, err := bolt.Open("blockchain.db", 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err = b.Put([]byte(key), []byte(value))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
}

func dbSearch(bucket string, username string) string {
	db, err := bolt.Open("blockchain.db", 0600, nil)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	var result []byte
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		result = b.Get([]byte(username))
		return nil
	})
	return string(result)
}

var (
	Secret     = "blockchain"
	ExpireTime = 3600
)

type JWTClaims struct {
	jwt.StandardClaims
	Password string `json:"password"`
	UserName string `json:"username"`
	Address  string `json:"address"`
}

func getToken(claims *JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func verifyToken(strToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, err
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}

func get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func post(url string, data []byte, contentType string) string {
	req, err := http.NewRequest(`POST`, url, bytes.NewBuffer(data))
	req.Header.Add(`content-type`, contentType)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}

type registerModel struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func register(c *gin.Context) {
	var registerInfo registerModel
	c.Bind(&registerInfo)

	if registerInfo.UserName == "bank" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "cannot register bank",
		})
		return
	}

	if dbSearch("account", registerInfo.UserName) != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "username already exists",
		})
	} else {
		dbInsert("account", registerInfo.UserName, registerInfo.Password)
		chainReturn := get("http://localhost:5002/WeBASE-Front/privateKey?type=0&userName=" + registerInfo.UserName)
		var data map[string]string
		json.Unmarshal([]byte(chainReturn), &data)
		dbInsert("address", registerInfo.UserName, data["address"])
		claims := &JWTClaims{
			UserName: registerInfo.UserName,
			Password: registerInfo.Password,
			Address:  data["address"],
		}
		claims.IssuedAt = time.Now().Unix()
		claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
		signedToken, _ := getToken(claims)
		c.SetCookie("jwt-token", signedToken, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	}
}

type loginModel struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func logout(c *gin.Context) {
	c.SetCookie("jwt-token", "", 0, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func login(c *gin.Context) {
	var loginInfo loginModel
	c.Bind(&loginInfo)
	status := "not defined"
	password := dbSearch("account", loginInfo.UserName)
	if password != loginInfo.Password {
		//如果未查询到对应字段则...
		status = "not found"
	} else {
		status = "success"
	}

	claims := &JWTClaims{
		UserName: loginInfo.UserName,
		Password: loginInfo.Password,
		Address:  dbSearch("address", loginInfo.UserName),
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()
	signedToken, _ := getToken(claims)
	c.SetCookie("jwt-token", signedToken, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

func self(c *gin.Context) {
	strToken, err := c.Cookie("jwt-token")
	claims, err := verifyToken(strToken)
	if err != nil {
		c.String(401, err.Error())
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := getToken(claims)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.SetCookie("jwt-token", signedToken, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"username": claims.UserName,
		"address":  claims.Address,
	})
}

type transModel struct {
	FuncName  string      `json:"funcName"`
	FuncParam interface{} `json:"funcParam"`
}

func trans(c *gin.Context) {
	var transInfo transModel
	c.Bind(&transInfo)

	strToken, err := c.Cookie("jwt-token")
	claims, err := verifyToken(strToken)
	if err != nil {
		c.String(401, err.Error())
		return
	}
	claims.ExpiresAt = time.Now().Unix() + (claims.ExpiresAt - claims.IssuedAt)
	signedToken, err := getToken(claims)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.SetCookie("jwt-token", signedToken, 3600, "/", "", false, true)

	abiString := "[{\"constant\":false,\"inputs\":[],\"name\":\"get_total_receipt\",\"outputs\":[{\"name\":\"total_receipt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"bill_id\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"create_bill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"companies\",\"outputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"asset\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"getAsset\",\"outputs\":[{\"name\":\"asset\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"get_total_debt\",\"outputs\":[{\"name\":\"total_debt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bank\",\"outputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"addr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"finance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"issue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"asset\",\"type\":\"uint256\"}],\"name\":\"register_company\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from_id\",\"type\":\"uint256\"},{\"name\":\"to_id\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer_bill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"bank_info\",\"outputs\":[{\"name\":\"bank_name\",\"type\":\"string\"},{\"name\":\"bank_address\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"bill_id\",\"type\":\"uint256\"}],\"name\":\"getReceipt\",\"outputs\":[{\"name\":\"borrer\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"bill_id\",\"type\":\"uint256\"}],\"name\":\"getDebt\",\"outputs\":[{\"name\":\"creditor\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Debt_query\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Receipt_query\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"asset\",\"type\":\"uint256\"}],\"name\":\"Register_company\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"account\",\"type\":\"uint256\"}],\"name\":\"Issue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Create_bill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"from_id\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"to_id\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Transfer_bill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Finance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"bill_id\",\"type\":\"uint256\"}],\"name\":\"Settle\",\"type\":\"event\"}]"
	var tempABI interface{}
	err = json.Unmarshal([]byte(abiString), &tempABI)
	if err != nil {
		fmt.Println(err)
	}

	toChain := make(map[string]interface{})
	toChain["user"] = claims.Address
	toChain["contractName"] = "supply_chain"
	toChain["contractAddress"] = "0xdcdf32d05308d16bbd84564b7b9a5818ff7c4256"
	toChain["contractAbi"] = tempABI
	toChain["groupId"] = "1"
	toChain["funcName"] = transInfo.FuncName
	toChain["funcParam"] = transInfo.FuncParam

	toChainData, _ := json.Marshal(toChain)
	txReturn := post("http://localhost:5002/WeBASE-Front/trans/handle", toChainData, "application/json")

	tempToDecode := make(map[string]string)
	err = json.Unmarshal([]byte(txReturn), &tempToDecode)
	if err != nil {
		fmt.Println(err)
	}

	if tempToDecode["status"] != "0x0" {
		c.JSON(http.StatusForbidden, gin.H{
			"result": "forbidden",
		})
		return
	}

	toChainDecode := make(map[string]interface{})
	toChainDecode["input"] = tempToDecode["input"]
	toChainDecode["output"] = tempToDecode["output"]
	toChainDecode["abiList"] = tempABI
	toChainDecode["decodeType"] = 2
	toChainDecode["returnType"] = 2
	toChainDecodeString, _ := json.Marshal(toChainDecode)

	decodeData := post("http://localhost:5002/WeBASE-Front/tool/decode", toChainDecodeString, "application/json")

	tempResult := make(map[string]interface{})
	err = json.Unmarshal([]byte(decodeData), &tempResult)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": tempResult["result"],
	})
}

//CrosHandler 跨域访问：cross  origin resource share
func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json") //设置返回格式是json

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "OK")
		}

		//处理请求
		context.Next()
	}
}
