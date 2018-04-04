package fb

import "testing"

func TestFormLead_EncodeJSON(t *testing.T) {
	leads := []struct {
		Lead FormLead
		JSON string
	}{
		{
			Lead: FormLead{},
			JSON: `{"created_time":"","id":"","field_data":[]}`,
		},
		{
			Lead: FormLead{
				CreatedTime: "12345",
				ID:          "2342342342",
				FieldData: []struct {
					Name   string
					Values []string
				}{
					{"a", []string{"b"}},
					{"c", []string{"d", "e"}},
				},
			},
			JSON: `{"created_time":"12345","id":"2342342342","field_data":[{"name":"a","values":["b"]},{"name":"c","values":["d","e"]}]}`,
		},
	}
	for i, lead := range leads {
		encoded := string(lead.Lead.EncodeJSON())
		if encoded != lead.JSON {
			t.Errorf("failed encoding text index %d; got %v", i, encoded)
		}
	}
}
