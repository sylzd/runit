package runit

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New("ls", "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewNoCommand(t *testing.T) {
	_, err := New("", "abc")
	if err == nil {
		t.Fatal("cmd empty should be err")
	}
}

func TestRun(t *testing.T) {
	runner, err := New("echo", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRunKeepAlive(t *testing.T) {
	runner, err := New("echo", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(true)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	runner.Shutdown()
}

func TestKill(t *testing.T) {
	runner, err := New("test/test.sh", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(true)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	err = runner.Kill()
	if err != nil {
		t.Fatal(err)
	}
	runner.Shutdown()
}

func TestRestart(t *testing.T) {
	runner, err := New("test/test.sh", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Restart()
	if err != nil {
		t.Fatal(err)
	}
	runner.Kill()
}

func TestRestartListen(t *testing.T) {
	runner, err := New("echo", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)
	if err != nil {
		t.Fatal(err)
	}
	runner.restartChan <- true
	runner.Kill()
}

func TestWatch(t *testing.T) {
	runner, err := New("echo", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = runner.Run(false)

	// create file
	err = ioutil.WriteFile("test/test.txt", []byte("hello"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing creating test file")
	time.Sleep(1 * time.Second)

	// write file
	err = ioutil.WriteFile("test/test.txt", []byte("goodbye"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing writing to test file")
	time.Sleep(1 * time.Second)

	// rename file
	err = os.Rename("test/test.txt", "test/test1.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing renaming test file")
	time.Sleep(1 * time.Second)

	// remove
	err = os.Remove("test/test1.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing removing test file")
	time.Sleep(1 * time.Second)

	t.Log("cleanup test dir")
	err = os.RemoveAll("test/test")
	if err != nil {
		t.Fatal(err)
	}

	// create dir
	err = os.MkdirAll("test/test", 0777)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("testing creating test dir")
	time.Sleep(1 * time.Second)
	err = os.RemoveAll("test/test")
	if err != nil {
		t.Fatal(err)
	}

	runner.Shutdown()
}

func TestWatchInvalidPath(t *testing.T) {
	_, err := New("echo", "nothere")
	if err == nil {
		t.Fatal("should get no folders to watch here")
	}
}
