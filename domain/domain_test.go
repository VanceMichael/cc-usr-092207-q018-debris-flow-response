package debrisflowresponse

import (
    "encoding/json"
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
    if value.Version < 2 { t.Fatalf("共享资料版本过期: %d", value.Version) }
    if value.Situation.Found + value.Situation.StillMissing != value.Situation.MissingInitial {
        t.Fatalf("失联人员统计与名单不一致: %d + %d != %d", value.Situation.Found, value.Situation.StillMissing, value.Situation.MissingInitial)
    }
    if value.Situation.StillMissing == 0 { t.Fatal("搜救仍在持续，未找到人数不应为零") }
    if len(value.Situation.SyncRequired) == 0 { t.Fatal("必须声明需同步更新的名单、搜索区域和对外表述") }
    if !value.EventPolicy.AppendOnly { t.Fatal("事件规则必须只追加") }
    if !value.FieldSync.SignedTasks || !value.FieldSync.ReportBack { t.Fatal("弱网现场必须支持签名任务与回传") }
    if len(value.AccessRules) < 2 { t.Fatal("家属联络与新闻发布的可见性规则不完整") }
    if len(value.CommanderView) < 4 { t.Fatal("指挥员视图要素不完整") }
}

func TestStaleVersionRejected(t *testing.T) {
    raw := []byte(`{"domain":"debris-flow-response","version":1,"sample_id":"record-018","actors":["应急指挥人员","现场搜救队伍"],"facts":["局部强降雨","五人失联"],"constraints":["失联人员核验","搜索网格分派"]}`)
    _, err := Parse(raw)
    if err == nil { t.Fatal("过期版本不应通过校验") }
    if !strings.Contains(err.Error(), "版本过期") { t.Fatalf("过期版本应提示同步更新: %v", err) }
}

func TestInconsistentSituationRejected(t *testing.T) {
    raw, err := os.ReadFile("../fixtures/domain.json")
    if err != nil { t.Fatal(err) }
    value, err := Parse(raw)
    if err != nil { t.Fatal(err) }
    value.Situation.StillMissing++
    mutated, err := json.Marshal(value)
    if err != nil { t.Fatal(err) }
    if _, err := Parse(mutated); err == nil { t.Fatal("统计与名单不一致的记录不应通过校验") }
}
