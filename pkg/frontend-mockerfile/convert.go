package frontendmockerfile

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

func Mockerfile2LLB(c *Config) (llb.State, *Image) {
	current := c.Images[0]
	s := llb.Image(current.From)

	if current.Package != nil {
		s = packages(s, current.Package)
	}

	for _, e := range current.External {
		downloaded := external(e)
		s = copy(downloaded, e.Destination, s, e.Destination)
	}

	imageCfg := NewImageConfig(c)

	return s, imageCfg
}

func packages(base llb.State, p *Package) llb.State {
	base = base.Run(shf("apt-get update && apt-get install apt-transport-https -y")).Root()
	for _, repo := range p.Repo {
		base = base.Run(shf("echo \"%s\" >> /etc/apt/sources.list", repo)).Root()
	}
	for _, key := range p.Gpg {
		base = aptAddKey(base, key)
	}

	if len(p.Install) > 0 {
		pkgs := strings.Join(p.Install, " ")
		base = base.Run(shf("apt-get update && apt-get install --no-install-recommends --no-install-suggests -y %s", pkgs)).Root()
	}

	return base
}

func external(e *ExternalFile) llb.State {
	downloadDst := e.Destination
	isTarGz := strings.HasSuffix(e.Source, ".tar.gz")
	isExecutable := true
	if isTarGz {
		downloadDst = "tmp.tar.gz"
		isExecutable = false
	}

	curlCmd := fmt.Sprintf("curl -Lo %s %s", downloadDst, e.Source)
	if isExecutable {
		curlCmd = fmt.Sprintf("%s && chmod +x %s", curlCmd, downloadDst)
	}

	s := curl().
		Run(sh(curlCmd)).
		Root()
	if e.Sha256 != "" {
		s = s.Run(shf("echo \"%s %s\" | sha256sum -c -", e.Sha256, downloadDst)).Root()
	}
	if isTarGz {
		s = s.Run(shf("mkdir -p %[2]s && tar -zxvf %[1]s -C %[2]s && rm %[1]s", downloadDst, e.Destination)).Root()
	}
	if len(e.Install) > 0 {
		installCmds := strings.Join(e.Install, " && ")
		s = s.Run(shf(installCmds)).Root()
	}

	return s
}

func shf(cmd string, v ...interface{}) llb.RunOption {
	return llb.Args([]string{"/bin/sh", "-c", fmt.Sprintf(cmd, v...)})
}

func aptAddKey(dst llb.State, url string) llb.State {
	downloadSt := curl().
		Run(llb.Shlexf("curl -Lo /key.gpg %s", url)).
		Root()
	dst = copy(downloadSt, "/key.gpg", dst, "/key.gpg")

	return dst.
		Run(sh("apt-key add /key.gpg && rm /key.gpg")).
		Root()
}

func curl() llb.State {
	return llb.Image("docker.io/library/alpine:3.6").
		Run(llb.Shlexf("apk add --no-cache curl")).
		Root()
}

func copy(src llb.State, srcPath string, dest llb.State, destPath string) llb.State {
	cpImage := llb.Image("docker.io/library/alpine:3.6")
	cp := cpImage.Run(llb.Shlexf("cp -a /src%s /dest%s", srcPath, destPath))
	cp.AddMount("/src", src, llb.Readonly)

	return cp.AddMount("/dest", dest)
}

func sh(cmd string) llb.RunOption {
	return llb.Args([]string{"/bin/sh", "-c", cmd})
}
