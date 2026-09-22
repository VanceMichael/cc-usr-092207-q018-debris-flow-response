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
    AsOf string `json:"as_of"`
    Actors []string `json:"actors"`
    Facts []string `json:"facts"`
    SituationElements []string `json:"situation_elements"`
    EventLogTypes []string `json:"event_log_types"`
    Constraints []string `json:"constraints"`
}

// Parse 读取并检查带版本的业务资料。
func Parse(raw []byte) (Record, error) {
    var value Record
    if err := json.Unmarshal(raw, &value); err != nil { return Record{}, err }
    if value.Domain == "" || value.Version < 1 || value.SampleID == "" || value.AsOf == "" ||
        len(value.Actors) < 2 || len(value.Facts) < 2 ||
        len(value.SituationElements) < 2 || len(value.EventLogTypes) < 2 ||
        len(value.Constraints) < 2 {
        return Record{}, errors.New("共享资料缺少必要字段")
    }
    return value, nil
}
