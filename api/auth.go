package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"ynp/env"

	"github.com/labstack/echo/v4"
)

// @Summary Get User Info
// @Description Get users information
// @Accept json
// @Produce json
// @Param name path string true "Desired Name"
// @Success 200 {object} Res
// @Router /users/get/{name} [get]
func GetUserInfo(c echo.Context) error {
	name := c.Param("name")

	rows, err := env.MyDB.Query("SELECT USER_ID, PW, NAME, EMAIL, TEL, LAST_ACCESS_DT, UPDATE_DT, CREATE_DT FROM USERS where name = ?;", name)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ui Users
	for rows.Next() {
		err := rows.Scan(&ui.UserId, &ui.PW, &ui.Name, &ui.Email, &ui.Tel, &ui.LastAccessDt, &ui.UpdateDt, &ui.CreateDt)
		if err != nil {
			return err
		}
	}

	r := &Res{
		Status:   http.StatusOK,
		Response: ui,
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

// @Summary New User Info
// @Description Create users information - SignUp
// @Accept json
// @Produce json
// @Param name formData string true "User's name"
// @Param password formData string true "User's password"
// @Param email formData string true "User's email"
// @Param tel formData string true "User's tel number"
// @Success 200 {object} Res
// @Router /users/signup [post]
func NewUserInfo(c echo.Context) error {

	ui := &Users{
		Name:  c.FormValue("name"),
		PW:    c.FormValue("password"),
		Email: c.FormValue("email"),
		Tel:   c.FormValue("tel"),
	}
	if ui.Name == "" || ui.PW == "" || ui.Email == "" || ui.Tel == "" {
		return c.JSON(http.StatusMethodNotAllowed, "invalid name pw email tel")
	}

	_, err := env.MyDB.Exec("INSERT INTO ghldnjs.USERS(NAME, PW, EMAIL, TEL, LAST_ACCESS_DT) VALUES(?, ?, ?, ?, ?);", ui.Name, ui.PW, ui.Email, ui.Tel, time.Now())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusMethodNotAllowed, "insert fail")
	}

	r := &Res{
		Status:   http.StatusOK,
		Response: "SignUp Seccess.",
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

// @Summary User's name duplicate check
// @Description Check a name for SignUp
// @Accept json
// @Produce json
// @Param name path string true "User's name"
// @Success 200 {object} Res
// @Router /users/exist/{name} [get]
func IdDuplicateCheck(c echo.Context) error {
	name := c.Param("name")
	id, err := queryId(env.MyDB, name)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return c.JSON(http.StatusMethodNotAllowed, "error!")
		}
		r := &Res{
			Status:   http.StatusOK,
			Response: "not exist",
		}
		return c.JSONPretty(http.StatusOK, r, " ")
	}
	r := &Res{
		Status:   http.StatusOK,
		Response: id,
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}

// @Summary Sign in
// @Description Sign in function
// @Accept json
// @Produce json
// @Param name formData string true "User's name"
// @Param password formData string true "User's password"
// @Success 200 {object} Res
// @Router /users/signin [post]
func SignIn(c echo.Context) error {
	params := make(map[string]string)
	var id int

	name := c.FormValue("name")
	password := c.FormValue("password")

	if name == "" || password == "" {
		return c.JSON(http.StatusUnauthorized, "invalid name pw")
	}

	id, err := queryPw(env.MyDB, name, password)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, "db error")
	}
	_, err = env.MyDB.Exec("update ghldnjs.USERS SET LAST_ACCESS_DT = ?;", time.Now())
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, "db error")
	}

	accessToken, err := generateToken(c, id, name)
	if err != nil {
		params["token"] = fmt.Sprint(err)
		return c.JSON(http.StatusMethodNotAllowed, params["token"])
	}

	c.Response().Header().Set("Cache-Control", "no-store no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0")
	c.Response().Header().Add("Last-Modified", time.Now().String())
	c.Response().Header().Add("pragma", "no-cache")
	// c.Response().Header().Add("Expires", "-1")
	cookie := new(http.Cookie)
	cookie.Name = "access-token"
	cookie.Value = accessToken
	cookie.Expires = time.Now().Add(ExpirationTime)
	c.SetCookie(cookie)
	params["token"] = accessToken

	r := &Res{
		Status:   http.StatusOK,
		Response: params["token"],
	}
	return c.JSONPretty(http.StatusOK, r, " ")
}
