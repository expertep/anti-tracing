package controller

import (
	"context"
	"errors"
	"jaeger-ot/model"
	"time"

	"github.com/gin-gonic/gin"
)

func ValidateUser(user model.User) error {
	time.Sleep(500 * time.Millisecond)

	if user.Username == "" {
		return errors.New("Username is required")
	}
	return nil
}

func ValidateUserV2(ot otelCus.IOtelPattern, user model.User) error {
	ot.ChildSpan("ValidateUser")
	defer ot.SpanEnd()
	time.Sleep(500 * time.Millisecond)

	if user.Username == "" {
		return errors.New("Username is required")
	}
	return nil
}

func CreateUser(user model.User) error {
	time.Sleep(500 * time.Millisecond)

	if user.Username == "a" {
		return errors.New("Username already exists")
	}
	return nil
}

/* ewewe */
func SignupV2() func(c *gin.Context) {
	return func(c *gin.Context) {
		var ctx context.Context
		ot := otelCus.CreateOtel(ctx, "signup", "signup")
		ot.ParentSpan("signup")

		// ot.ChildSpan("ShouldBind")
		var user model.User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			ot.AddEventParent("Error1", ot.Att("Error", err.Error()))
			ot.HandlerFail("Error12324", err.Error())
			ot.OtelEnd()
			return
		}
		time.Sleep(500 * time.Millisecond)
		// ot.SpanEnd()

		ot.AddEventByStructParent("obj", user)
		ot.SetAttributesParent(ot.St("username1", user.Username))

		ot.ChildSpan("ValidateUser")
		if err := ValidateUserV2(ot, user); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			ot.AddEventParent("Error1", ot.Att("Error", err.Error()))
			ot.HandlerFail("Error", err.Error())
			ot.OtelEnd()
			return
		}
		ot.SpanEnd()

		ot.ChildSpan("CreateUser")
		if err := CreateUser(user); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			ot.AddEventParent("Error1", ot.Att("Error", err.Error()))
			ot.HandlerFail("Error", err.Error())
			ot.OtelEnd()
			return
		}
		ot.SpanEnd()

		/* _, spanChild := tracer.Start(ctxParent, "ValidateUser")
		defer spanChild.End() */

		ot.SpanEndParent()
		ot.AddEventParent("Success", ot.Att("Success", "Check cashback current"))
		ot.HandlerSuccess("Check cashback current")
		ot.HandlerSuccessParent("withdraw.cashback Success")
		ot.SetAttributesByStructParent(user)
		ot.OtelEnd()

		c.JSON(200, gin.H{
			"code":    200,
			"message": "Success",
		})

		return
	}
}

// otel
func SignupV3() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}
		time.Sleep(500 * time.Millisecond)

		//o_p user:$user
		//o_p username:$user.Username

		if err := ValidateUser(user); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

		if err := CreateUser(user); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

		/* o_p "Create user success"
		$username */

		c.JSON(200, gin.H{
			"code":    200,
			"message": "Success",
		})

		return
	}
}
