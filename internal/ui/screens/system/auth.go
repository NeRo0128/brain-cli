package system

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/NeRo0128/brain-cli/internal/core/auth"
	"github.com/NeRo0128/brain-cli/internal/ui/keys"
	"github.com/NeRo0128/brain-cli/internal/ui/screens"
	"github.com/NeRo0128/brain-cli/internal/ui/styles"
	authuc "github.com/NeRo0128/brain-cli/internal/usecases/auth"
	"github.com/rs/zerolog"
)

type authState int

const (
	stateLoading authState = iota
	stateIdle
	statePending
	stateConnected
	stateError
)

type authStatusMsg struct {
	user *auth.User
	err  error
}

type authDoneMsg struct {
	user *auth.User
	err  error
}

type AuthScreen struct {
	mgr      *authuc.Manager
	log      zerolog.Logger
	styles   *styles.Styles
	state    authState
	user     *auth.User
	code     *auth.AuthCode
	err      error
	spinner  spinner.Model
	width    int
	height   int
	clientID string
}

func NewAuthScreen(
	mgr *authuc.Manager,
	log zerolog.Logger,
	s *styles.Styles,
) AuthScreen {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = s.SpinnerStyle

	clientID := ""
	if mgr != nil {
		clientID = "configured"
	}

	return AuthScreen{
		mgr:      mgr,
		log:      log.With().Str("screen", "auth").Logger(),
		styles:   s,
		state:    stateLoading,
		spinner:  sp,
		clientID: clientID,
	}
}

func (m AuthScreen) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return m.spinner.Tick() },
		func() tea.Msg { return tea.RequestWindowSize() },
		m.loadSession(),
	)
}

func (m AuthScreen) loadSession() tea.Cmd {
	return func() tea.Msg {
		if m.mgr == nil {
			return authStatusMsg{err: nil}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := m.mgr.CurrentUser(ctx)
		if err != nil {
			return authStatusMsg{err: err}
		}
		return authStatusMsg{user: user}
	}
}

func (m AuthScreen) Keys() []string {
	return []string{keys.ViewAuth, keys.NavBack, keys.ViewHelp, keys.ActionExecute}
}

func (m AuthScreen) Update(msg tea.Msg) (screens.ScreenI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case authStatusMsg:
		if msg.user != nil {
			m.user = msg.user
			m.state = stateConnected
		} else if m.clientID == "" {
			m.state = stateError
			m.err = &configError{msg: "GitHub no configurado. Añade github.client_id en config.yaml"}
		} else {
			m.state = stateIdle
		}
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if m.state == stateIdle && m.clientID != "" {
				return m.startLogin()
			}
			if m.state == stateError {
				return m.startLogin()
			}

		case "esc":
			if m.state == statePending {
				m.mgr.CancelLogin()
				m.state = stateIdle
				m.code = nil
				m.err = nil
				return m, nil
			}
			return m, screens.Back()

		case "ctrl+g":
			return m, nil
		}

	case screens.ActionMsg:
		switch msg.ID {
		case keys.NavBack:
			if m.state == statePending {
				m.mgr.CancelLogin()
				m.state = stateIdle
				m.code = nil
				m.err = nil
				return m, nil
			}
			return m, screens.Back()

		case keys.ActionExecute:
			if m.state == stateIdle || m.state == stateError {
				return m.startLogin()
			}
			return m, nil

		case keys.ViewHelp:
			return m, screens.OpenHelp()
		}

	case authDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			m.code = nil
		} else {
			m.user = msg.user
			m.state = stateConnected
			m.code = nil
		}
		return m, nil
	}

	return m, nil
}

func (m AuthScreen) startLogin() (screens.ScreenI, tea.Cmd) {
	m.state = statePending
	var err error
	m.code, err = m.mgr.StartLogin(context.Background())
	if err != nil {
		m.err = err
		m.state = stateError
		return m, nil
	}

	// Abrir el navegador en background. Si falla, mostramos la URL
	// igual en pantalla (el usuario puede abrirla manualmente).
	if m.code.VerificationURL != "" {
		if err := openBrowser(m.code.VerificationURL); err != nil {
			m.log.Warn().Err(err).Msg("no se pudo abrir el navegador")
		}
	}

	return m, m.waitForToken()
}

