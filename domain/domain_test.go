package debrisflowresponse

import (
    "os"
    "strings"
    "testing"
)

func TestFixtureMatchesDomain(t *testing.T) {
    raw, err := os.ReadFile("../fixtures/domain.json")
    if err != nil { t.Fatal(err) }
    value, err := Parse(raw)
    if err != nil { t.Fatal(err) }
    if value.Domain != "debris-flow-response" { t.Fatalf("领域标识不一致: %s", value.Domain) }
    if value.Version < 2 { t.Fatalf("十九日深夜后的资料应为 v2 或更高，实际: %d", value.Version) }
}

func TestSituationElementsAndEvents(t *testing.T) {
    raw, err := os.ReadFile("../fixtures/domain.json")
    if err != nil { t.Fatal(err) }
    value, err := Parse(raw)
    if err != nil { t.Fatal(err) }

    requiredElements := []string{
        "灾害边界", "气象地质警报", "家属核验", "搜索网格",
        "队伍能力", "车辆装备", "进出记录", "人员交接",
    }
    for _, name := range requiredElements {
        if !containsPrefix(value.SituationElements, name) {
            t.Errorf("实时态势缺少要素: %s", name)
        }
    }

    requiredEvents := []string{"重复报警", "坐标修正", "再降雨停工", "网格重搜", "轮班"}
    for _, name := range requiredEvents {
        if !containsPrefix(value.EventLogTypes, name) {
            t.Errorf("只追加事件缺少类型: %s", name)
        }
    }
}

func TestRejectsMissingV2Fields(t *testing.T) {
    raw := []byte(`{
        "domain": "debris-flow-response",
        "version": 2,
        "sample_id": "record-018",
        "actors": ["应急指挥人员", "现场搜救队伍"],
        "facts": ["事实一", "事实二"],
        "constraints": ["约束一", "约束二"]
    }`)
    if _, err := Parse(raw); err == nil {
        t.Fatal("缺少 as_of、态势要素和事件类型的资料应被拒绝")
    }
}

func containsPrefix(items []string, prefix string) bool {
    for _, item := range items {
        if strings.HasPrefix(item, prefix) { return true }
    }
    return false
}
