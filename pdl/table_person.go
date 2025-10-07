package pdl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tablePerson() *plugin.Table {
	return &plugin.Table{
		Name:        "pdl_person",
		Description: "Retrieve enriched person profiles using the People Data Labs API.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"email", "phone", "profile", "name"}),
			Hydrate:    listPdlPerson,
		},
		Columns: []*plugin.Column{
			{Name: "status", Type: proto.ColumnType_INT, Description: "Status value returned by the API.", Transform: transform.FromField("status")},
			{Name: "likelihood", Type: proto.ColumnType_DOUBLE, Description: "Match likelihood score.", Transform: transform.FromField("likelihood")},
			{Name: "warnings", Type: proto.ColumnType_JSON, Description: "Any warnings returned by the API.", Transform: transform.FromField("warnings")},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Enriched person data payload.", Transform: transform.FromField("data")},
			{Name: "raw_response", Type: proto.ColumnType_JSON, Description: "Full API response payload.", Transform: transform.FromValue()},
		},
	}
}

func listPdlPerson(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := getConfig(d.Connection)
	if err != nil {
		return nil, err
	}

	apiKey, err := conn.apiKey(ctx, d)
	if err != nil {
		return nil, err
	}

	if apiKey == "" {
		return nil, fmt.Errorf("api_key must be configured or the PDL_API_KEY environment variable must be set")
	}

	payload := map[string]interface{}{}

	quals := d.EqualsQuals
	if quals == nil {
		quals = map[string]*plugin.QualValue{}
	}

	if q := quals["email"]; q != nil {
		payload["email"] = q.GetStringValue()
	}
	if q := quals["phone"]; q != nil {
		payload["phone"] = q.GetStringValue()
	}
	if q := quals["profile"]; q != nil {
		payload["profile"] = q.GetStringValue()
	}
	if q := quals["name"]; q != nil {
		payload["name"] = q.GetStringValue()
	}

	if len(payload) == 0 {
		return nil, fmt.Errorf("at least one key column (email, phone, profile, name) must be provided")
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.peopledatalabs.com/v5/person/enrich", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received error response %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, result)

	return nil, nil
}

func getConfig(conn *plugin.Connection) (*Config, error) {
	if conn == nil || conn.Config == nil {
		return &Config{}, nil
	}

	cfg, ok := conn.Config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid connection config type")
	}

	return cfg, nil
}
