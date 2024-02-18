package docs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mshrtsr/gh-iteration/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type AdditionalInformation struct {
	Description  string
	Installation string
	Examples     string
	Links        []Link
}

type Link struct {
	Label string
	URL   string
}

var ErrDirEmptyString = errors.New("dir must not be empty")

func GenMarkdownTree(cmd *cobra.Command, info AdditionalInformation, dir string) error {
	if len(dir) == 0 {
		return ErrDirEmptyString
	}

	err := os.MkdirAll(dir, 0o755) //nolint:gomnd
	if err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	err = genIndexMarkdown(cmd, info, dir)
	if err != nil {
		return err
	}

	err = doc.GenMarkdownTree(cmd, dir)
	if err != nil {
		return fmt.Errorf("failed to generate markdown docs: %w", err)
	}
	return nil
}

func genIndexMarkdown(cmd *cobra.Command, info AdditionalInformation, dir string) error {
	filename := filepath.Join(dir, "index.md")

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	err = genIndexMarkdownCustom(cmd, info, file)
	if err != nil {
		return err
	}
	return nil
}

//nolint:funlen,cyclop
func genIndexMarkdownCustom(cmd *cobra.Command, info AdditionalInformation, writer io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	// Heading
	name := cmd.CommandPath()
	buf.WriteString(fmt.Sprintf("## %s\n\n", name))

	// Index
	buf.WriteString("### Index\n\n")
	if len(info.Installation) > 0 {
		buf.WriteString("- [Installation](#installation)\n")
	}
	buf.WriteString("- [Commands](#commands) \n")
	if len(info.Examples) > 0 {
		buf.WriteString("- [Examples](#examples)\n")
	}
	if len(info.Links) > 0 {
		buf.WriteString("- [Links](#links) \n")
	}
	buf.WriteString("\n")

	// Installation
	if len(info.Installation) > 0 {
		buf.WriteString("### Installation\n\n")
		buf.WriteString(info.Installation)
		buf.WriteString("\n")
	}

	// Commands
	buf.WriteString("### Commands\n\n")
	buf.WriteString("|Command|Description|\n")
	buf.WriteString("|-|-|\n")
	for _, subCmd := range cmd.Commands() {
		if !subCmd.IsAvailableCommand() || subCmd.IsAdditionalHelpTopicCommand() {
			continue
		}
		cName := subCmd.CommandPath()
		link := strings.ReplaceAll(cName, " ", "_") + ".md"
		buf.WriteString(fmt.Sprintf("|[%s](%s)|%s|\n", cName, link, subCmd.Short))
	}
	buf.WriteString("\n")

	// Examples
	if len(info.Examples) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(info.Examples)
		buf.WriteString("\n")
	}

	// Links
	if len(info.Links) > 0 {
		buf.WriteString("### Links\n\n")
		for _, link := range info.Links {
			buf.WriteString(fmt.Sprintf("- [%s](%s)\n", link.Label, link.URL))
		}
		buf.WriteString("\n")
	}

	// Footer
	if !cmd.DisableAutoGenTag {
		buf.WriteString(fmt.Sprintf("###### Generated on %s\n\n", time.Now().Format("2-Jan-2006")))
	}

	log.Debug(buf.String())
	_, err := buf.WriteTo(writer)
	if err != nil {
		return fmt.Errorf("failed to write docs to writer: %w", err)
	}
	return nil
}
