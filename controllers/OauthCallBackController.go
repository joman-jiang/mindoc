package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"strings"
	"time"
)

// OauthCallBackController 认证登陆
type OauthCallBackController struct {
	BaseController
}

// LoginBack 1、验证code
// 2、记录重定向
// 3、重定
func (c *OauthCallBackController) LoginBack() {
	code := c.GetString("code")
	redirectUrl := c.GetSession("REDIRECT_URL").(string)
	appId, _ := web.AppConfig.String("oauth_appId")
	appSecret, _ := web.AppConfig.String("oauth_appSecret")
	oauthRedirectUrl, _ := web.AppConfig.String("oauth_redirectUrl")
	tokenUrl, _ := web.AppConfig.String("oauth_tokenUrl")
	userUrl, _ := web.AppConfig.String("oauth_userUrl")
	// 获取code
	var url = tokenUrl + "?grant_type=authorization_code&code=_CODE&client_id=" + appId + "&client_secret=" + appSecret + "&redirect_uri=" + oauthRedirectUrl
	url_ := strings.Replace(url, "_CODE", code, 1)
	// fmt.Println(url_)
	req := httplib.Get(url_)
	var result Result
	req.ToJSON(&result)
	token := result.Access_token
	//fmt.Println(result.Access_token)
	c.Ctx.SetCookie("mindoc_token", token, time.Now().Add(time.Hour*24*30).Unix())
	// 获取用户
	var url_user = userUrl + "?access_token=_TOKEN"
	url_user_ := strings.Replace(url_user, "_TOKEN", token, 1)
	// fmt.Println(url_user_)
	reqUser := httplib.Get(url_user_)
	var user User
	reqUser.ToJSON(&user)

	fmt.Println(user)

	// 保存用户
	member, err := models.NewMember().Find(user.MemberId)
	logs.Debug(member)
	member.MemberId = user.MemberId
	member.RealName = user.RealName
	member.Avatar = user.Avatar
	member.Phone = user.Phone
	member.Account = user.Account
	member.LastLoginTime = time.Now()
	fmt.Println(err)
	if err != nil {
		fmt.Println("has no")
		member.Email = user.Email
		hash, _ := utils.PasswordHash("joman_123456")
		member.Password = hash
		member.Status = 0
		member.CreateAt = 0
		member.AuthMethod = "ldap"
		member.Role = conf.MemberGeneralRole
		if err1 := member.Add(); err1 != nil {
			fmt.Println(err1)
		}
	} else {
		fmt.Println("has")
		member.Update("phone", "account", "avatar", "real_name", "last_login_time")
	}
	// 设置缓存
	c.SetMember(*member)
	c.Redirect(redirectUrl, 302)

}

type Result struct {
	Access_token string `json:"access_token"`
}
type User struct {
	MemberId int    `json:"id"`
	Account  string `json:"english_name"`
	RealName string `json:"name"`
	//认证方式: local 本地数据库 /ldap LDAP
	Email  string `json:"email"`
	Phone  string `json:"tel"`
	Avatar string `json:"avatar"`
}
