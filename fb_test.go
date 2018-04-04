package fb

import "testing"

func TestLeadGenEntry_MarshalJSON(t *testing.T) {
	leads := []struct {
		Lead LeadGenEntry
		JSON string
	}{
		{
			Lead: LeadGenEntry{},
			JSON: `{"ad_id":"","form_id":"","leadgen_id":"","page_id":"","adgroup_id":"","created_time":0}`,
		},
		{
			Lead: LeadGenEntry{
				AdID:        "123",
				FormID:      "23425",
				LeadgenID:   "432532509",
				PageID:      "32097042",
				AdgroupID:   "9253195",
				CreatedTime: 1522862162,
			},
			JSON: `{"ad_id":"123","form_id":"23425","leadgen_id":"432532509","page_id":"32097042","adgroup_id":"9253195","created_time":1522862162}`,
		},
	}
	for i, lead := range leads {
		b, err := lead.Lead.MarshalJSON()
		if err != nil {
			// There must never be an error.
			t.Fatalf("got an error encoding: %v", err)
		}
		encoded := string(b)
		if encoded != lead.JSON {
			t.Errorf("failed encoding text index %d; got %v", i, encoded)
		}
	}
}
