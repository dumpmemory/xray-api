package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"xray-api/config"

	cmap "github.com/orcaman/concurrent-map"
)

var (
	Config     config.CONF
	Protocol   string
	users      = cmap.New()
	InbountTag = "vryusers"
	xrayCtl    *XrayController
)

func Init() (err error) {
	Protocol = Config.Xray.Protocol
	xrayCtl = new(XrayController)
	err = xrayCtl.Init(&BaseConfig{
		APIAddress: "127.0.0.1",
		APIPort:    uint16(Config.Xray.Grpc),
	})
	if err != nil {
		return
	}
	return
}

type User struct {
	Uuid  string `json:"uuid"`
	Email string `json:"email"`
}

func Xuser(uuid, email string) UserInfo {
	return UserInfo{
		Uuid:       uuid,
		Email:      email,
		Password:   uuid,
		AlertId:    0,
		Level:      0,
		InTag:      InbountTag,
		CipherType: Config.Xray.Method,
	}
}

func AddUser(Uuid string, Email string) (err error) {
	user := User{
		Uuid:  Uuid,
		Email: Email,
	}
	xuser := Xuser(Uuid, Email)
	if Protocol == "vmess" || Protocol == "vless" {
		err = addVmessUser(xrayCtl.HsClient, &xuser)
	} else if Protocol == "shadowsocks" {
		err = addSSUser(xrayCtl.HsClient, &xuser)
	} else if Protocol == "trojan" {
		err = addTrojanUser(xrayCtl.HsClient, &xuser)
	}

	if err == nil {
		users.Set(Uuid, user)
	} else {
		log.Println(err)
	}
	return
}
func RemoveUser(Uuid string) (err error) {
	user, has := users.Get(Uuid)
	if !has {
		return errors.New("user not found")
	}
	xuser := Xuser(Uuid, user.(User).Email)
	err = removeUser(xrayCtl.HsClient, &xuser)
	if err == nil {
		users.Remove(Uuid)
	}
	return err
}
func Sync(newUsers []User) map[string]interface{} {
	S := make(map[string]User)
	for _, user := range newUsers {
		S[user.Uuid] = user
	}
	return SyncS(S)
}
func SyncS(S map[string]User) map[string]interface{} {
	for _, Uuid := range users.Keys() {
		_, has := S[Uuid]
		if has {
			delete(S, Uuid)
		} else {
			RemoveUser(Uuid)
		}
	}
	for _, user := range S {
		AddUser(user.Uuid, user.Email)
	}
	Users := users.Items()
	if Config.Syncfile != "" {
		data, err := json.Marshal(Users)
		if err == nil {
			err = ioutil.WriteFile(Config.Syncfile, data, 0644)
		}
		if err != nil {
			log.Println(err)
		}
	}
	return Users
}
func ReAddUsers() {
	for _, user := range users.Items() {
		AddUser(user.(User).Uuid, user.(User).Email)
	}
}
func Traffic(reset bool) map[string]uint64 {
	res := make(map[string]uint64)
	eres := make(map[string]uint64)
	ptn := ""
	stats, err := queryTraffic(xrayCtl.SsClient, ptn, reset)
	if err != nil {
		return res
	}
	Reg := regexp.MustCompile(`user>>>([^>]+)>>>traffic>>>([^>]+)link`)
	for _, stat := range stats {
		match := Reg.FindStringSubmatch(stat.Name)
		if len(match) > 0 {
			eres[match[1]] += uint64(stat.Value)
		}
	}
	for uuid, user := range users.Items() {
		traffic, has := eres[user.(User).Email]
		if has && traffic > 0 {
			res[uuid] = traffic
		}
	}
	return res
}
