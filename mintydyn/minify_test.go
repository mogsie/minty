package mintydyn

import (
	"strings"
	"testing"
)

func TestMinifyJS(t *testing.T) {
	// Sample JavaScript similar to what we generate
	input := `
// Dynamic Component: test
class DynamicComponent_test {
    constructor() {
        this.id = 'test';
        this.container = document.getElementById(this.id);
        this.config = this.loadConfig();
        this.managers = {};
        this.externals = {};
        this.hooks = this.config.hooks || {};
        this.state = {
            currentState: null,
            filters: {},
            dependencies: new Map(),
            initialized: false
        };
        
        this.initWithDependencies();
    }
    
    loadConfig() {
        const configScript = document.getElementById(this.id + '-config');
        return configScript ? JSON.parse(configScript.textContent) : {};
    }
    
    async initWithDependencies() {
        try {
            await this.loadExternalScripts();
            this.init();
        } catch (error) {
            console.error('Init failed:', error);
        }
    }
    
    handleClick(event) {
        const action = event.target.dataset.clientAction;
        if (action) {
            this.executeAction(action, event.target, event);
        }
    }
    
    async switchToState(stateId) {
        const prevState = this.state.currentState;
        if (this.managers.states) {
            this.managers.states.switchTo(stateId);
        }
        return true;
    }
}

// Auto-initialization
document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('test')) {
        window.DynComponent_test = new DynamicComponent_test();
    }
});
`

	minified := MinifyJS(input)
	
	t.Logf("Original size: %d bytes", len(input))
	t.Logf("Minified size: %d bytes", len(minified))
	t.Logf("Reduction: %.1f%%", 100.0*(1.0-float64(len(minified))/float64(len(input))))
	
	// Check keywords are intact
	keywords := []string{"constructor", "function", "return", "const", "class", "async", "await", "switch", "document", "getElementById", "addEventListener"}
	for _, kw := range keywords {
		if !strings.Contains(minified, kw) {
			t.Errorf("Keyword '%s' missing from minified output", kw)
		}
	}
	
	// Check no broken keywords
	broken := []string{"const ructor", "func tion", "re turn", "docu ment", "get ElementById", "add EventListener"}
	for _, b := range broken {
		if strings.Contains(minified, b) {
			t.Errorf("Broken keyword found: '%s'", b)
		}
	}
	
	// Output sample
	if len(minified) > 800 {
		t.Logf("Sample output:\n%s...", minified[:800])
	} else {
		t.Logf("Full output:\n%s", minified)
	}
}
