package api

import (
	"errors"
	"regexp"
)

var (
	users      = map[string]UserData{}
	APIAddress = "127.0.0.1"
	APIPort    = uint16(10085)
	InbountTag = "vryusers"
	Protocol   = "vmess"
	SSMethod   = "aes-256-gcm"
	xrayCtl    *XrayController
)

func Init(grpcAddress string, grpcPort uint16, protocol string) (err error) {
	APIAddress = grpcAddress
	APIPort = grpcPort
	Protocol = protocol
	var (
		cfg = &BaseConfig{
			APIAddress: APIAddress,
			APIPort:    APIPort,
		}
	)
	xrayCtl = new(XrayController)
	err = xrayCtl.Init(cfg)
	if err != nil {
		return
	}
	return
}

func SetSSmethod(method string) (err error) {
	SSMethod = method
	return
}

type UserData struct {
	Uuid  string `json:"uuid"`
	Email string `json:"email"`
}

func ToUser(Uuid string, Email string) UserData {
	return UserData{Uuid: Uuid, Email: Email}
}
func toXuser(user UserData) UserInfo {
	return UserInfo{
		Uuid:       user.Uuid,
		Email:      user.Email,
		Password:   user.Uuid,
		AlertId:    0,
		Level:      0,
		InTag:      InbountTag,
		CipherType: SSMethod,
	}
}

func AddUser(Uuid string, Email string) (err error) {
	user := ToUser(Uuid, Email)
	xuser := toXuser(user)
	if Protocol == "vmess" || Protocol == "vless" {
		err = addVmessUser(xrayCtl.HsClient, &xuser)
	} else if Protocol == "shadowsocks" {
		err = addSSUser(xrayCtl.HsClient, &xuser)
	} else if Protocol == "trojan" {
		err = addTrojanUser(xrayCtl.HsClient, &xuser)
	}

	if err == nil {
		users[Uuid] = user
	}
	return
}
func RemoveUser(Uuid string) (err error) {
	user, has := users[Uuid]
	if !has {
		errors.New("user not found")
	}
	xuser := toXuser(user)
	err = removeUser(xrayCtl.HsClient, &xuser)
	if err == nil {
		delete(users, Uuid)
	}
	return err
}
func Sync(newUsers *[]UserData) map[string]UserData {
	S := make(map[string]UserData)
	for _, user := range *newUsers {
		S[user.Uuid] = user
	}
	for Uuid := range users {
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
	return users
}
func Traffic(reset bool) (res map[string]uint64) {
	res = make(map[string]uint64)
	eres := make(map[string]uint64)
	ptn := ""
	stats, err := queryTraffic(xrayCtl.SsClient, ptn, reset)
	if err != nil {
		return
	}
	Reg := regexp.MustCompile(`user>>>([^>]+)>>>traffic>>>([^>]+)link`)
	for _, stat := range stats {
		match := Reg.FindStringSubmatch(stat.Name)
		if len(match) > 0 {
			eres[match[1]] += uint64(stat.Value)
		}
	}
	for _, user := range users {
		res[user.Uuid] = eres[user.Email]
	}
	return
}
