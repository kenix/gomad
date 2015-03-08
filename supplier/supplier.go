package supplier

import (
	"github.com/kenix/gomad/bytebuffer"
)

type Supplier interface {
	Get(bytebuffer.ByteBuffer)
}
