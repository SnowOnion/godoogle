package collect

// Not `goimports`-ed by design. Copy from candidates.go
import (
	// basic data types & data structures
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
