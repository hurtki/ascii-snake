package ws

import (
	"encoding/binary"

	"github.com/hurtki/ascii-snake/internal/srv/app"
)

// Serializes whole plot into binary format
// One cell: [[2b player id][1 byte value][1 byte isHead]]
func SerializePlot(plot [][]app.Cell) []byte {
	rows := len(plot)
	if rows == 0 {
		return nil
	}
	cols := len(plot[0])
	if cols == 0 {
		return nil
	}

	// allocating once
	res := make([]byte, rows*cols*4)
	idx := 0

	for i := range rows {
		for j := range cols {
			cell := &plot[i][j]

			// 1-2 bytes: PlayerID (uint16 LittleEndian)
			binary.LittleEndian.PutUint16(res[idx:], uint16(cell.PlayerID))

			// 3 byte: Value
			res[idx+2] = byte(cell.Value)

			// 4 byte: IsHead
			if cell.IsHead {
				res[idx+3] = 1
			} else {
				res[idx+3] = 0
			}

			idx += 4
		}
	}

	return res
}
