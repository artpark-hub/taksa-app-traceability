package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
// This allows Kratos to automatically inject these usecases into your services. [cite: 124, 214]
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewTraceabilityUsecase)
