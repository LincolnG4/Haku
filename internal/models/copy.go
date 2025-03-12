package models

import (
	"fmt"
	"io"
	"os"
)

type CopyFileActivity struct{}

func (c *CopyFileActivity) Execute(task Task) error {
	src := task.Config["source"].(string)
	dest := task.Config["destination"].(string)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	fmt.Println("File copied successfully")
	return nil
}
