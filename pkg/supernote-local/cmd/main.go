package main

import (
	"context"
	"fmt"
	"os"

	supernote "github.com/dylanmazurek/supernote-sync/pkg/supernote-local"
	"github.com/markkurossi/tabulate"
)

func main() {
	ctx := context.Background()

	opts := []supernote.Option{
		supernote.WithHost(os.Getenv("HOST")),
		supernote.WithPort(8089),
	}

	supernoteClient, err := supernote.New(ctx, opts...)
	if err != nil {
		panic(err)
	}

	tab := tabulate.New(tabulate.CompactUnicode)
	tab.Header("Name")
	tab.Header("Modified")
	tab.Header("Size")
	tab.Header("Depth")

	entries, err := supernoteClient.ListEntries("/", nil)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		row := tab.Row()
		row.Column(entry.Name)
		row.Column(entry.Date)
		row.Column(fmt.Sprintf("%.2f", entry.Size))
		row.Column(fmt.Sprintf("%d", entry.Depth))
	}

	fmt.Println(tab.String())
}
