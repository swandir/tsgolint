package jsx_no_leaked_render

import (
	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/microsoft/typescript-go/shim/checker"
	"github.com/microsoft/typescript-go/shim/jsnum"
	"github.com/typescript-eslint/tsgolint/internal/rule"
	"github.com/typescript-eslint/tsgolint/internal/utils"
)

func buildDefaultMessage() rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "default",
		Description: "Potential 0 or NaN render leak.",
	}
}

var JsxNoLeakedRenderRule = rule.Rule{
	Name: "jsx-no-leaked-render",
	Run: func(ctx rule.RuleContext, options any) rule.RuleListeners {
		var inspect func(node *ast.Node, visited map[*ast.Symbol]bool) *ast.Node
		inspect = func(node *ast.Node, visited map[*ast.Symbol]bool) *ast.Node {
			if node == nil {
				return nil
			}
			node = ast.SkipOuterExpressions(node, ast.OEKAll)
			switch node.Kind {
			case ast.KindJsxElement, ast.KindJsxFragment, ast.KindJsxSelfClosingElement:
				return nil
			case ast.KindBinaryExpression:
				bin := node.AsBinaryExpression()
				if bin.OperatorToken.Kind != ast.KindAmpersandAmpersandToken {
					return nil
				}
				left := ast.SkipOuterExpressions(bin.Left, ast.OEKAll)
				if left.Kind == ast.KindPrefixUnaryExpression &&
					left.AsPrefixUnaryExpression().Operator == ast.KindExclamationToken {
					return inspect(bin.Right, visited)
				}
				if allPartsAllowed(utils.GetConstrainedTypeAtLocation(ctx.TypeChecker, bin.Left)) {
					return inspect(bin.Right, visited)
				}
				return bin.Left
			case ast.KindConditionalExpression:
				cond := node.AsConditionalExpression()
				if r := inspect(cond.WhenTrue, visited); r != nil {
					return r
				}
				return inspect(cond.WhenFalse, visited)
			case ast.KindIdentifier:
				symbol := ctx.TypeChecker.GetSymbolAtLocation(node)
				if symbol == nil || visited[symbol] {
					return nil
				}
				visited[symbol] = true
				for _, decl := range symbol.Declarations {
					if decl.Kind == ast.KindVariableDeclaration {
						if init := decl.AsVariableDeclaration().Initializer; init != nil {
							return inspect(init, visited)
						}
					}
				}
			}
			return nil
		}

		return rule.RuleListeners{
			ast.KindJsxExpression: func(node *ast.Node) {
				expr := node.AsJsxExpression().Expression
				if expr == nil {
					return
				}
				if culprit := inspect(expr, map[*ast.Symbol]bool{}); culprit != nil {
					ctx.ReportNode(culprit, buildDefaultMessage())
				}
			},
		}
	},
}

func allPartsAllowed(t *checker.Type) bool {
	for _, part := range utils.UnionTypeParts(t) {
		if !isPartAllowed(part) {
			return false
		}
	}
	return true
}

// Number and bigint operands can be falsy (0, -0, NaN) and would render
// the value into the DOM. Only literal types whose value is known to be
// truthy are safe.
func isPartAllowed(t *checker.Type) bool {
	flags := checker.Type_flags(t)

	// An intersection like `number & { __brand }` is a number at runtime —
	// the brand member is type-only. Require every member to be safe.
	if flags&checker.TypeFlagsIntersection != 0 {
		for _, part := range t.Types() {
			if !isPartAllowed(part) {
				return false
			}
		}
		return true
	}

	if flags&checker.TypeFlagsNumberLiteral != 0 {
		val, _ := t.AsLiteralType().Value().(jsnum.Number)
		return val != 0 && !val.IsNaN()
	}
	if flags&checker.TypeFlagsBigIntLiteral != 0 {
		val, _ := t.AsLiteralType().Value().(jsnum.PseudoBigInt)
		return val.Sign() != 0
	}
	if flags&(checker.TypeFlagsNumberLike|checker.TypeFlagsBigIntLike) != 0 {
		return false
	}

	return flags&(checker.TypeFlagsAny|
		checker.TypeFlagsNever|
		checker.TypeFlagsNull|
		checker.TypeFlagsUndefined|
		checker.TypeFlagsVoid|
		checker.TypeFlagsBooleanLike|
		checker.TypeFlagsStringLike|
		checker.TypeFlagsEnumLike|
		checker.TypeFlagsObject|
		checker.TypeFlagsNonPrimitive|
		checker.TypeFlagsESSymbolLike) != 0
}
