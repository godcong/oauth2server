package oauth2server

import (
	"net/http"

	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/godcong/oauth2server/model"
)

//sub	string	Subject - Identifier for the End-User at the Issuer.
//name	string	End-User's full name in displayable form including all name parts, possibly including titles and suffixes, ordered according to the End-User's locale and preferences.
//given_name	string	Given name(s) or first name(s) of the End-User. Note that in some cultures, people can have multiple given names; all can be present, with the names being separated by space characters.
//family_name	string	Surname(s) or last name(s) of the End-User. Note that in some cultures, people can have multiple family names or no family name; all can be present, with the names being separated by space characters.
//middle_name	string	Middle name(s) of the End-User. Note that in some cultures, people can have multiple middle names; all can be present, with the names being separated by space characters. Also note that in some cultures, middle names are not used.
//nickname	string	Casual name of the End-User that may or may not be the same as the given_name. For instance, a nickname value of Mike might be returned alongside a given_name value of Michael.
//preferred_username	string	Shorthand name by which the End-User wishes to be referred to at the RP, such as janedoe or j.doe. This value MAY be any valid JSON string including special characters such as @, /, or whitespace. The RP MUST NOT rely upon this value being unique, as discussed in Section 2.5.3.
//profile	string	URL of the End-User's profile page. The contents of this Web page SHOULD be about the End-User.
//picture	string	URL of the End-User's profile picture. This URL MUST refer to an image file (for example, a PNG, JPEG, or GIF image file), rather than to a Web page containing an image. Note that this URL SHOULD specifically reference a profile photo of the End-User suitable for displaying when describing the End-User, rather than an arbitrary photo taken by the End-User.
//website	string	URL of the End-User's Web page or blog. This Web page SHOULD contain information published by the End-User or an organization that the End-User is affiliated with.
//email	string	End-User's preferred e-mail address. Its value MUST conform to the RFC 5322 [RFC5322] addr-spec syntax. The RP MUST NOT rely upon this value being unique, as discussed in Section 2.5.3.
//email_verified	boolean	True if the End-User's e-mail address has been verified; otherwise false. When this Claim Value is true, this means that the OP took affirmative steps to ensure that this e-mail address was controlled by the End-User at the time the verification was performed. The means by which an e-mail address is verified is context specific, and dependent upon the trust framework or contractual agreements within which the parties are operating.
//gender	string	End-User's gender. Values defined by this document are female and male. Other values MAY be used when neither of the defined values are applicable.
//birthdate	string	End-User's birthday, represented as an ISO 8601:2004 [ISO8601‑2004] YYYY-MM-DD format. The year MAY be 0000, indicating that it is omitted. To represent only the year, YYYY format is allowed. Note that depending on the underlying platform's date related function, providing just year can result in varying month and day, so the implementers need to take this factor into account to correctly process the dates.
//zoneinfo	string	String from zoneinfo [zoneinfo] time zone database representing the End-User's time zone. For example, Europe/Paris or America/Los_Angeles.
//locale	string	End-User's locale, represented as a BCP47 [RFC5646] language tag. This is typically an ISO 639-1 Alpha-2 [ISO639‑1] language code in lowercase and an ISO 3166-1 Alpha-2 [ISO3166‑1] country code in uppercase, separated by a dash. For example, en-US or fr-CA. As a compatibility note, some implementations have used an underscore as the separator rather than a dash, for example, en_US; Relying Parties MAY choose to accept this locale syntax as well.
//phone_number	string	End-User's preferred telephone number. E.164 [E.164] is RECOMMENDED as the format of this Claim, for example, +1 (425) 555-1212 or +56 (2) 687 2400. If the phone number contains an extension, it is RECOMMENDED that the extension be represented using the RFC 3966 [RFC3966] extension syntax, for example, +1 (604) 555-1234;ext=5678.
//phone_number_verified	boolean	True if the End-User's phone number has been verified; otherwise false. When this Claim Value is true, this means that the OP took affirmative steps to ensure that this phone number was controlled by the End-User at the time the verification was performed. The means by which a phone number is verified is context specific, and dependent upon the trust framework or contractual agreements within which the parties are operating. When true, the phone_number Claim MUST be in E.164 format and any extensions MUST be represented in RFC 3966 format.
//address	JSON object	End-User's preferred postal address. The value of the address member is a JSON [RFC4627] structure containing some or all of the members defined in Section 2.5.1.
//updated_at	number	Time the End-User's information was last updated. Its value is a JSON number representing the number of seconds from 1970-01-01T00:00:00Z as measured in UTC until the date/time.
func userinfoGet(c *gin.Context) {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	j, e := GetRedis().Do("GET", c.Request.Form.Get("access_token"))
	fmt.Println("j:", string(j.([]byte)), e)

	if j == nil || e != nil {
		ResponseError(c, E_INVALID_REQUEST)
		return
	}

	user := model.User{}
	author := model.Authorize{}
	cli := model.Client{}

	model.Gorm().First(&author, "sub_id = ?", string(j.([]byte)))
	model.Gorm().Model(&author).Related(&user)
	model.Gorm().Model(&author).Related(&cli)

	log.Println("userinfoGet user: " + user.ID.String())
	log.Println("userinfoGet client: " + cli.ID.String())

	//<<<<<<< HEAD
	//	c.JSON(http.StatusAccepted, gin.H{
	//		"sub": author.SubID,
	//		"nickname": user.Nickname,
	//		"name": user.Username,
	//		"phone_number": user.Mobile,
	//		"email": user.Mail,
	//		"picture": "",
	//	})

	sc, err := GetRedis().Do("GET", cli.ID.String())

	if err != nil {
		ResponseError(c, E_INVALID_REQUEST)
		return
	}
	if sc == nil {
		c.JSON(http.StatusAccepted, gin.H{
			"sub":          author.SubID,
			"nickname":     user.Nickname,
			"name":         user.Username,
			"phone_number": user.Mobile,
			"email":        user.Mail,
			"picture":      "",
		})
		return
	}

	sco := cli.GetScope()
	scope := CheckScope(string(sc.([]byte)), sco)
	info := GetInfo(scope, &user, &author)
	fmt.Println("info:", info)

	c.JSON(http.StatusAccepted, info)

}

func GetInfo(scope []string, user *model.User, author *model.Authorize) map[string]string {

	info := make(map[string]string)
	info["sub"] = author.SubID

	for _, v := range scope {
		switch v {
		case "profile":
			info["nickname"] = user.Nickname
			info["name"] = user.Username
			info["picture"] = ""
		case "email":
			info["email"] = user.Mail
		case "address":
			info["address"] = ""
		case "phone":
			info["phone_number"] = user.Mobile
		}

	}
	return info
}
func CheckScope(s, s1 string) []string {

	scope := []string{}
	sc := strings.Split(s, " ")
	for _, v := range sc {
		if v != "openid" && strings.Contains(s1, v) {
			scope = append(scope, v)
		}

	}
	return scope
}

//func infoMake(scope string) {

//	//=======
//	//user.FirstByID(string(j.([]byte)))

//	aginh := gin.H{
//		"sub":          author.SubID,
//		"nickname":     user.Nickname,
//		"name":         user.Username,
//		"phone_number": user.Mobile,
//		"email":        user.Mail,
//		"picture":      ""}
//	log.Println(aginh)
//	c.JSON(http.StatusAccepted, aginh)
//	//>>>>>>> 61ccd1cfb55d69afe946a99abd5151b5ce3cad0b
//}
