package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"net/url"
)

// OauthController 认证登陆
type OauthController struct {
	AccountController
}

// Login 1、系统不存在session时调用这个页面
// 2、调整到认证系统
func (c *OauthController) Login() {

	appId, _ := web.AppConfig.String("oauth_appId")
	redirectUrl, _ := web.AppConfig.String("oauth_redirectUrl")
	authorizeUrl, _ := web.AppConfig.String("oauth_authorizeUrl")
	oauthLogoutUrl, _ := web.AppConfig.String("oauth_logoutUrl")
	reUrl := c.referer()
	c.SetSession("REDIRECT_URL", reUrl)
	var authUrl = authorizeUrl + "?response_type=code&grant_type=authorization_code&scope=code&client_id=" + appId + "&redirect_uri=" + redirectUrl
	escape := url.QueryEscape(authUrl)
	c.Redirect(oauthLogoutUrl+escape, 302)
}
