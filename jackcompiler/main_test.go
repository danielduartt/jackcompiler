package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListarJack(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "Backup.jack")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	for _, nome := range []string{"Square.jack", "Main.jack", "Main.vm", "Backup.jack/Outro.jack"} {
		if err := os.WriteFile(filepath.Join(dir, nome), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
	mainJack := filepath.Join(dir, "Main.jack")
	for _, tc := range []struct {
		entrada string
		want    []string
	}{
		{dir, []string{mainJack, filepath.Join(dir, "Square.jack")}},
		{mainJack, []string{mainJack}},
	} {
		got, err := listarJack(tc.entrada)
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("listarJack(%q) = %v, %v; esperado %v", tc.entrada, got, err, tc.want)
		}
	}
	if _, err := listarJack(filepath.Join(dir, "Main.vm")); err == nil {
		t.Fatal("arquivo .vm aceito como entrada Jack")
	}
}
