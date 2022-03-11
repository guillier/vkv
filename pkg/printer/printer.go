package printer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

type outputFormat int

const (
	maskChar = "*"

	// MaxValueLength maximum length of passwords.
	MaxValueLength = 12

	yaml outputFormat = iota
	json
	export
)

var defaultWriter = os.Stdout

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	secrets     map[string]interface{}
	format      outputFormat
	writer      io.Writer
	onlyKeys    bool
	onlyPaths   bool
	showSecrets bool
	valueLength int
}

// CustomValueLength option for trimming down the output of secrets.
func CustomValueLength(length int) Option {
	return func(p *Printer) {
		p.valueLength = length
	}
}

// OnlyKeys flag for only showing secrets keys.
func OnlyKeys(b bool) Option {
	return func(p *Printer) {
		if b {
			p.onlyKeys = true
			p.printOnlykeys()
		}
	}
}

// OnlyPaths flag for only printing kv paths.
func OnlyPaths(b bool) Option {
	return func(p *Printer) {
		if b {
			p.onlyPaths = true
			p.printOnlyPaths()
		}
	}
}

// ToYAML outputformat to yaml.
func ToYAML(b bool) Option {
	return func(p *Printer) {
		if b {
			p.format = yaml
		}
	}
}

// ToJSON outputformat to yaml.
func ToJSON(b bool) Option {
	return func(p *Printer) {
		if b {
			p.format = json
		}
	}
}

// ToExportFormat option for printing out variables so they can be exported into the shell.
func ToExportFormat(b bool) Option {
	return func(p *Printer) {
		p.format = export
		p.showSecrets = true
	}
}

// WithWriter option for passing a custom io.Writer.
func WithWriter(w io.Writer) Option {
	return func(p *Printer) {
		p.writer = w
	}
}

// ShowSecrets flag for unmasking secrets in output.
func ShowSecrets(b bool) Option {
	return func(p *Printer) {
		p.showSecrets = b
	}
}

// NewPrinter return a new printer struct.
func NewPrinter(m map[string]interface{}, opts ...Option) *Printer {
	p := &Printer{
		secrets:     m,
		writer:      defaultWriter,
		valueLength: MaxValueLength,
	}

	for _, opt := range opts {
		opt(p)
	}

	if !p.showSecrets {
		p.maskSecrets()
	}

	return p
}

// Out prints out the secrets according all configured options.
func (p *Printer) Out() error {
	switch p.format {
	case yaml:
		out, err := utils.ToYAML(p.secrets)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	case json:
		out, err := utils.ToJSON(p.secrets)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	case export:
		for _, s := range utils.SortMapKeys(p.secrets) {
			for k, v := range p.secrets[s].(map[string]interface{}) {
				fmt.Fprintf(p.writer, "export %s=%v\n", k, v)
			}
		}
	default:
		for _, k := range utils.SortMapKeys(p.secrets) {
			fmt.Fprintf(p.writer, "%s\n", k)
			p.printSecrets(p.secrets[k])
		}
	}

	return nil
}

func (p *Printer) printOnlykeys() {
	for k := range p.secrets {
		m, ok := p.secrets[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			m[k] = ""
		}
	}
}

func (p *Printer) printOnlyPaths() {
	for k := range p.secrets {
		p.secrets[k] = nil
	}
}

func (p *Printer) maskSecrets() {
	for k := range p.secrets {
		m, ok := p.secrets[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			secret := fmt.Sprintf("%v", m[k])

			if len(secret) > p.valueLength {
				m[k] = strings.Repeat(maskChar, p.valueLength)
			} else {
				m[k] = strings.Repeat(maskChar, len(secret))
			}
		}
	}
}

func (p *Printer) printSecrets(s interface{}) {
	m, ok := s.(map[string]interface{})
	if ok {
		for _, k := range utils.SortMapKeys(m) {
			if p.onlyKeys {
				fmt.Fprintf(p.writer, "\t%s\n", k)
			}

			if !p.onlyKeys && !p.onlyPaths {
				fmt.Fprintf(p.writer, "\t%s=%v\n", k, m[k])
			}
		}
	}
}
