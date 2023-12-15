package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"ynp/env"

	"github.com/labstack/echo/v4"
)

func GetUserInfo(c echo.Context) error {
	name := c.Param("name")

	rows, err := env.MyDB.Query("SELECT USER_ID, PW, NAME, EMAIL, TEL, LAST_ACCESS_DT, UPDATE_DT, CREATE_DT FROM USERS where name = ?", name)
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

	buff, err := json.Marshal(ui)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(buff))
}

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

	_, err := env.MyDB.Exec("INSERT INTO ghldnjs.USERS(NAME, PW, EMAIL, TEL) VALUES(?, ?, ?, ?);", ui.Name, ui.PW, ui.Email, ui.Tel)
	if err != nil {
		log.Println(err)
	}
	return c.String(http.StatusOK, "SignUp Success.")
}

func IdDuplicateCheck(c echo.Context) error {
	name := c.Param("name")
	id, err := queryId(env.MyDB, name)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return c.JSON(http.StatusMethodNotAllowed, "error!")
		}
		return c.JSON(http.StatusOK, "not exist")
	}
	return c.JSON(http.StatusOK, id)
}

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
		return c.JSON(http.StatusMethodNotAllowed, err)
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
