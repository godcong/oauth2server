# Auth2接口文档

### 配置域名：oa.mana.com
### 测试域名：oatest.mana.com

# Step 1. 登陆授权
## 跳转到：/authorize  

### 参数：  
>       必填 response_type : "code"   //此处填写code  
        必填 client_id        //玛娜花园的授权客户端id  
        选填 state            //状态码  
        选填 redirect_uri     //回调地址，必须与服务端预留地址一致  
### 返回值：  
>    code  //用于获取token
     state //返回请求的state值
### 说明：  
客户端跳转到授权界面获取登陆授权。授权成功服务器跳转至客户端注册的redirect url地址，并附带参数code，state

# Step 2. 获取token
## 地址：/token

### 协议：POST  
### 参数：  
>       必填 client_id        //玛娜花园的授权客户端id  
        必填 client_secret    //玛娜花园预留安全码
        必填 grant_type       //此处填authorization_code
        必填 code             //step1返回的code
        必填 redirect_uri     //回调地址，必须与服务端预留地址一致  
### 返回值：    
>    access_token           //用户访问授权码，通过该授权码获取内容
     token_type             //返回bearer
     refresh_token          //访问授权码过期时，通过该授权码刷新访问授权码
     expires_in             //过期时间
    
# Step 3. 获取信息
## 地址：/userinfo

### 协议：GET  
>
### 参数：  
>       必填 access_token //用户访问授权码，通过该授权码获取内容
### 返回值：  
>     "sub"               //用户唯一识别ID，同openid         
      "nickname":         //用户昵称
      "name":             //用户名
      "phone_number":     //手机号
      "email":            //邮箱

    
