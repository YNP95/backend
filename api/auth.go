package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetUserInfo(c echo.Context) error {
	name := c.Param("name")

	db := NewDb()
	defer CloseDb(db)

	rows, err := db.Query("SELECT id, name, password, created, updated  FROM user_table where name = ?", name)
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
		Name: c.FormValue("name"),
		PW:   c.FormValue("password"),
	}
	if ui.Name == "" || ui.PW == "" {
		return c.JSON(http.StatusMethodNotAllowed, "invalid name pw")
	}

	db := NewDb()
	defer CloseDb(db)
	_, err := db.Exec("INSERT INTO USERS(NAME, PW) VALUES(?, ?)", ui.Name, ui.PW)
	if err != nil {
		log.Println(err)
	}
	return c.String(http.StatusOK, "SignUp Success.")
}

func SignIn(c echo.Context) error {
	db := NewDb()
	defer CloseDb(db)
	params := make(map[string]string)
	var id int

	name := c.FormValue("name")
	password := c.FormValue("password")

	if name == "" || password == "" {
		return c.JSON(http.StatusMethodNotAllowed, "invalid name pw")
	}

	id, err := queryPw(db, name, password)
	if err != nil {
		return c.JSON(http.StatusMethodNotAllowed, err)
	}

	// if password != dbPassword {
	// 	params["pwd"] = "no match"
	// 	_ = c.Bind(&params)

	// 	return c.JSON(http.StatusMethodNotAllowed, params["pwd"])
	// }

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

	return c.JSON(http.StatusOK, params["token"])
}
