package utils

import (
	"encoding/json"
	"testing"
)

func TestIDResponse_Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		want    IDResponse
		wantErr bool
	}{
		{
			name: "valid JSON bytes",
			src:  []byte(`{"id": 123}`),
			want: IDResponse{ID: 123},
		},
		{
			name: "nil source",
			src:  nil,
			want: IDResponse{},
		},
		{
			name:    "invalid type",
			src:     123,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			src:     []byte(`{invalid}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r IDResponse
			err := r.Scan(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("IDResponse.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && r.ID != tt.want.ID {
				t.Errorf("IDResponse.Scan() = %v, want %v", r, tt.want)
			}
		})
	}
}

func TestSuccessResponseMap_JSON(t *testing.T) {
	resp := SuccessResponseMap[string]{
		RequestID: "req-123",
		Status:    200,
		Data:      "test data",
		TS:        "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal SuccessResponseMap: %v", err)
	}

	var decoded SuccessResponseMap[string]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal SuccessResponseMap: %v", err)
	}

	if decoded.RequestID != resp.RequestID {
		t.Errorf("RequestID = %v, want %v", decoded.RequestID, resp.RequestID)
	}
	if decoded.Status != resp.Status {
		t.Errorf("Status = %v, want %v", decoded.Status, resp.Status)
	}
	if decoded.Data != resp.Data {
		t.Errorf("Data = %v, want %v", decoded.Data, resp.Data)
	}
}

func TestErrorResponseMap_JSON(t *testing.T) {
	resp := ErrorResponseMap{
		RequestID: "req-456",
		Status:    400,
		Error:     "bad request",
		TS:        "2024-01-01T00:00:00Z",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorResponseMap: %v", err)
	}

	var decoded ErrorResponseMap
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ErrorResponseMap: %v", err)
	}

	if decoded.Error != resp.Error {
		t.Errorf("Error = %v, want %v", decoded.Error, resp.Error)
	}
}

func TestPaginationResponse_JSON(t *testing.T) {
	next := 2
	resp := PaginationResponse[[]string]{
		Items: []string{"a", "b", "c"},
		Total: 100,
		Meta: PaginationResponseMeta{
			Next:        &next,
			Prev:        nil,
			HasNextPage: true,
			HasPrevPage: false,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal PaginationResponse: %v", err)
	}

	var decoded PaginationResponse[[]string]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal PaginationResponse: %v", err)
	}

	if decoded.Total != resp.Total {
		t.Errorf("Total = %v, want %v", decoded.Total, resp.Total)
	}
	if decoded.Meta.HasNextPage != resp.Meta.HasNextPage {
		t.Errorf("HasNextPage = %v, want %v", decoded.Meta.HasNextPage, resp.Meta.HasNextPage)
	}
	if decoded.Meta.Next == nil || *decoded.Meta.Next != next {
		t.Errorf("Next = %v, want %v", decoded.Meta.Next, next)
	}
}

func TestDatabaseResponse_JSON(t *testing.T) {
	resp := DatabaseResponse[string]{
		Data:          "test",
		StatusCode:    200,
		StatusMessage: "OK",
		TS:            "2024-01-01",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal DatabaseResponse: %v", err)
	}

	var decoded DatabaseResponse[string]
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal DatabaseResponse: %v", err)
	}

	if decoded.Data != resp.Data {
		t.Errorf("Data = %v, want %v", decoded.Data, resp.Data)
	}
	if decoded.StatusCode != resp.StatusCode {
		t.Errorf("StatusCode = %v, want %v", decoded.StatusCode, resp.StatusCode)
	}
}
