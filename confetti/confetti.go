package confetti

import (
	"math/rand"
	"time"

	"github.com/maaslalani/confetty/array"
	"github.com/maaslalani/confetty/simulation"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	framesPerSecond = 30.0
	numParticles    = 75
)

var (
	colors     = []string{"#a864fd", "#29cdff", "#78ff44", "#ff718d", "#fdff6a"}
	characters = []string{"█", "▓", "▒", "░", "▄", "▀"}
)

type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/framesPerSecond, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

// Confetti model
type model struct {
	system  *simulation.System
	counter int
}

func Spawn(width, height int) []*simulation.Particle {
	particles := []*simulation.Particle{}
	for i := 0; i < numParticles; i++ {
		x := float64(width / 2)
		y := float64(0)

		p := simulation.Particle{
			Physics: harmonica.NewProjectile(
				harmonica.FPS(framesPerSecond),
				harmonica.Point{X: x + (float64(width/4) * (rand.Float64() - 0.5)), Y: y, Z: 0},
				harmonica.Vector{X: (rand.Float64() - 0.5) * 100, Y: rand.Float64() * 50, Z: 0},
				harmonica.TerminalGravity,
			),
			Char: lipgloss.NewStyle().
				Foreground(lipgloss.Color(array.Sample(colors))).
				Render(array.Sample(characters)),
		}

		particles = append(particles, &p)
	}
	return particles
}

func InitialModel() model {
	return model{system: &simulation.System{
		Particles: []*simulation.Particle{},
		Frame:     simulation.Frame{},
	}}
}

// Init initializes the confetti after a small delay
func (m model) Init() tea.Cmd {
	return animate()
}

// Update updates the model every frame, it handles the animation loop and
// updates the particle physics every frame
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		return m, nil
	case frameMsg:
		m.counter = (m.counter + 1) % 18
		if m.counter == 0 {
			m.system.Particles = append(m.system.Particles, Spawn(m.system.Frame.Width, m.system.Frame.Height)...)
		}
		m.system.Update()
		return m, animate()
	case tea.WindowSizeMsg:
		if m.system.Frame.Width == 0 && m.system.Frame.Height == 0 {
			// For the first frameMsg spawn a system of particles
			m.system.Particles = Spawn(msg.Width, msg.Height)
		}
		m.system.Frame.Width = msg.Width
		m.system.Frame.Height = msg.Height
		return m, nil
	default:
		return m, nil
	}
}

// View displays all the particles on the screen
func (m model) View() string {
	return m.system.Render()
}
