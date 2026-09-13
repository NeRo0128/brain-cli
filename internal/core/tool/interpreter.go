package tool

import "os/exec"

// Interpreter describe un intérprete de scripts soportado.
// Se detecta al arrancar la app y se cachea en Deps.
type Interpreter struct {
	// Name es el binario a buscar en PATH: "bash", "python3".
	Name string
	// Display es el nombre para mostrar en la UI: "Bash", "Python 3".
	Display string
	// ScriptType es el tipo de Tool que produce este intérprete.
	ScriptType ScriptType
	// Ext es la extensión del archivo temporal: ".sh", ".py".
	Ext string
	// Path es la ruta encontrada. Vacío si no está instalado.
	Path string
}

// Available indica si el intérprete está instalado.
func (i Interpreter) Available() bool { return i.Path != "" }

// Detect busca los intérpretes soportados en PATH.
// Se llama UNA vez al arrancar la app; el resultado se cachea.
//
// Nota: Go queda pendiente hasta que exista su executor.
// Nota: python3 va antes que python — si ambos existen, ganamos
//
//	python3 (ver Available()).
func Detect() []Interpreter {
	all := []Interpreter{
		{Name: "bash", Display: "Bash", ScriptType: ScriptTypeBash, Ext: ".sh"},
		{Name: "python3", Display: "Python 3", ScriptType: ScriptTypePython, Ext: ".py"},
		{Name: "python", Display: "Python", ScriptType: ScriptTypePython, Ext: ".py"},
	}
	for i := range all {
		if p, err := exec.LookPath(all[i].Name); err == nil {
			all[i].Path = p
		}
	}
	return all
}

// Available filtra la lista y deduplica por ScriptType.
// Ej: si python3 y python están instalados, solo conserva python3
// (el primero en la lista con ese ScriptType).
func Available(list []Interpreter) []Interpreter {
	seen := make(map[ScriptType]bool, len(list))
	out := make([]Interpreter, 0, len(list))
	for _, i := range list {
		if !i.Available() {
			continue
		}
		if seen[i.ScriptType] {
			continue
		}
		seen[i.ScriptType] = true
		out = append(out, i)
	}
	return out
}

// Missing devuelve los intérpretes NO disponibles.
// Útil para mostrar warnings al arrancar.
func Missing(list []Interpreter) []Interpreter {
	var out []Interpreter
	for _, i := range list {
		if !i.Available() {
			out = append(out, i)
		}
	}
	return out
}
