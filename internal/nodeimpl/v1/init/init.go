package init

import (
	_ "github.com/BullionBear/sequex/internal/nodeimpl/v1/app/bar"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/v1/app/trade"   // Import to register Trade node
	_ "github.com/BullionBear/sequex/internal/nodeimpl/v1/example/rng" // Import to register RNG node
	_ "github.com/BullionBear/sequex/internal/nodeimpl/v1/example/sum" // Import to register Sum node
)
