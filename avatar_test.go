package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gmntest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	//when declaring and no assign variable is nil
	//but Go has default initialization in go it is
	//acceptable to call a method on a nil object
	//provided that the method doesn't try to access
	//a field
	var authAvatar AuthAvatar
	testUser := &gmntest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar should return ErrNoAvatarURL when no value present")
	}
	//set value
	testURL := "http://url-to-gravatar/"
	testUser = &gmntest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testURL, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}

	if url != testURL {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatar.GetAvatarURL(user)
	if err != nil {
		t.Error("Gravatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	var fsGravatar FileSystemAvatar
	//creating fake avatar file for user "abc"
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)

	user := &chatUser{uniqueID: "abc"}
	url, err := fsGravatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
