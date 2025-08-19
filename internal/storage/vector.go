package storage

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
)

// Vector represents a PostgreSQL vector type
type Vector []float32

// Value implements the driver.Valuer interface.
func (v Vector) Value() (driver.Value, error) {
	if len(v) == 0 {
		return nil, nil
	}

	// Calculate the total size: 2 bytes for dimension, then 4 bytes per float32
	size := 2 + len(v)*4
	buf := make([]byte, size)

	// Write dimension (uint16, big endian)
	binary.BigEndian.PutUint16(buf[0:2], uint16(len(v)))

	// Write float32 values
	for i, f := range v {
		binary.BigEndian.PutUint32(buf[2+i*4:], math.Float32bits(f))
	}

	return buf, nil
}

// Scan implements the sql.Scanner interface.
func (v *Vector) Scan(src interface{}) error {
	if src == nil {
		*v = nil
		return nil
	}

	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("pgvector: expected []byte, got %T", src)
	}

	if len(b) < 2 {
		return fmt.Errorf("pgvector: invalid vector format")
	}

	// Read dimension
	dim := binary.BigEndian.Uint16(b[0:2])
	if len(b) != 2+int(dim)*4 {
		return fmt.Errorf("pgvector: expected %d bytes, got %d", 2+dim*4, len(b))
	}

	// Read float32 values
	*v = make(Vector, dim)
	for i := uint16(0); i < dim; i++ {
		bits := binary.BigEndian.Uint32(b[2+i*4:])
		(*v)[i] = math.Float32frombits(bits)
	}

	return nil
}
