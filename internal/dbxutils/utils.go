package dbxutils

import (
	"blog/internal/types"
	"blog/pkg/dbx"
	"fmt"
	"io"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/sharing"
)

type DbxCommiter struct {
	dbxArgs dbx.DbxArgs
}

func NewDbxCommiter(dbxArgs dbx.DbxArgs) *DbxCommiter {
	return &DbxCommiter{
		dbxArgs: dbxArgs,
	}
}

func (d *DbxCommiter) Upload(fileArgs types.FileArgs) (string, error) {
	defer fileArgs.File.Close()

	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(fileArgs.File)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileBytes)

	fileToUpload, err := os.Open(tempFile.Name())
	if err != nil {
		return "", err
	}
	defer fileToUpload.Close()

	dropboxPath := "/" + fileArgs.Handler.Filename
	res, err := d.dbxArgs.Files.Upload(files.NewUploadArg(dropboxPath), fileToUpload)
	if err != nil {
		return "", err
	}

	_ = res.ContentHash

	sharingClient := sharing.New(d.dbxArgs.Cfg)
	sharedLink, err := sharingClient.CreateSharedLinkWithSettings(
		&sharing.CreateSharedLinkWithSettingsArg{Path: dropboxPath},
	)
	if err != nil {
		return "", err
	}
	
	var linkUrl string

	switch v := sharedLink.(type) {
	case *sharing.FileLinkMetadata:
		linkUrl = v.Url
	case *sharing.FolderLinkMetadata:
		linkUrl = v.Url
	default:
		return "", fmt.Errorf("unknown type: %s", v)
	}

	return linkUrl, nil
}
