package providers_test

import (
	"testing"

	"github.com/lixuanqun/voice-realtime/internal/providers"

	_ "github.com/lixuanqun/voice-realtime/internal/providers/bailian"
	_ "github.com/lixuanqun/voice-realtime/internal/providers/stepfun"
	_ "github.com/lixuanqun/voice-realtime/internal/providers/volcengine"
	_ "github.com/lixuanqun/voice-realtime/internal/providers/zhipu"
)

func TestRegistryListsProviders(t *testing.T) {
	names := providers.List()
	if len(names) != 4 {
		t.Fatalf("expected 4 providers, got %v", names)
	}
	for _, want := range []string{"bailian", "stepfun", "volcengine", "zhipu"} {
		if _, err := providers.Get(want); err != nil {
			t.Fatalf("missing provider %s: %v", want, err)
		}
	}
}
