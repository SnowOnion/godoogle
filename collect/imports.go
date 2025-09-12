package collect

// Not `goimports`-ed by design. Copy from candidates.go
import (
	// basic data types & data structures
	_ "github.com/bytedance/gg/collection"
	_ "github.com/bytedance/gg/collection/list"
	_ "github.com/bytedance/gg/collection/set"
	_ "github.com/bytedance/gg/collection/skipmap"
	_ "github.com/bytedance/gg/collection/skipset"
	_ "github.com/bytedance/gg/collection/tuple"
	_ "github.com/bytedance/gg/gcond"
	_ "github.com/bytedance/gg/gconv"
	_ "github.com/bytedance/gg/gfunc"
	_ "github.com/bytedance/gg/gmap"
	_ "github.com/bytedance/gg/goption"
	_ "github.com/bytedance/gg/gptr"
	_ "github.com/bytedance/gg/gresult"
	_ "github.com/bytedance/gg/gslice"
	_ "github.com/bytedance/gg/gson"
	_ "github.com/bytedance/gg/gstd/gsync"
	_ "github.com/bytedance/gg/gvalue"
	_ "github.com/dominikbraun/graph"
	_ "github.com/samber/lo"
	_ "github.com/samber/mo"
	_ "github.com/thoas/go-funk"
	// Date and Time
	_ "github.com/jinzhu/now"
	// concurrent & parallel
	_ "github.com/sourcegraph/conc"
	// HTTP
	_ "github.com/go-resty/resty/v2"
	// ORM
	_ "gorm.io/gorm"
	_ "xorm.io/xorm"
	// Testing
	_ "github.com/stretchr/testify/assert"
)
