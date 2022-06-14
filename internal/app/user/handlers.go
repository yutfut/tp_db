package user

import (
	"encoding/json"
	"example.com/greetings/internal/app/models"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Handlers struct {
	UserRepo UserRep
}

func (h *Handlers) CreateUser(ctx *fasthttp.RequestCtx) {
	var user models.User

	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		body, _ := json.Marshal(err.Error())
		ctx.SetBody(body)
		return
	}
	user.Nickname = fmt.Sprintf("%s", ctx.UserValue("nickname"))

	user1, err1 := h.UserRepo.GetUserByNickname(user.Nickname)
	user2, err2 := h.UserRepo.GetUserByEmail(user.Email)

	if err1 == nil || err2 == nil {
		var users []models.User
		if err1 == nil {
			users = append(users, user1)
		}
		if err2 == nil && user1.About != user2.About {
			users = append(users, user2)
		}
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(users)
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	_, err = h.UserRepo.CreateUser(user)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(user)
	ctx.SetStatusCode(http.StatusCreated)
	ctx.SetBody(body)
}

func (h *Handlers) GetProfileByNickname(ctx *fasthttp.RequestCtx) {
	nickname := fmt.Sprintf("%s", ctx.UserValue("nickname"))
	user, err := h.UserRepo.GetUserByNickname(nickname)

	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: "Can't find user by nickname:"})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(user)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}

func (h *Handlers) UpdateProfile(ctx *fasthttp.RequestCtx) {
	nickname := fmt.Sprintf("%s", ctx.UserValue("nickname"))
	newUserData, err := h.UserRepo.GetUserByNickname(nickname)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message: "Can't find user by nickname:"})
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	err = json.Unmarshal(ctx.PostBody(), &newUserData)
	if err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(http.StatusBadRequest)
		body, _ := json.Marshal(err.Error())
		ctx.SetBody(body)
		return
	}

	checkUser, err := h.UserRepo.GetUserByEmail(newUserData.Email)
	if !checkUser.IsEmpty() && checkUser.Nickname != newUserData.Nickname {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(models.MessageError{Message:"This email is already registered by user:"})
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetBody(body)
		return
	}

	user, err := h.UserRepo.UpdateProfile(newUserData)
	if err != nil {
		ctx.SetContentType("application/json")
		body, _ := json.Marshal(err.Error())
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBody(body)
		return
	}

	if newUserData.About == "" {
		user.About = newUserData.About
	}

	ctx.SetContentType("application/json")
	body, _ := json.Marshal(user)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(body)
}
