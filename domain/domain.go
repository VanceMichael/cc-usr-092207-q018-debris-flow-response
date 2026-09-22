package debrisflowresponse

import (
    "encoding/json"
    "errors"
)

// Record 表示项目共享的领域资料。
type Record struct {
    Domain string `json:"domain"`
    Version int `json:"version"`
    SampleID string `json:"sample_id"`
    Actors []string `json:"actors"`
    Facts []string `json:"facts"`
    Constraints []string `json:"constraints"`
    Situation Situation `json:"situation"`
    SituationElements []string `json:"situation_elements"`
    EventPolicy EventPolicy `json:"event_policy"`
    FieldSync FieldSync `json:"field_sync"`
    AccessRules []string `json:"access_rules"`
    CommanderView []string `json:"commander_view"`
}

// Situation 表示最新态势快照，统计数字必须与失联人员名单保持一致。
type Situation struct {
    AsOf string `json:"as_of"`
    HazardArea string `json:"hazard_area"`
    MissingInitial int `json:"missing_initial"`
    Found int `json:"found"`
    StillMissing int `json:"still_missing"`
    SyncRequired []string `json:"sync_required"`
}

// EventPolicy 表示事件记录规则，现场事件只按顺序追加，不改写历史。
type EventPolicy struct {
    AppendOnly bool `json:"append_only"`
    EventTypes []string `json:"event_types"`
}

// FieldSync 表示弱网现场协同能力。
type FieldSync struct {
    WeakNetwork bool `json:"weak_network"`
    SignedTasks bool `json:"signed_tasks"`
    ReportBack bool `json:"report_back"`
}

// Parse 读取并检查带版本的业务资料，过期版本会被拒绝。
func Parse(raw []byte) (Record, error) {
    var value Record
    if err := json.Unmarshal(raw, &value); err != nil { return Record{}, err }
    if value.Domain == "" || value.Version < 2 || value.SampleID == "" || len(value.Actors) < 2 || len(value.Facts) < 2 || len(value.Constraints) < 2 {
        return Record{}, errors.New("共享资料缺少必要字段或版本过期")
    }
    if value.Situation.AsOf == "" || value.Situation.HazardArea == "" || len(value.Situation.SyncRequired) == 0 {
        return Record{}, errors.New("态势快照缺少必要字段")
    }
    if value.Situation.Found + value.Situation.StillMissing != value.Situation.MissingInitial {
        return Record{}, errors.New("失联人员统计与名单不一致")
    }
    if len(value.SituationElements) == 0 || !value.EventPolicy.AppendOnly || len(value.EventPolicy.EventTypes) == 0 {
        return Record{}, errors.New("态势要素或事件规则不完整")
    }
    if !value.FieldSync.WeakNetwork || !value.FieldSync.SignedTasks || !value.FieldSync.ReportBack {
        return Record{}, errors.New("弱网协同能力不完整")
    }
    if len(value.AccessRules) == 0 || len(value.CommanderView) == 0 {
        return Record{}, errors.New("岗位可见性或指挥员视图不完整")
    }
    return value, nil
}
