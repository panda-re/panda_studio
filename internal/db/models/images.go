package models

import "github.com/panda-re/panda_studio/internal/db"

type Architecture string

const (
	X86_64 Architecture = "x86_64"
	I386 Architecture = "i386"
	ARM Architecture = "arm"
	AARCH64 Architecture = "aarch64"
)

// Entities
type Image struct {
	ID db.ObjectID `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Files []*ImageFile `bson:"files" json:"files"`
	Config *ImageConfiguration `bson:"config" json:"config"`
}

type ImageFile struct {
	ID db.ObjectID  `bson:"_id" json:"id"`
	ImageID db.ObjectID `bson:"-" json:"imageId"`
	FileName string `bson:"file_name" json:"file_name"`
	FileType string `bson:"file_type" json:"file_type"`
	IsUploaded bool `bson:"is_uploaded" json:"is_uploaded"`
	Size int64		`bson:"size" json:"size"`
	Sha256 string   `bson:"sha256" json:"sha256,omitempty"`
}

type ImageConfiguration struct {
	Architecture Architecture `bson:"arch" json:"arch"`
	DiskImage db.ObjectID `bson:"disk_image" json:"disk_image"`
	// Kernel string
	// Initrd string
	// DeviceTree string
}

// Data transfer objects
type ImageFileCreateRequest struct {
	ImageID db.ObjectID `json:"image_id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
}

type ImageFileUploadRequest struct {
	ImageId db.ObjectID
	FileId db.ObjectID
}
