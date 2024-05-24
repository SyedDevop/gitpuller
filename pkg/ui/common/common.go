package common

// TODO: Add Features.
// - mouse:(zone "github.com/lrstanley/bubblezone")
// - highlight: ("github.com/alecthomas/chroma/v2/lexers")
import (
	"context"
	"io"

	"github.com/SyedDevop/gitpuller/pkg/ui/keymap"
	"github.com/SyedDevop/gitpuller/pkg/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

type contextKey struct {
	name string
}

// Keys to use for context.Context.
var (
	ConfigKey = &contextKey{"config"}
	RepoKey   = &contextKey{"repo"}
)

// Common is a struct all components should embed.
type Common struct {
	ctx          context.Context
	Styles       *styles.Styles
	KeyMap       *keymap.KeyMap
	Renderer     *lipgloss.Renderer
	Output       *termenv.Output
	Logger       *log.Logger
	Width        int
	Height       int
	HideCloneCmd bool
}

// NewCommon returns a new Common struct.
func NewCommon(ctx context.Context, fileLogger io.Writer, out *lipgloss.Renderer, width, height int) Common {
	if ctx == nil {
		ctx = context.TODO()
	}
	logger := log.FromContext(ctx).WithPrefix("ui")
	logger.SetOutput(fileLogger)
	return Common{
		ctx:      ctx,
		Width:    width,
		Height:   height,
		Renderer: out,
		Output:   out.Output(),
		Styles:   styles.DefaultStyles(out),
		KeyMap:   keymap.DefaultKeyMap(),
		Logger:   logger,
	}
}

// SetValue sets a value in the context.
func (c *Common) SetValue(key, value interface{}) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

// SetSize sets the width and height of the common struct.
func (c *Common) SetSize(width, height int) {
	c.Width = width
	c.Height = height
}

// Context returns the context.
func (c *Common) Context() context.Context {
	return c.ctx
}

// Config returns the server config.
// func (c *Common) Config() *config.Config {
// 	return config.FromContext(c.ctx)
// }

// Backend returns the Soft Serve backend.
// func (c *Common) Backend() *backend.Backend {
// 	return backend.FromContext(c.ctx)
// }

// Repo returns the repository.
// func (c *Common) Repo() *git.Repository {
// 	v := c.ctx.Value(RepoKey)
// 	if r, ok := v.(*git.Repository); ok {
// 		return r
// 	}
// 	return nil
// }

// PublicKey returns the public key.
// func (c *Common) PublicKey() ssh.PublicKey {
// 	v := c.ctx.Value(ssh.ContextKeyPublicKey)
// 	if p, ok := v.(ssh.PublicKey); ok {
// 		return p
// 	}
// 	return nil
// }

// CloneCmd returns the clone command string.
// func (c *Common) CloneCmd(publicURL, name string) string {
// 	if c.HideCloneCmd {
// 		return ""
// 	}
// 	return fmt.Sprintf("git clone %s", RepoURL(publicURL, name))
// }

// IsFileMarkdown returns true if the file is markdown.
// It uses chroma lexers to analyze and determine the language.
// func IsFileMarkdown(content, ext string) bool {
// 	var lang string
// 	lexer := lexers.Match(ext)
// 	if lexer == nil {
// 		lexer = lexers.Analyse(content)
// 	}
// 	if lexer != nil && lexer.Config() != nil {
// 		lang = lexer.Config().Name
// 	}
// 	return lang == "markdown"
// }
