package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

type TryAvatars []Avatar

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar}

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

//ErrNoAvatarURL is the error that is returned when the
//Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

//Avatar represents types capable of representing
//user profile pictures.
type Avatar interface {
	//GetAvatarURL get the avatar URL for the specified client,
	//or returns an error if something goes wront.
	//ErrNoAvatarURL is returned if the object is unable to get
	//a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

//Gravatar is an implementation of Avatar interface
//which relies on the OAuth authentication info site
type AuthAvatar struct{}

//UseAuthAvatar is the reference structure to use the
//AuthAvatar implementation: users can avoid worrying
//about the object instantiation
var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

//Gravatar is an implementation of Avatar interface
//which relies on the https://en.gravatar.com/ site
type GravatarAvatar struct{}

//UseGravatar is the reference structure to use the
//Gravatar implementation: users can avoid worrying
//about the object instantiation
var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return fmt.Sprintf("//www.gravatar.com/avatar/%s", u.UniqueID()), nil
}

//FileSystemAvatar is an implementation of Avatar interface
//which relies on user uploaded data
type FileSystemAvatar struct{}

//UseFileSystemAvatar is the reference structure to use the
//FileSystemAvatar implementation: users can avoid worrying
//about the object instantiation
var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	userID := u.UniqueID()

	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}

	//check for user avatar in filesystem
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := path.Match(userID+"*", file.Name()); match {
			return fmt.Sprintf("/avatars/%s", file.Name()), nil
		}
	}

	return "", ErrNoAvatarURL
}
