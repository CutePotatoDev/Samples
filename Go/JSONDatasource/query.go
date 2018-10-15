package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlekSi/zabbix"
	"github.com/iancoleman/orderedmap"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	str "strings"
	"time"
)

type Target struct {
	Target string `json:"target"`
	RefId  string `json:"refId"`
	Type   string `json:"type"`
}

type Range struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Query struct {
	PanelId       int      `json:"panelId"`
	Range         Range    `json:"range,omitempty"`
	Targets       []Target `json:"targets,omitempty"`
	MaxDataPoints int      `json:"maxDataPoints"`
}

type TableResp struct {
	Type    string          `json:"type"`
	Columns []Column        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}

type Column struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

func isJSON(str string) (bool, *orderedmap.OrderedMap) {
	omap := orderedmap.New()
	err := json.Unmarshal([]byte(str), &omap)
	return err == nil, omap
}

func inCollumns(val Column, array []Column) bool {
	for _, value := range array {
		if value == val {
			return true
		}
	}
	return false
}

func HandlerDefault(target Target, result interface{}) interface{} {
	datain := result.([]interface{})

	data := [][]interface{}{}
	columns := []Column{}
	// notjson := false

	columns = append(columns, Column{"Time", "time"})

	for _, item := range datain {
		item := item.(map[string]interface{})

		nano, _ := strconv.Atoi(item["ns"].(string))
		timestamp, _ := strconv.Atoi(item["clock"].(string) + strconv.Itoa(nano/1e6))

		dat := []interface{}{timestamp}

		if ok, values := isJSON(item["value"].(string)); ok {

			for _, key := range values.Keys() {
				value, _ := values.Get(key)

				col := Column{key, "string"}
				if !inCollumns(col, columns) {
					columns = append(columns, col)
				}

				dat = append(dat, value)
			}

			data = append(data, dat)

		} else {
			// notjson = true
			// fmt.Println(item["value"])
			// data = append(data, []interface{}{item["value"]})
			fmt.Println(item)
		}
	}

	// if notjson {
	// 	columns = append(columns, Column{"Value", "string"})
	// }

	resp := TableResp{
		target.Type,
		columns,
		data,
	}

	return []TableResp{resp}
}

func HandlerSuspiciousAgents(target Target, result interface{}) interface{} {
	datain := result.([]interface{})

	data := [][]interface{}{}
	columns := []Column{}

	columns = append(columns, Column{"Time", "time"})
	columns = append(columns, Column{"AppServer", "string"})

	for _, item := range datain {
		item := item.(map[string]interface{})

		nano, _ := strconv.Atoi(item["ns"].(string))
		timestamp, _ := strconv.Atoi(item["clock"].(string) + strconv.Itoa(nano/1e6))

		dat := []interface{}{timestamp, str.Replace(str.Split(target.Target, " ")[1], ".xxx.xxx.xxx.net", "", 1)} //Data.
		data = append(data, dat)

		scanner := bufio.NewScanner(str.NewReader(item["value"].(string)))
		i := 0
		for scanner.Scan() {
			i++

			if i < 3 {
				continue
			}

			elements := str.Fields(scanner.Text())
			for i, element := range elements {
				dat = append(dat, element)

				col := Column{strconv.Itoa(i), "string"}
				if !inCollumns(col, columns) {
					columns = append(columns, col)
				}
			}

			data = append(data, dat)
		}
	}

	resp := TableResp{
		target.Type,
		columns,
		data,
	}

	return []TableResp{resp}
}

var HandlerMap = map[string]func(Target, interface{}) interface{}{
	"HandlerDefault":          HandlerDefault,
	"HandlerSuspiciousAgents": HandlerSuspiciousAgents,
}

func (h *Handler) query(c echo.Context) error {
	query := new(Query)

	if err := c.Bind(query); err != nil {
		return err
	}

	var resp []interface{}

	for _, target := range query.Targets {
		queryparts := str.Split(target.Target, "/")

		groups, err := h.GetZabbixGroups(queryparts[0])
		CheckErr(err, c)

		hosts, err := h.GetZabbixHosts(queryparts[1], groups[0].GroupId)
		CheckErr(err, c)

		items, err := h.GetZabbixItems(queryparts[2], hosts[0].HostId)
		CheckErr(err, c)

		timefrom, _ := time.Parse("2006-01-02T15:04:05.000Z", query.Range.From)
		timeto, _ := time.Parse("2006-01-02T15:04:05.000Z", query.Range.To)

		conf := zabbix.Params{
			"itemids":   items[0].ItemId,
			"history":   items[0].ValueType,
			"sortfield": "clock",
			"sortorder": "DESC",
			"limit":     query.MaxDataPoints,
			"time_from": timefrom.Unix(),
			"time_till": timeto.Unix(),
		}

		hist, _ := h.zabbix_api.Call("history.get", conf)
		// fmt.Println(target.Target)

		resp = append(resp, HandlerMap[GetDataHandler(target.Target)](target, hist.Result).([]TableResp)[0])
	}
	return c.JSON(http.StatusOK, resp)
}
