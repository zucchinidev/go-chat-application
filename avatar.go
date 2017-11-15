package main

import (
	"errors"
	"fmt"
)

// ErrNoAvatar is the error that is returned when the
// Avatar instance is unable to provide an avatar URL
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.

type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct {

}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatarUrl"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct {

}

var UseGravatarAvatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userId"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIdStr), nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct {

}


var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userId"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return fmt.Sprintf("/avatars/%s%s", userIdStr, ".jpg"), nil
		}
	}
	return "", ErrNoAvatarURL
}