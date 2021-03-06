package rbuser

import (
	"fmt"
	"os"
)

// CreateHome Create a users home dir and chown it to them
func (user *RbUser) CreateHome() error {
	if err := os.MkdirAll(user.HomeDirectory, os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(fmt.Sprintf("%s/.forward", user.HomeDirectory))
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(fmt.Sprintf("%s\n", user.Altmail)); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	return os.Chown(user.HomeDirectory, user.UIDNumber, user.GidNumber)
}

// CreateWebDir Create a users Web dir and chown it to them
func (user *RbUser) CreateWebDir() error {
	folder := fmt.Sprintf("/webtree/%d/%s", []rune(user.UID)[0], user.UID)
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}
	return os.Chown(folder, user.UIDNumber, user.GidNumber)
}

// LinkPublicHTML Link a users Webdir to their home dir
func (user *RbUser) LinkPublicHTML() error {
	return os.Symlink(fmt.Sprintf("/webtree/%d/%s", []rune(user.UID)[0], user.UID), fmt.Sprintf("%s/public_html", user.HomeDirectory))
}

// MigrateHome migrate a users home dir and chown it to them
func (user *RbUser) MigrateHome(newHome string) error {
	if err := os.Rename(user.HomeDirectory, newHome); err != nil {
		return err
	}
	user.HomeDirectory = newHome
	return user.LinkPublicHTML()
}

// DelWebDir Delete a users web dir
func (user *RbUser) DelWebDir() error {
	return os.RemoveAll(fmt.Sprintf("/webtree/%d/%s", []rune(user.UID)[0], user.UID))
}

// DelHomeDir Delete a users home dir
func (user *RbUser) DelHomeDir() error {
	return os.RemoveAll(user.HomeDirectory)
}

// DelExtraFiles delete leftover files from user
func (user *RbUser) DelExtraFiles() error {
	for _, file := range []string{
		"/local/share/agreement/statedir/%s",
		"/var/mail/%s",
		"/var/spool/cron/crontabs/%s",
	} {
		if err := os.Remove(fmt.Sprintf(file, user.UID)); err != nil {
			return err
		}
	}
	return nil
}
