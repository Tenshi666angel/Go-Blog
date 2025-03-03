package types

import "mime/multipart"

type FileArgs struct {
	File    multipart.File
	Handler multipart.FileHeader
}

	
