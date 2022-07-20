package frontendmockerfile

import (
	"runtime"
	"time"

	"github.com/moby/buildkit/util/system"
	"github.com/moby/moby/api/types/strslice"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type Image struct {
	specs.Image

	Config ImageConfig `json:"config,omitempty"`
}

type ImageConfig struct {
	specs.ImageConfig

	Healthcheck *HealthConfig `json:",omitempty"`
	ArgsEscaped bool          `json:",omitempty"`

	OnBuild     []string
	StopTimeout *int              `json:",omitempty"`
	Shell       strslice.StrSlice `json:",omitempty"`
}

type HealthConfig struct {
	Test []string `json:",omitempty"`

	Interval    time.Duration `json:",omitempty"`
	Timeout     time.Duration `json:",omitempty"`
	StartPeriod time.Duration `json:"omitempty"`

	Retries int `json:",omitempty"`
}

func NewImageConfig(c *Config) *Image {
	img := emptyImage()
	img.Config.Cmd = []string{"ls"}

	return img
}

func emptyImage() *Image {
	img := &Image{
		Image: specs.Image{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	img.RootFS.Type = "layers"
	img.Config.WorkingDir = "/"
	img.Config.Env = []string{"PATH=" + system.DefaultPathEnv(runtime.GOOS)}

	return img
}

func clone(src Image) Image {
	img := src
	img.Config = src.Config
	img.Config.Env = append([]string{}, src.Config.Env...)
	img.Config.Cmd = append([]string{}, src.Config.Cmd...)
	img.Config.Entrypoint = append([]string{}, src.Config.Entrypoint...)

	return img
}
