package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"path/filepath"
	"os"

)

func Test_commandsExist(t *testing.T) {
	// we expect these subcommands to be handled
	expect := []string{"install", "enable", "disable", "uninstall", "update"}
	for _, c := range expect {
		_, ok := cmds[c]
		if !ok {
			t.Fatalf("cmd '%s' is not handled", c)
		}
	}
}

func Test_commands_shouldReportStatus(t *testing.T) {
	// - certain extension invocations are supposed to write 'N.status' files and some do not.

	// these subcommands should NOT report status
	require.False(t, cmds["install"].shouldReportStatus, "install should not report status")
	require.False(t, cmds["uninstall"].shouldReportStatus, "uninstall should not report status")

	// these subcommands SHOULD report status
	require.True(t, cmds["enable"].shouldReportStatus, "enable should report status")
	require.True(t, cmds["disable"].shouldReportStatus, "disable should report status")
	require.True(t, cmds["update"].shouldReportStatus, "update should report status")
}

func Test_checkAndSaveSeqNum_fails(t *testing.T) {
	// pass in invalid seqnum format
	_, err := checkAndSaveSeqNum(log.NewNopLogger(), 0, "/non/existing/dir")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), `failed to save the sequence number`)
}

func Test_checkAndSaveSeqNum(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	fp := filepath.Join(dir, "seqnum")
	defer os.RemoveAll(dir)

	nop := log.NewNopLogger()

	// no sequence number, 0 comes in.
	shouldExit, err := checkAndSaveSeqNum(nop, 0, fp)
	require.Nil(t, err)
	require.False(t, shouldExit)

	// file=0, seq=0 comes in. (should exit)
	shouldExit, err = checkAndSaveSeqNum(nop, 0, fp)
	require.Nil(t, err)
	require.True(t, shouldExit)

	// file=0, seq=1 comes in.
	shouldExit, err = checkAndSaveSeqNum(nop, 1, fp)
	require.Nil(t, err)
	require.False(t, shouldExit)

	// file=1, seq=1 comes in. (should exit)
	shouldExit, err = checkAndSaveSeqNum(nop, 1, fp)
	require.Nil(t, err)
	require.True(t, shouldExit)

	// file=1, seq=0 comes in. (should exit)
	shouldExit, err = checkAndSaveSeqNum(nop, 1, fp)
	require.Nil(t, err)
	require.True(t, shouldExit)
}

func Test_runCmd_success(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	require.Nil(t, runCmd(log.NewNopLogger(), dir, handlerSettings{
		publicSettings: publicSettings{CommandToExecute: "date"},
	}), "command should run successfully")

	// check stdout stderr files
	_, err = os.Stat(filepath.Join(dir, "stdout"))
	require.Nil(t, err, "stdout should exist")
	_, err = os.Stat(filepath.Join(dir, "stderr"))
	require.Nil(t, err, "stderr should exist")
}

func Test_runCmd_pass(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	err = runCmd(log.NewNopLogger(), dir, handlerSettings{})
	require.Nil(t, err)
}

func Test_unzip_fail(t *testing.T) {
	_, err := unzip(log.NewNopLogger(), "","agent")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), `failed to open zip`)

	_, err = unzip(log.NewNopLogger(), "hello.zip","agent")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), `failed to open zip`)
}

func Test_unzip_pass(t *testing.T) {
	filenames, err := unzip(log.NewNopLogger(), "../agent/GCAgentx64.zip", "agent")
	require.Nil(t, err)
	require.NotEmpty(t, filenames)
	t.Log(filenames)
}

/*
func Test_downloadFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	srv := httptest.NewServer(httpbin.GetMux())
	defer srv.Close()

	err = downloadFiles(log.NewContext(log.NewNopLogger()),
		dir,
		handlerSettings{
			publicSettings: publicSettings{
				FileURLs: []string{
					srv.URL + "/bytes/10",
					srv.URL + "/bytes/100",
					srv.URL + "/bytes/1000",
				}},
		})
	require.Nil(t, err)

	// check the files
	f := []string{"10", "100", "1000"}
	for _, fn := range f {
		fp := filepath.Join(dir, fn)
		_, err := os.Stat(fp)
		require.Nil(t, err, "%s is missing from download dir", fp)
	}
} */