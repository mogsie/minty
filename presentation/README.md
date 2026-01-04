# Presentation Layer Status

The presentation layer packages (`presentation/`) currently have API mismatches with the core library and do not compile. These packages provide domain-specific UI components that bridge the domain libraries with theme implementations.

## Issues to Address

1. **If/IfElse API mismatch**: Code uses `If(condition, Node, Node)` but the actual signature is `If(condition, H)`. Need to use `IfElse` for two-branch conditionals.

2. **Each return type**: `mintyex.Each` now returns `[]Node` instead of `[]H`. Code using `Each` with `GridLayout` needs updating.

3. **Node vs H confusion**: Many places pass `Node` values where `H` (template functions) are expected. Need to wrap with `miex.WrapNode()` or restructure the code.

## Affected Packages

- `presentation/mintycartui/` - E-commerce UI components
- `presentation/mintyfinui/` - Finance UI components
- `presentation/mintymoveui/` - Logistics UI components

## Recommended Fix Approach

1. Replace `miex.If(cond, node1, node2)` with `miex.IfElse(cond, func1, func2)`
2. Wrap Node values in H functions: `func(b *mi.Builder) mi.Node { return existingNode }`
3. Update `GridLayout` calls to work with `[]Node` instead of `[]H`
4. Remove unused sample data declarations

These packages are included for reference but should be considered **work in progress**.
