# go-iliad-it-sdk

> The unofficial Golang SDK to access [`iliad.it`](https://www.iliad.it/) services

## Usage
```go
iClient := iliad.NewClient()

tkn, err := iClient.Login("your-iliad-it-userid", "your-iliad-it-password")
if err != nil {
  log.Printf("Error during login: %s", err)
  return
}

log.Println("User token:", tkn)

usr, usrErr := iClient.GetUserInfo()
if usrErr != nil {
  log.Printf("Error during GetUserInfo: %s", usrErr)
  return
}

log.Printf("UserInfo: %+v", usr)

crdInfo, crdErr := iClient.GetUserCreditInfo()
if crdErr != nil {
  log.Printf("Error during GetUserCreditInfo: %s", crdErr)
  return
}

log.Printf("UserCreditInfo: %+v", crdInfo)
```

## Notice

This is an unofficial API developed for learning purposes. Iliad Italia S.p.A. is not responsible in any way.

This program comes with ABSOLUTELY NO WARRANTY. This is free software, and you are welcome to redistribute it.

---

Made with :sparkles: & :heart: by [Mattia Panzeri](https://github.com/panz3r) and [contributors](https://github.com/panz3r/go-iliad-it-sdk/graphs/contributors)