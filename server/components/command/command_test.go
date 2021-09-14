package command

import "testing"

func TestTransformExec(t *testing.T) {
	exec, err := TransformExec("ECHO")
	if err != nil {
		t.Fatalf(`TransformExec("ECHO") throws exception: %s`, err)
	}
	if exec != EchoExec {
		t.Fatalf(`TransformExec("ECHO") doesn't matches %d`, EchoExec)
	}
}
