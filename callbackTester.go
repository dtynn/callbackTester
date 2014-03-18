package main

import (
    "os"
    "fmt"

    "encoding/json"

    "github.com/dtynn/tornago"
    "github.com/qiniu/api/rs"
    . "github.com/qiniu/api/conf"
)

type Config struct {
    Host       string
    Bucket     string
    AccessKey  string
    SecretKey  string
}

var (
    Host         string
    Bucket       string
    CallbackUrl  string
    CallbackBody string
)

func loadConfig() (err error) {
    r, err := os.Open("conf.json")
    if err != nil {
        return
    }
    conf := Config{}
    d := json.NewDecoder(r)
    err = d.Decode(&conf)
    ACCESS_KEY = conf.AccessKey
    SECRET_KEY = conf.SecretKey
    Bucket = conf.Bucket
    Host = conf.Host
    CallbackUrl = Host + "/callback"
    return
}

func makeToken() (token string) {
    policy := rs.PutPolicy{}
    policy.Scope = Bucket + ":testCallback"
    policy.CallbackUrl = CallbackUrl
    policy.CallbackBody = CallbackBody
    token = policy.Token(nil)
    return
}

var form = `
<html>
 <body>
  <form method="post" action="http://up.qiniu.com/" enctype="multipart/form-data">
   <input name="token" type="hidden" value="%s">
   <input name="key" type="hidden" value="testCallback">
   File:<input name="file" type="file"/><br>
   <input type="submit" value="Upload">
  </form>
 </body>
</html>
`

func UpHdl(hdl *tornago.RequestHandler) {
    token := makeToken()
    content := fmt.Sprintf(form, token)
    hdl.Output(200, []byte(content))
}

func CallbackHdl(hdl *tornago.RequestHandler) {
    body := hdl.GetRawBody()
    auth := hdl.GetHeader("Authorization")
    fmt.Println("Authorization:", auth)
    fmt.Println("Body:", body)
    hdl.OutputJson(200, "success")
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Need a callbackBody")
        os.Exit(1)
    }
    CallbackBody = os.Args[1]
    err := loadConfig()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    routerConf := tornago.Config{
        Listen: "0.0.0.0:51234",
    }
    r := tornago.NewRouter(routerConf)
    r.HandlerFunc("POST", "/callback", CallbackHdl)
    r.HandlerFunc("GET", "/upload", UpHdl)
    err = r.Run()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
