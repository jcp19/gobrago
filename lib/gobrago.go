package lib

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	. "github.com/jcp19/goprelude/utils"
)

type (
	MceMode      string
	MoreJoins    string
	ViperBackend string
)

const (
	MceOn  MceMode = "on"
	MceOd  MceMode = "od"
	MceOff MceMode = "off"

	MoreJoinsAll    MoreJoins = "all"
	MoreJoinsImpure MoreJoins = "impure"
	MoreJoinsOff    MoreJoins = "off"

	Silicon            ViperBackend = "SILICON"
	Carbon             ViperBackend = "CARBON"
	SiliconViperServer ViperBackend = "VSWITHSILICON"
	CarbonViperServer  ViperBackend = "VSWITHCARBON"
)

type GobraInstallCfg struct {
	JarPath    string   `json:"jar_path"`
	JvmOptions []string `json:"jvm_options"`
	Z3Path     string   `json:"z3_path"`
}

type VerificationJobCfg struct {
	AssumeInjectivityInhale   bool         `json:"assume_injectivity_inhale"`
	Backend                   ViperBackend `json:"backend"`
	CheckConsistency          bool         `json:"check_consistency"`
	CheckOverflow             bool         `json:"overflow"`
	ConditionalizePermissions bool         `json:"conditionalize_permissions"`
	HeaderOnly                bool         `json:"only_files_with_header"`
	IncludePaths              []string     `json:"includes"`
	InputFilePaths            []string     `json:"input_files"`
	Mce                       MceMode      `json:"mce_mode"`
	Module                    string       `json:"module"`
	MoreJoins                 MoreJoins    `json:"more_joins"`
	PackagePath               string       `json:"pkg_path"`
	ParallelizeBranches       bool         `json:"parallelize_branches"`
	PrintVpr                  bool         `json:"print_vpr"`
	ProjectRoot               string       `json:"project_root"`
	Recursive                 bool         `json:"recursive"`
	RequireTriggers           bool         `json:"require_triggers"`
	OtherFlags                []string     `json:"other"`
}

func DefaultGobraInstallCfg() GobraInstallCfg {
	return GobraInstallCfg{}
}

func DefaultVerificationJobCfg() VerificationJobCfg {
	return VerificationJobCfg{
		AssumeInjectivityInhale:   true,
		Backend:                   Silicon,
		CheckConsistency:          true,
		CheckOverflow:             false,
		ConditionalizePermissions: false,
		HeaderOnly:                true,
		Mce:                       MceOd,
		ParallelizeBranches:       false,
		PrintVpr:                  false,
		Recursive:                 false,
		RequireTriggers:           true,
	}
}

func (g *GobraInstallCfg) ResolvePaths(basePath string) error {
	if !filepath.IsAbs(g.JarPath) {
		g.JarPath = filepath.Join(basePath, g.JarPath)
	}
	if g.Z3Path != "" && !filepath.IsAbs(g.Z3Path) {
		g.Z3Path = filepath.Join(basePath, g.Z3Path)
	}
	return nil
}

func (v *VerificationJobCfg) ResolvePaths(basePath string) error {
	for i, val := range v.IncludePaths {
		if !filepath.IsAbs(val) {
			v.IncludePaths[i] = filepath.Join(basePath, val)
		}
	}
	for i, val := range v.InputFilePaths {
		if !filepath.IsAbs(val) {
			v.InputFilePaths[i] = filepath.Join(basePath, val)
		}
	}
	if v.PackagePath != "" && !filepath.IsAbs(v.PackagePath) {
		v.PackagePath = filepath.Join(basePath, v.PackagePath)
	}
	if v.ProjectRoot != "" && !filepath.IsAbs(v.ProjectRoot) {
		v.ProjectRoot = filepath.Join(basePath, v.ProjectRoot)
	}
	return nil
}

func ExpandCfgToCmd(installCfg GobraInstallCfg, jobCfg VerificationJobCfg) (string, error) {
	components := make([]string, 0, 1024)
	// expand gobra install config:
	Append(&components, "java")
	Append(&components, installCfg.JvmOptions...)
	if installCfg.JarPath == "" {
		return "", errors.New("path to Gobra's jar was not provided")
	}
	Append(&components, "-jar", installCfg.JarPath)
	if installCfg.Z3Path != "" {
		Append(&components, "--z3Exe", installCfg.Z3Path)
	}

	// expand verification job config:
	if jobCfg.AssumeInjectivityInhale {
		Append(&components, "--assumeInjectivityOnInhale")
	} else {
		Append(&components, "--noassumeInjectivityOnInhale")
	}
	if jobCfg.Backend != "" {
		Append(&components, "--backend", string(jobCfg.Backend))
	}
	if jobCfg.CheckConsistency {
		Append(&components, "--checkConsistency")
	}
	if jobCfg.CheckOverflow {
		Append(&components, "--overflow")
	}
	if jobCfg.ConditionalizePermissions {
		Append(&components, "--conditionalizePermissions")
	}
	if jobCfg.HeaderOnly {
		Append(&components, "--onlyFilesWithHeader")
	}
	if len(jobCfg.IncludePaths) > 0 {
		Append(&components, "-I")
		Append(&components, jobCfg.IncludePaths...)
	}
	if len(jobCfg.InputFilePaths) > 0 {
		Append(&components, "-i")
		Append(&components, jobCfg.InputFilePaths...)
	}
	Append(&components, "--mceMode="+string(jobCfg.Mce))
	if jobCfg.Module != "" {
		Append(&components, "-m", jobCfg.Module)
	}
	if jobCfg.MoreJoins != "" {
		Append(&components, "--moreJoins", string(jobCfg.MoreJoins))
	}
	if jobCfg.PackagePath != "" {
		Append(&components, "-p", jobCfg.PackagePath)
	}
	if jobCfg.ParallelizeBranches {
		Append(&components, "--parallelizeBranches")
	}
	if jobCfg.PrintVpr {
		Append(&components, "--printVpr")
	}
	if jobCfg.ProjectRoot != "" {
		Append(&components, "--projectRoot", jobCfg.ProjectRoot)
	}
	if jobCfg.Recursive {
		Append(&components, "-r")
	}
	if jobCfg.RequireTriggers {
		Append(&components, "--requireTriggers")
	}
	Append(&components, jobCfg.OtherFlags...)

	return strings.Join(components, " "), nil
}

func GenCmd(installCfgPath, jobCfgPath string) (string, error) {
	installCfgRaw, err := os.ReadFile(installCfgPath)
	installCfgDir := filepath.Dir(installCfgPath)
	if err != nil {
		return "", err
	}
	jobCfgRaw, err := os.ReadFile(jobCfgPath)
	jobCfgDir := filepath.Dir(jobCfgPath)
	if err != nil {
		return "", err
	}
	// set defaults before unmarshalling
	installCfg := DefaultGobraInstallCfg()
	jobCfg := DefaultVerificationJobCfg()
	if err := json.Unmarshal(installCfgRaw, &installCfg); err != nil {
		return "", err
	}
	if err := json.Unmarshal(jobCfgRaw, &jobCfg); err != nil {
		return "", err
	}
	if err := installCfg.ResolvePaths(installCfgDir); err != nil {
		return "", err
	}
	if err := jobCfg.ResolvePaths(jobCfgDir); err != nil {
		return "", err
	}
	return ExpandCfgToCmd(installCfg, jobCfg)
}
