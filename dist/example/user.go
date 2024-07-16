package controller

import (
	"context"
	"errors"
	"jaeger-ot/model"
	"time"
	"github.com/gin-gonic/gin"
	fmt "fmt"
)

func ValidateUser(user model.User, ot otelCus.IOtelPattern) error {
	fmt.Println("ValidateUser")
	defer fmt.Println("close ValidateUser")
	time.Sleep(500 * time.Millisecond)
	if user.Username == "" {
		return errors.New("Username is required")
	}
	return nil
}
func ValidateUserV2(ot otelCus.IOtelPattern, ot otelCus.IOtelPattern) error {
	fmt.Println("ValidateUserV2")
	defer fmt.Println("close ValidateUserV2")
	ot.ChildSpan("ValidateUser")
	defer ot.SpanEnd()
	time.Sleep(500 * time.Millisecond)
	if user.Username == "" {
		return errors.New("Username is required")
	}
	return nil
}
func CreateUser(user model.User, ot otelCus.IOtelPattern) error {
	fmt.Println("CreateUser")
	defer fmt.Println("close CreateUser")
	time.Sleep(500 * time.Millisecond)
	if user.Username == "a" {
		return errors.New("Username already exists")
	}
	return nil
}
func SignupV2() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()
		ot := otelCus.CreateOtel(ctx, "SignupV2", "SignupV2")
		fmt.Println("SignupV2")
		defer fmt.Println("close SignupV2")
		var ctx context.Context
		ot := otelCus.CreateOtel(ctx, "signup", "signup")
		ot.ParentSpan("signup")
		var user model.User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error()})
			ot.AddEventParent("Error1", ot.Att("Error", err.Error()))
			ot.HandlerFail("Error12324", err.Error())
			ot.OtelEnd()
			return
		}
		time.Sleep(500 * time.Millisecond)
		ot.AddEventByStructParent("obj", user)
		ot.SetAttributesParent(ot.St("username1", user.Username))
		ot.ChildSpan("ValidateUser")
		if err := ValidateUserV2(ot, user); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error()})
			ot.AddEventParent("Error1", ot.Att("Error", err.Error()))
			ot.HandlerFail("Error", err.Error())
			ot.OtelEnd()
			return
		}
		ot.SpanEnd()
		ot.ChildSpan("CreateUser")
		if err := CreateUser(user); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error()})
			ot.AddEventParent("Error1", ot.Att("Error", err.Error()))
			ot.HandlerFail("Error", err.Error())
			ot.OtelEnd()
			return
		}
		ot.SpanEnd()
		ot.SpanEndParent()
		ot.AddEventParent("Success", ot.Att("Success", "Check cashback current"))
		ot.HandlerSuccess("Check cashback current")
		ot.HandlerSuccessParent("withdraw.cashback Success")
		ot.SetAttributesByStructParent(user)
		ot.OtelEnd()
		c.JSON(200, gin.H{"code": 200, "message": "Success"})
		return
	}
}
func SignupV3() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()
		ot := otelCus.CreateOtel(ctx, "SignupV3", "SignupV3")
		fmt.Println("SignupV3")
		defer fmt.Println("close SignupV3")
		var user model.User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}
		time.Sleep(500 * time.Millisecond)
		if err := ValidateUser(user); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}
		if err := CreateUser(user); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 200, "message": "Success"})
		return
	}
}