func (m AuthScreen) waitForToken() tea.Cmd {
	return func() tea.Msg {
		if m.code == nil {
			return authDoneMsg{err: authuc.ErrNoPendingLogin}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		session, err := m.mgr.CompleteLogin(ctx)
		if err != nil {
			return authDoneMsg{err: err}
		}
		return authDoneMsg{user: &session.User}
	}
}

func (m AuthScreen) View() tea.View {

	switch m.state {
	case stateLoading:
		return tea.NewView(m.centered(m.spinner.View() + " Comprobando sesión..."))

	case stateIdle:
		if m.clientID == "" {
			return tea.NewView(m.renderError("GitHub no configurado", "Añade github.client_id en config.yaml\n(o exporta GITHUB_CLIENT_ID)"))
		}
		return tea.NewView(m.renderIdle())

	case statePending:
		return tea.NewView(m.renderPending())

	case stateConnected:
		return tea.NewView(m.renderConnected())

	case stateError:
		return tea.NewView(m.renderError("Error", m.err.Error()))
	}

	return tea.NewView("")
}

func (m AuthScreen) renderIdle() string {
	var b strings.Builder

	b.WriteString(m.styles.Title.Render("🔐 GitHub"))
	b.WriteString("\n\n")

	b.WriteString(m.styles.Subtitle.Render("Estado:  ⚠  No conectado"))
	b.WriteString("\n\n")
	b.WriteString("Pulsa ")
	b.WriteString(m.styles.Key.Render("Enter"))
	b.WriteString(" para conectar con GitHub.")

	return m.centeredBox(b.String())
}

func (m AuthScreen) renderPending() string {
	p := m.styles.Theme.Resolve(m.styles.Dark)
	if m.code == nil {
		return m.centered(m.spinner.View() + " Iniciando...")
	}

	var b strings.Builder

	b.WriteString(m.styles.Title.Render("🔐 GitHub"))
	b.WriteString("\n\n")

	boxW := 60
	if m.width > 0 && m.width-10 < boxW {
		boxW = m.width - 10
	}
	if boxW < 40 {
		boxW = 40
	}

	userCode := m.formatUserCode(m.code.UserCode)

	content := lipgloss.NewStyle().Width(boxW - 4).Render(
		m.styles.Subtitle.Render("1. Visita:") + "\n" +
			m.styles.Key.Render(m.code.VerificationURL) + "\n\n" +
			m.styles.Subtitle.Render("2. Introduce el código:") + "\n" +
			m.styles.Key.Render(userCode) + "\n\n" +
			m.spinner.View() + " " + m.styles.Subtitle.Render("Esperando autorización..."),
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.Primary).
		Padding(1, 2).
		Width(boxW).
		Render(content)

	b.WriteString(box)
	b.WriteString("\n\n")
	b.WriteString("Pulsa ")
	b.WriteString(m.styles.Key.Render("Esc"))
	b.WriteString(" para cancelar")

	return m.centeredBox(b.String())
}

func (m AuthScreen) formatUserCode(code string) string {
	clean := strings.Map(func(r rune) rune {
		switch r {
		case '-', ' ', '\t', '\n':
			return -1 // descartar
		}
		return r
	}, code)

	if len(clean) >= 8 {
		return clean[:4] + "-" + clean[4:8]
	}
	return clean
}

func (m AuthScreen) renderConnected() string {
	var b strings.Builder

	b.WriteString(m.styles.Title.Render("🔐 GitHub"))
	b.WriteString("\n\n")

	b.WriteString(m.styles.SuccessStyle.Render("Estado:  ✅ Conectado como @" + m.user.Login))
	b.WriteString("\n")
	if m.user.Name != "" {
		b.WriteString(m.styles.Subtitle.Render("Nombre:  " + m.user.Name))
		b.WriteString("\n")
	}
	if m.user.Email != "" {
		b.WriteString(m.styles.Subtitle.Render("Email:   " + m.user.Email))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString("Pulsa ")
	b.WriteString(m.styles.Key.Render("Esc"))
	b.WriteString(" para volver")

	return m.centeredBox(b.String())
}

func (m AuthScreen) renderError(title, msg string) string {
	p := m.styles.Theme.Resolve(m.styles.Dark)
	var b strings.Builder

	b.WriteString(m.styles.Title.Render("🔐 GitHub"))
	b.WriteString("\n\n")

	boxW := 60
	if m.width > 0 && m.width-10 < boxW {
		boxW = m.width - 10
	}
	if boxW < 40 {
		boxW = 40
	}

	content := m.styles.ErrorStyle.Render("❌ "+title) + "\n\n" + msg

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.Error).
		Padding(1, 2).
		Width(boxW).
		Render(content)

	b.WriteString(box)
	b.WriteString("\n\n")

	if m.clientID != "" {
		b.WriteString("Pulsa ")
		b.WriteString(m.styles.Key.Render("Enter"))
		b.WriteString(" para reintentar · ")
		b.WriteString(m.styles.Key.Render("Esc"))
		b.WriteString(" para volver")
	} else {
		b.WriteString("Pulsa ")
		b.WriteString(m.styles.Key.Render("Esc"))
		b.WriteString(" para volver")
	}

	return m.centeredBox(b.String())
}

func (m AuthScreen) centered(s string) string {
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
	}
	return "\n" + s + "\n"
}

func (m AuthScreen) centeredBox(s string) string {
	if m.width > 0 && m.height > 0 {
		h := max(m.height-4, 12)
		return lipgloss.Place(m.width, h, lipgloss.Center, lipgloss.Center, s)
	}
	return "\n" + s + "\n"
}

type configError struct {
	msg string
}

func (e *configError) Error() string {
	return e.msg
}

// openBrowser intenta abrir la URL en el navegador por defecto.
// Best-effort: si falla, no rompe el login.
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("SO no soportado: %s", runtime.GOOS)
	}
	return cmd.Start()
}
