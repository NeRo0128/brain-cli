package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// rowScanner abstrae *sql.Row y *sql.Rows (ambos tienen Scan).
// Permite que los helpers de scan funcionen con cualquiera.
type rowScanner interface {
	Scan(dest ...any) error
}

// boolToInt convierte bool a INTEGER (0/1) para SQLite.
// SQLite no tiene BOOLEAN nativo; guardamos como INTEGER.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// nullTimeToPtr convierte sql.NullTime a *time.Time.
// nil si la columna era NULL.
func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// isUniqueViolation detecta violaciones de UNIQUE en SQLite.
// Usamos string matching porque modernc no expone un código
// portable de forma estable. Es lo que hacen proyectos como
// golang-migrate y se considera aceptable.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

// nullInt64ToPtr convierte sql.NullInt64 a *int (nil si NULL).
func nullInt64ToPtr(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	i := int(n.Int64)
	return &i
}

// ptrIntToArg convierte *int a any para pasarlo como argumento SQL.
func ptrIntToArg(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}

// marshalJSONMap serializa map → JSON. Vacío → "{}".
func marshalJSONMap(m map[string]string) (string, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}
	return string(b), nil
}

// unmarshalJSONMap parsea JSON → map. "{}" o "" → nil.
func unmarshalJSONMap(s string) (map[string]string, error) {
	if s == "" || s == "{}" {
		return nil, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	return m, nil
}

// placeholders genera "?,?,?" para N argumentos.
func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}
