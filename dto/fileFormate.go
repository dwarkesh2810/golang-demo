package dto

type FileType struct {
	FileType    string `json:"file_type" form:"file_type"`
	Limit       int    `json:"limit" form:"limit"`
	SratingFrom int    `json:"starting_from" form:"starting_from"`
}
