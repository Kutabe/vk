# vk
vk is a golang package that provides tools to interact with the API of the social network [VK](http://vk.com)
## Install
```
go get github.com/Kutabe/vk
```
## Usage example
```
package main
 
import (
    "fmt"
 
    "github.com/Kutabe/vk"
)
 
func main() {
    login := "me@mail.org"   //your VK login
    password := "mypassword" // your VK password
 
    user, err := vk.Auth(login, password)
    if err != nil {
        panic(err)
    }
 
    if user.Error != "" {
        fmt.Println(user.ErrorDescription)
    } else {
        parameters := make(map[string]string)
        parameters["user_id"] = "123456789"              // receiver's user ID
        parameters["message"] = "Greetings from Golang!" // message
        parameters["version"] = "5.59"                   // VK API version
 
        resp, err := vk.Request("messages.send", parameters, user)
        if err != nil {
            panic(err)
        }
        fmt.Printf("%s", resp) //This should be your message id
    }
}
```

## Used in
- [vkgetmusic](https://github.com/wingrime/vkgetmusic)

