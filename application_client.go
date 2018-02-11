package oauth2server

//"github.com/RangelReale/osin"

const HOST = "http://localhost:7788"

//func DownloadAccessToken(url string, auth *osin.BasicAuth, output map[string]interface{}) error {
//	// download access token
//	preq, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return err
//	}
//	if auth != nil {
//		preq.SetBasicAuth(auth.Username, auth.Password)
//	}
//
//	pclient := &http.Client{}
//	presp, err := pclient.Do(preq)
//	if err != nil {
//		return err
//	}
//
//	if presp.StatusCode != 200 {
//		return errors.New("Invalid status code")
//	}
//
//	jdec := json.NewDecoder(presp.Body)
//	err = jdec.Decode(&output)
//	return err
//}
//
//// Application destination - CODE
//func appCode(c *gin.Context) {
//	r := c.Request
//	w := c.Writer
//	r.ParseForm()
//
//	code := r.Form.Get("code")
//
//	w.Write([]byte("<html><body>"))
//	w.Write([]byte("APP AUTH - CODE<br/>"))
//	defer w.Write([]byte("</body></html>"))
//
//	if code == "" {
//		w.Write([]byte("Nothing to do"))
//		return
//	}
//
//	jr := make(map[string]interface{})
//
//	// build access code url
//	aurl := fmt.Sprintf("/token?grant_type=authorization_code&client_id=1234&client_secret=aabbccdd&state=xyz&redirect_uri=%s&code=%s",
//		url.QueryEscape(HOST+"/app/code"), url.QueryEscape(code))
//
//	// if parse, download and parse json
//	if r.Form.Get("doparse") == "1" {
//		err := DownloadAccessToken(fmt.Sprintf(HOST+"%s", aurl),
//			&osin.BasicAuth{"1234", "aabbccdd"}, jr)
//		if err != nil {
//			w.Write([]byte(err.Error()))
//			w.Write([]byte("<br/>"))
//		}
//	}
//
//	// show json error
//	if erd, ok := jr["error"]; ok {
//		w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
//	}
//
//	// show json access token
//	if at, ok := jr["access_token"]; ok {
//		w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
//	}
//
//	w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))
//
//	// output links
//	w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Goto Token URL</a><br/>", aurl)))
//
//	cururl := *r.URL
//	curq := cururl.Query()
//	curq.Add("doparse", "1")
//	cururl.RawQuery = curq.Encode()
//	w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Download Token</a><br/>", cururl.String())))
//
//	if rt, ok := jr["refresh_token"]; ok {
//		rurl := fmt.Sprintf("/app/refresh?code=%s", rt)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
//	}
//
//	if at, ok := jr["access_token"]; ok {
//		rurl := fmt.Sprintf("/app/info?code=%s", at)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
//	}
//}
//
//func appToken(c *gin.Context) {
//	r := c.Request
//	w := c.Writer
//	r.ParseForm()
//
//	w.Write([]byte("<html><body>"))
//	w.Write([]byte("APP AUTH - TOKEN<br/>"))
//
//	w.Write([]byte("Response data in fragment - not acessible via server - Nothing to do"))
//
//	w.Write([]byte("</body></html>"))
//}
//
//func appPassword(c *gin.Context) {
//	r, w := c.Request, c.Writer
//	r.ParseForm()
//
//	w.Write([]byte("<html><body>"))
//	w.Write([]byte("APP AUTH - PASSWORD<br/>"))
//
//	jr := make(map[string]interface{})
//
//	// build access code url
//	aurl := fmt.Sprintf("/token?grant_type=password&scope=everything&username=%s&password=%s",
//		"test", "test")
//
//	// download token
//	err := DownloadAccessToken(fmt.Sprintf(HOST+"%s", aurl),
//		&osin.BasicAuth{Username: "1234", Password: "aabbccdd"}, jr)
//	if err != nil {
//		w.Write([]byte(err.Error()))
//		w.Write([]byte("<br/>"))
//	}
//
//	// show json error
//	if erd, ok := jr["error"]; ok {
//		w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
//	}
//
//	// show json access token
//	if at, ok := jr["access_token"]; ok {
//		w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
//	}
//
//	w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))
//
//	if rt, ok := jr["refresh_token"]; ok {
//		rurl := fmt.Sprintf("/app/refresh?code=%s", rt)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
//	}
//
//	if at, ok := jr["access_token"]; ok {
//		rurl := fmt.Sprintf("/app/info?code=%s", at)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
//	}
//
//	w.Write([]byte("</body></html>"))
//}
//
//func appClientCredentials(c *gin.Context) {
//	r, w := c.Request, c.Writer
//	r.ParseForm()
//
//	w.Write([]byte("<html><body>"))
//	w.Write([]byte("APP AUTH - CLIENT CREDENTIALS<br/>"))
//
//	jr := make(map[string]interface{})
//
//	// build access code url
//	aurl := fmt.Sprintf("/token?grant_type=client_credentials")
//
//	// download token
//	err := DownloadAccessToken(fmt.Sprintf(HOST+"%s", aurl),
//		&osin.BasicAuth{Username: "1234", Password: "aabbccdd"}, jr)
//	if err != nil {
//		w.Write([]byte(err.Error()))
//		w.Write([]byte("<br/>"))
//	}
//
//	// show json error
//	if erd, ok := jr["error"]; ok {
//		w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
//	}
//
//	// show json access token
//	if at, ok := jr["access_token"]; ok {
//		w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
//	}
//
//	w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))
//
//	if rt, ok := jr["refresh_token"]; ok {
//		rurl := fmt.Sprintf("/app/refresh?code=%s", rt)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
//	}
//
//	if at, ok := jr["access_token"]; ok {
//		rurl := fmt.Sprintf("/app/info?code=%s", at)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
//	}
//
//	w.Write([]byte("</body></html>"))
//}
//
//func appAssertion(c *gin.Context) {
//	r, w := c.Request, c.Writer
//	r.ParseForm()
//
//	w.Write([]byte("<html><body>"))
//	w.Write([]byte("APP AUTH - ASSERTION<br/>"))
//
//	jr := make(map[string]interface{})
//
//	// build access code url
//	aurl := fmt.Sprintf("/token?grant_type=assertion&assertion_type=urn:osin.example.complete&assertion=osin.data")
//
//	// download token
//	err := DownloadAccessToken(fmt.Sprintf(HOST+"%s", aurl),
//		&osin.BasicAuth{Username: "1234", Password: "aabbccdd"}, jr)
//	if err != nil {
//		w.Write([]byte(err.Error()))
//		w.Write([]byte("<br/>"))
//	}
//
//	// show json error
//	if erd, ok := jr["error"]; ok {
//		w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
//	}
//
//	// show json access token
//	if at, ok := jr["access_token"]; ok {
//		w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
//	}
//
//	w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))
//
//	if rt, ok := jr["refresh_token"]; ok {
//		rurl := fmt.Sprintf("/app/refresh?code=%s", rt)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
//	}
//
//	if at, ok := jr["access_token"]; ok {
//		rurl := fmt.Sprintf("/app/info?code=%s", at)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
//	}
//
//	w.Write([]byte("</body></html>"))
//}
//
//func appRefresh(c *gin.Context) {
//	r, w := c.Request, c.Writer
//	r.ParseForm()
//
//	w.Write([]byte("<html><body>"))
//	w.Write([]byte("APP AUTH - REFRESH<br/>"))
//	defer w.Write([]byte("</body></html>"))
//
//	code := r.Form.Get("code")
//
//	if code == "" {
//		w.Write([]byte("Nothing to do"))
//		return
//	}
//
//	jr := make(map[string]interface{})
//
//	// build access code url
//	aurl := fmt.Sprintf("/token?grant_type=refresh_token&refresh_token=%s", url.QueryEscape(code))
//
//	// download token
//	err := DownloadAccessToken(fmt.Sprintf(HOST+"%s", aurl),
//		&osin.BasicAuth{Username: "1234", Password: "aabbccdd"}, jr)
//	if err != nil {
//		w.Write([]byte(err.Error()))
//		w.Write([]byte("<br/>"))
//	}
//
//	// show json error
//	if erd, ok := jr["error"]; ok {
//		w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
//	}
//
//	// show json access token
//	if at, ok := jr["access_token"]; ok {
//		w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
//	}
//
//	w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))
//
//	if rt, ok := jr["refresh_token"]; ok {
//		rurl := fmt.Sprintf("/app/refresh?code=%s", rt)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
//	}
//
//	if at, ok := jr["access_token"]; ok {
//		rurl := fmt.Sprintf("/app/info?code=%s", at)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
//	}
//}
//
//func appInfo(c *gin.Context) {
//	r, w := c.Request, c.Writer
//	r.ParseForm()
//
//	w.Write([]byte("<html><body>"))
//	w.Write([]byte("APP AUTH - INFO<br/>"))
//	defer w.Write([]byte("</body></html>"))
//
//	code := r.Form.Get("code")
//
//	if code == "" {
//		w.Write([]byte("Nothing to do"))
//		return
//	}
//
//	jr := make(map[string]interface{})
//
//	// build access code url
//	aurl := fmt.Sprintf("/info?code=%s", url.QueryEscape(code))
//
//	// download token
//	err := DownloadAccessToken(fmt.Sprintf(HOST+"%s", aurl),
//		&osin.BasicAuth{Username: "1234", Password: "aabbccdd"}, jr)
//	if err != nil {
//		w.Write([]byte(err.Error()))
//		w.Write([]byte("<br/>"))
//	}
//
//	// show json error
//	if erd, ok := jr["error"]; ok {
//		w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
//	}
//
//	// show json access token
//	if at, ok := jr["access_token"]; ok {
//		w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
//	}
//
//	w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))
//
//	if rt, ok := jr["refresh_token"]; ok {
//		rurl := fmt.Sprintf("/app/refresh?code=%s", rt)
//		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
//	}
//}
