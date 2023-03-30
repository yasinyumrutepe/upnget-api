package models

import (
	"auction/globals"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	// slicehelper "github.com/roofdomain/probnet/sliceHelper"
)

type File struct {
	ID        uint64
	TableType string
	TableID   uint
	URL       string
	DataType  uint8
}

type Files struct {
	File    []File                  `json:"file"`
	ErrFile []*multipart.FileHeader `json:"err_file"`
}

func (fileses *Files) SaveFile(tableName string, fileTypes []uint8, c *fiber.Ctx, files []*multipart.FileHeader) error {
	var savedFiles []*multipart.FileHeader
	var expectedFileTypes []uint8
	var acceptedTypes = map[uint8][]string{
		1: {"jpg", "png", "jpeg", "svg", "webp", "jfif", "pjpeg", "pjp"},
		2: {"txt", "doc", "docx", "odt", "rtf", "tex", "wks", "wps", "wpd", "pdf"},
		3: {"mp4", "mov", "wmv", "avi", "avchd", "mkv", "mp3", "wav"},
		4: {"zip", "rar", "7z", "tar", "gz", "bz2", "xz", "iso"},
	}
	if len(fileTypes) == 0 {
		return errors.New("File Types Empty")
	}

	//DESC - Create if there is no folder named database.
	if _, err := os.Stat(fmt.Sprintf("./storage/%s", tableName)); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fmt.Sprintf("./storage/%s", tableName), os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	//DESC - The Random creates file name and compares the file types accepted with the file type.
	for _, file := range files {
		name := uuid.New().String()
		fileSplit := strings.Split(file.Filename, ".")
		fileType := fileSplit[len(fileSplit)-1]
		for _, types := range fileTypes {
			ok := globals.InSlice(strings.ToLower(fileType), acceptedTypes[types])
			if ok {
				expectedFileTypes = append(expectedFileTypes, types)
			}
		}
		file.Filename = name + "." + fileType
		savedFiles = append(savedFiles, file)
	}

	//DESC - Register if the number of files is equal with the right file types
	if len(expectedFileTypes) == len(savedFiles) {
		for idx, file := range savedFiles {
			err := c.SaveFile(file, fmt.Sprintf("./storage/%s/%s", tableName, file.Filename))
			if err != nil {
				fileses.ErrFile = append(fileses.ErrFile, file)
				deleteErr := fileses.DeleteFile()
				fileses.File = nil

				if deleteErr != nil {
					log.Println(deleteErr)
				}
				return err
			}
			fileses.File = append(fileses.File, File{
				URL:      fmt.Sprintf("/storage/%s/%s", tableName, file.Filename),
				DataType: expectedFileTypes[idx],
			})
		}

	} else {
		return errors.New("File type is not suitable")
	}
	return nil

}
func (fileses *Files) DeleteFile() error {
	for _, fileURL := range fileses.File {
		err := os.Remove(fmt.Sprintf(".%s", fileURL.URL))
		if err != nil {
			return err
		}
	}
	return nil
}
