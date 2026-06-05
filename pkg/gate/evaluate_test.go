package gate

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

func fireRule(id string, lvl Level, fire bool) Rule {
	return Rule{
		Meta: RuleMeta{ID: id, Level: lvl, Desc: id},
		Check: func(Context) (bool, quest.Fact) {
			return fire, quest.Fact{Where: id}
		},
	}
}

func TestEvaluateLevelAggregation(t *testing.T) {
	cases := []struct {
		name  string
		rules []Rule
		want  quest.Outcome
	}{
		{"all silent", []Rule{fireRule("a", LevelFail, false), fireRule("b", LevelReview, false)}, quest.OutPass},
		{"review only", []Rule{fireRule("a", LevelFail, false), fireRule("b", LevelReview, true)}, quest.OutReview},
		{"fail wins over review", []Rule{fireRule("a", LevelFail, true), fireRule("b", LevelReview, true)}, quest.OutFail},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := Evaluate(c.rules, Context{})
			if v.Outcome != c.want {
				t.Fatalf("outcome = %s, want %s", v.Outcome, c.want)
			}
		})
	}
}

func TestEvaluateStampsRuleID(t *testing.T) {
	v := Evaluate([]Rule{fireRule("r1", LevelFail, true)}, Context{})
	if len(v.Facts) != 1 || v.Facts[0].Rule != "r1" {
		t.Fatalf("facts = %+v, want one fact stamped with r1", v.Facts)
	}
}

func TestCatalog(t *testing.T) {
	metas := Catalog([]Rule{fireRule("a", LevelFail, false), fireRule("b", LevelReview, false)})
	if len(metas) != 2 || metas[0].ID != "a" || metas[1].ID != "b" {
		t.Fatalf("catalog = %+v", metas)
	}
}
