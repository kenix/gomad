package supplier

import (
	"github.com/mentopolis/gomad/bytebuffer"
)

type Supplier interface {
	Get(bytebuffer.ByteBuffer)
}
