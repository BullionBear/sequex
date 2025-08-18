package init

import (
	_ "github.com/BullionBear/sequex/internal/nodeimpl/app/bar"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/app/trade"   // Import to register Trade node
	_ "github.com/BullionBear/sequex/internal/nodeimpl/example/rng" // Import to register RNG node
	_ "github.com/BullionBear/sequex/internal/nodeimpl/example/sum" // Import to register Sum node
)
