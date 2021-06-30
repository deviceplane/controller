package query

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/deviceplane/controller/pkg/models"
	"github.com/stretchr/testify/require"
)

type Scenario struct {
	desc  string
	in    []models.Device
	query models.Query
	out   []models.Device
}

func testScenario(t *testing.T, scenario Scenario) {
	t.Helper()

	selectedDevices, _, err := QueryDevices(QueryDependencies{}, scenario.in, scenario.query)
	require.NoError(t, err, scenario.desc)
	require.Equal(t, scenario.out, selectedDevices, scenario.desc)
}

func testEmptyScenario(t *testing.T, scenario Scenario) {
	t.Helper()

	selectedDevices, _, err := QueryDevices(QueryDependencies{}, scenario.in, scenario.query)
	require.Error(t, err, scenario.desc)
	require.Len(t, selectedDevices, 0, scenario.desc)
}

func TestQueryDevices(t *testing.T) {
	t.Run("device properties", func(t *testing.T) {
		scenarios := []Scenario{
			{
				desc: "Query online device for online status",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.DevicePropertyCondition,
							Params: map[string]interface{}{
								"property": "status",
								"operator": models.OperatorIs,
								"value":    "online",
							},
						},
					},
				},
				out: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
					},
				},
			},
			{
				desc: "Query offline device for online status",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOffline,
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.DevicePropertyCondition,
							Params: map[string]interface{}{
								"property": "status",
								"operator": models.OperatorIs,
								"value":    "online",
							},
						},
					},
				},
				out: []models.Device{},
			},
		}

		for _, scenario := range scenarios {
			testScenario(t, scenario)
		}
	})

	t.Run("label", func(t *testing.T) {
		scenarios := []Scenario{
			{
				desc: "Query labeled device for matching label+value",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorIs,
								"value":    "b",
							},
						},
					},
				},
				out: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
			},
			{
				desc: "Query labeled device for matching label's existence",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelExistenceCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorExists,
							},
						},
					},
				},
				out: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
			},
			{
				desc: "Query labeled device for missing label's non-existence",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelExistenceCondition,
							Params: map[string]interface{}{
								"key":      "c",
								"operator": models.OperatorNotExists,
							},
						},
					},
				},
				out: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
			},

			{
				desc: "Query labeled device for matching label with different value",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorIs,
								"value":    "x",
							},
						},
					},
				},
				out: []models.Device{},
			},
			{
				desc: "Query labeled device for missing label",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelExistenceCondition,
							Params: map[string]interface{}{
								"key":      "c",
								"operator": models.OperatorExists,
							},
						},
					},
				},
				out: []models.Device{},
			},
			{
				desc: "Query labeled device for missing label",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelExistenceCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorNotExists,
							},
						},
					},
				},
				out: []models.Device{},
			},
		}

		for _, scenario := range scenarios {
			testScenario(t, scenario)
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		scenarios := []Scenario{
			{
				desc: "Empty query",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOffline,
					},
				},
				query: models.Query{},
				out: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOffline,
					},
				},
			},
		}
		for _, scenario := range scenarios {
			testScenario(t, scenario)
		}
	})

	t.Run("queries that should error", func(t *testing.T) {
		scenarios := []Scenario{
			{
				desc: "LabelValueCondition with an OperatorExists",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorExists,
								"value":    "b",
							},
						},
					},
				},
			},
			{
				desc: "LabelExistenceCondition with an OperatorIs",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelExistenceCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorIs,
								"value":    "b",
							},
						},
					},
				},
			},
			{
				desc: "models.DevicePropertyCondition with an OperatorExists",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.DevicePropertyCondition,
							Params: map[string]interface{}{
								"property": "a",
								"operator": models.OperatorExists,
								"value":    "b",
							},
						},
					},
				},
			},
			{
				desc: "LabelValueCondition with 'property' instead of 'key'",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"property": "a",
								"operator": models.OperatorExists,
								"value":    "b",
							},
						},
					},
				},
			},
			{
				desc: "Empty LabelExistenceCondition",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type:   models.LabelExistenceCondition,
							Params: map[string]interface{}{},
						},
					},
				},
			},
			{
				desc: "LabelExistenceCondition without operator",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelExistenceCondition,
							Params: map[string]interface{}{
								"key": "a",
							},
						},
					},
				},
			},
			{
				desc: "LabelExistenceCondition without operator",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelExistenceCondition,
							Params: map[string]interface{}{
								"key": "a",
							},
						},
					},
				},
			},
			{
				desc: "models.DevicePropertyCondition with invalid property",
				in: []models.Device{
					{
						ID:     "one",
						Status: models.DeviceStatusOnline,
						Labels: map[string]string{
							"a": "b",
						},
					},
				},
				query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.DevicePropertyCondition,
							Params: map[string]interface{}{
								"property": "qweroiweqroijfdsfafdew",
								"operator": models.OperatorIs,
								"value":    "qweiofioweweiweofewi",
							},
						},
					},
				},
			},
		}

		for _, scenario := range scenarios {
			testEmptyScenario(t, scenario)
		}
	})
}

func TestFiltersFromQuery(t *testing.T) {
	filtersA := models.Filter{
		models.Condition{
			Type: models.DevicePropertyCondition,
			Params: map[string]interface{}{
				"property": "status",
				"operator": string(models.OperatorIs),
				"value":    "online",
			},
		},
		models.Condition{
			Type: models.DevicePropertyCondition,
			Params: map[string]interface{}{
				"property": "status",
				"operator": string(models.OperatorIs),
				"value":    "offline",
			},
		},
	}

	filtersB := models.Filter{
		models.Condition{
			Type: models.DevicePropertyCondition,
			Params: map[string]interface{}{
				"property": "status",
				"operator": string(models.OperatorIs),
				"value":    "online",
			},
		},
	}

	jsonFilterA, _ := json.Marshal(filtersA)
	encodedFilterA := base64.StdEncoding.EncodeToString(jsonFilterA)
	jsonFilterB, _ := json.Marshal(filtersB)
	encodedFilterB := base64.StdEncoding.EncodeToString(jsonFilterB)

	query := map[string][]string{
		"filter": {
			encodedFilterA,
			encodedFilterB,
		},
	}

	result, err := FiltersFromQuery(query)
	require.NoError(t, err)

	require.Len(t, result, 2)
	require.Len(t, result[0], 2)
	require.Len(t, result[1], 1)
	require.Equal(t, filtersA[0], result[0][0])
	require.Equal(t, filtersA[1], result[0][1])
	require.Equal(t, filtersB[0], result[1][0])
}
