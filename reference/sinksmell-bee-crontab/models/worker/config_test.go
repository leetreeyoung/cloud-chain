package worker

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	err := InitConfig("../../conf/worker.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", Conf)
}
