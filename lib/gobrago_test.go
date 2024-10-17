package lib_test

import (
	"gobrago/lib"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expandCfgToCmdTests = map[string]struct {
	installCfg lib.GobraInstallCfg
	jobCfg     lib.VerificationJobCfg
	expected   string
}{
	"simple": {
		installCfg: lib.GobraInstallCfg{
			JarPath:    "gobra.jar",
			JvmOptions: []string{"-Xss128m"},
		},
		jobCfg: lib.VerificationJobCfg{
			PackagePath: "/test/",
			Backend:     lib.Silicon,
		},
		expected: "java -Xss128m -jar gobra.jar --noassumeInjectivityOnInhale --backend SILICON --mceMode= -p /test/",
	},
}

func TestExpandCfgToCmd(t *testing.T) {
	for desc, test := range expandCfgToCmdTests {
		cmd, err := lib.ExpandCfgToCmd(test.installCfg, test.jobCfg)
		assert.NoError(t, err, desc)
		assert.Equal(t, test.expected, cmd, desc)
	}
}

var genCmdTests = map[string]struct {
	gobraCfgFile string
	jobCfgFile   string
	expected     string
}{
	"scion path package": {
		gobraCfgFile: "testdata/gobra/gobra_cfg.json",
		jobCfgFile:   "testdata/jobs/scion_path.json",
		expected: "java -Xss1g -Xmx4g -jar /gobra.jar --assumeInjectivityOnInhale " +
			"--backend SILICON --checkConsistency --onlyFilesWithHeader " +
			"-I testdata/jobs testdata/jobs/verification/dependencies --mceMode=on " +
			"-m github.com/scionproto/scion --moreJoins off " +
			"-p testdata/jobs/pkg/slayers/path --parallelizeBranches " +
			"--printVpr --requireTriggers --disableNL",
	},
	"scion slayers package": {
		gobraCfgFile: "testdata/gobra/gobra_cfg.json",
		jobCfgFile:   "testdata/jobs/scion_slayers.json",
		expected: "java -Xss1g -Xmx4g -jar /gobra.jar --assumeInjectivityOnInhale --backend VSWITHSILICON --checkConsistency " +
			"--onlyFilesWithHeader -I testdata/jobs testdata/jobs/verification/dependencies -i testdata/jobs/pkg/slayers/extn.go " +
			"testdata/jobs/pkg/slayers/extn_spec.gobra testdata/jobs/pkg/slayers/l4.go testdata/jobs/pkg/slayers/layertypes_spec.gobra " +
			"testdata/jobs/pkg/slayers/scion.go testdata/jobs/pkg/slayers/scion_spec.gobra testdata/jobs/pkg/slayers/scmp.go " +
			"testdata/jobs/pkg/slayers/scmp_msg.go testdata/jobs/pkg/slayers/scmp_msg_spec.gobra testdata/jobs/pkg/slayers/scmp_spec.gobra " +
			"testdata/jobs/pkg/slayers/scmp_typecode.go testdata/jobs/pkg/slayers/scmp_typecode_spec.gobra --mceMode=on -m " +
			"github.com/scionproto/scion --moreJoins off --parallelizeBranches --printVpr --requireTriggers",
	},
}

func TestGenCmd(t *testing.T) {
	for desc, test := range genCmdTests {
		cmd, err := lib.GenCmd(test.gobraCfgFile, test.jobCfgFile)
		assert.NoError(t, err, desc)
		assert.Equal(t, test.expected, cmd, desc)
	}
}
