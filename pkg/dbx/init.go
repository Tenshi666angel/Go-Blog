package dbx

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type DbxArgs struct {
	Cfg   dropbox.Config
	Files files.Client
}

func InitDbx(token string) *DbxArgs {
	config := dropbox.Config{
		Token: token,
	}
	dbx := files.New(config)

	return &DbxArgs{
		Cfg:   config,
		Files: dbx,
	}
}
