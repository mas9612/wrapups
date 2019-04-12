package command

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/mas9612/wrapups/pkg/wrapups"
)

func printWrapup(doc *pb.Wrapup) {
	fmt.Printf("ID: %s\n", doc.Id)
	fmt.Printf("Title: %s\n", doc.Title)
	fmt.Printf("Wrapup: %s\n", doc.Wrapup)
	fmt.Printf("Comment: %s\n", doc.Comment)
	fmt.Printf("Note: %s\n", doc.Note)
	t, err := ptypes.Timestamp(doc.CreateTime)
	if err != nil {
		fmt.Printf("CreateTime: <invalid>\n")
	} else {
		fmt.Printf("CreateTime: %s\n", t.String())
	}
}
