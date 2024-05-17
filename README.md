

git add .
git commit -m "My first release 2.0.4"
git tag -a 2.0.4 -m "My first release 2.0.4"
git push origin 2.0.4


## 添加依赖


```shell
go get -u github.com/mdp/qrterminal/v3 v3.0.0

go get -u github.com/wechaty/go-wechaty

go get -u github.com/xxl6097/go-http@v0.0.7

go get -u github.com/gorilla/mux

go get -u github.com/google/uuid

go get -u github.com/dgrijalva/jwt-go

go get -u gopkg.in/yaml.v3

```

```azure
curl -F "file=@/Users/uuxia/Desktop/work/code/go/go-upload/main.go" -F "token=55" http://localhost:5555/upload

curl -F "files=@/Users/uuxia/Desktop/work/code/go/go-upload/main.go" -F "token=het002402" http://uuxia.cn:8087/upload

curl -F "file=@/Users/uuxia/Desktop/work/code/go/go-upload/main.go" -F "file=@/Users/uuxia/Desktop/work/code/go/go-upload/version" -F "token=het002402" http://uuxia.cn:8087/upload


```

```azure
curl -F "file=@$1" -F "token=$2" https://uuxia.cn/v1/api/file/upfile

```


docker-compose.yaml

```yaml
version: "3.3"
services:
  homepage:
    image: xxl6097/go-upload:0.0.2
    restart: no
    container_name: go-upload
    ports:
      - 3008:8087
    volumes:
      - ./conf/files:/app/files
    environment:
      ENV_PORT: 8087
      ENV_TOKEN: het002402

```


curl -F "file=@./teamide" -F "token=het002402" http://uuxia.cn:8087/upload


### 指令上传示例：
```shell
curl -F "file=@/root/xxx.log" -F "token=44" http://localhost:4444/upload
```

bash <(curl -s -S -L http://uuxia.cn:8087/files/2024/03/12/test.sh)  /Users/uuxia/Desktop/work/code/go/go-upload/go.mod

```ssh
ssh -T git@github.com
```