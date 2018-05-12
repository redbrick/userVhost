package rbldap

import (
	"errors"

	"github.com/redbrick/rbldap/pkg/rbuser"
	"github.com/urfave/cli"
)

// Update a user in ldap
func Update(ctx *cli.Context) error {
	p := newPrompt()
	// Get User from arg if there prompt if not
	username := ""
	if ctx.NArg() > 0 {
		username = ctx.Args().First()
	} else {
		name, err := p.ReadString("Enter Username")
		if err != nil {
			return err
		}
		username = name
	}
	rb, err := rbuser.NewRbLdap(
		ctx.GlobalString("user"),
		ctx.GlobalString("password"),
		ctx.GlobalString("host"),
		ctx.GlobalInt("port"),
		ctx.GlobalString("smtp"),
	)
	if err != nil {
		return err
	}
	defer rb.Conn.Close()

	user, err := rb.Search(filterAnd(filter("uid", username)))
	if user.UID == "" || err == nil {
		return errors.New("User not found")
	}
	// Prompt for details to change
	return rb.Update(user)
}
