package encoding

import (
	"github.com/JBrVJxsc/xk6-encoding/encoding"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/encoding", encoding.New())
}
