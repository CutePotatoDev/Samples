package main

import (
	"fmt"
	"github.com/AlekSi/zabbix"
	"github.com/labstack/echo"
	"net/http"
	str "strings"
)

type Handler struct {
	zabbix_api *zabbix.API
}

func main() {

	api := zabbix.NewAPI("https://zabbix/api_jsonrpc.php")
	_, err := api.Login("", "")
	if err != nil {
		panic(err)
	}

	ec := echo.New()

	// ec.Use(middleware.Logger())
	// ec.Use(middleware.Recover())

	h := &Handler{zabbix_api: api}

	ec.GET("/", hey)
	ec.GET("/zabbix", h.zabbix)
	ec.POST("/search", h.search)
	ec.POST("/query", h.query)
	ec.Start(":9090")

	// for _, s := range conf.Handlers {
	// 	fmt.Printf("%s/n", s.Message)
	// }

}

func hey(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{"Working."})
}

func (h *Handler) zabbix(c echo.Context) error {
	query := c.QueryParam("query")

	queryparts := str.Split(query, ".")

	fmt.Println(cap(queryparts))
	fmt.Println(len(queryparts[0]))

	groups, err := h.GetZabbixGroups(queryparts[0])
	CheckErr(err, c)

	var data []string

	if cap(queryparts) < 2 {
		for _, e := range groups {
			data = append(data, e.Name)
		}

		return c.JSON(http.StatusOK, data)
	} else if cap(queryparts) < 3 {

		hosts, err := h.GetZabbixHosts(queryparts[1], groups[0].GroupId)
		CheckErr(err, c)

		for _, e := range hosts {
			data = append(data, queryparts[0]+"."+e.Host)
		}

		return c.JSON(http.StatusOK, data)
	} else if cap(queryparts) < 4 {

		hosts, err := h.GetZabbixHosts(queryparts[1], groups[0].GroupId)
		CheckErr(err, c)

		items, err := h.GetZabbixItems(queryparts[2], hosts[0].HostId)
		CheckErr(err, c)

		for _, e := range items {
			data = append(data, queryparts[0]+"."+queryparts[1]+"."+e.Name)
		}

		return c.JSON(http.StatusOK, data)
	} else if cap(queryparts) < 5 {
		hosts, err := h.GetZabbixHosts(queryparts[1], groups[0].GroupId)
		CheckErr(err, c)

		items, err := h.GetZabbixItems(queryparts[2], hosts[0].HostId)
		CheckErr(err, c)

		fmt.Println(items)

		conf := zabbix.Params{
			"itemids":   items[0].ItemId,
			"history":   items[0].ValueType,
			"sortfield": "clock",
			"sortorder": "DESC",
			"limit":     100,
		}

		hist, _ := h.zabbix_api.Call("history.get", conf)

		return c.JSON(http.StatusOK, hist)
	}

	return c.JSON(http.StatusOK, echo.Map{"Status": query})
}

func CheckErr(err error, c echo.Context) {
	if err != nil {
		c.JSON(http.StatusOK, echo.Map{"Error": err.Error()})
	}
}

func (h *Handler) GetZabbixGroups(groupname string) (zabbix.HostGroups, error) {
	conf := zabbix.Params{
		"output": []string{"name"},
		"search": map[string]string{
			"name": groupname,
		},
	}

	groups, err := h.zabbix_api.HostGroupsGet(conf)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (h *Handler) GetZabbixHosts(hostname string, group string) (zabbix.Hosts, error) {
	conf := zabbix.Params{
		"output":   []string{"host"},
		"groupids": group,
		"search": map[string]string{
			"name": hostname,
		},
	}

	hosts, err := h.zabbix_api.HostsGet(conf)
	if err != nil {
		return nil, err
	}

	return hosts, nil
}

func (h *Handler) GetZabbixItems(itemname string, hostid string) (zabbix.Items, error) {
	conf := zabbix.Params{
		"output":      []string{"name", "value_type"},
		"hostids":     hostid,
		"startSearch": true,
		"search": map[string]string{
			"name": itemname,
		},
	}

	items, err := h.zabbix_api.ItemsGet(conf)
	if err != nil {
		return nil, err
	}

	return items, nil
}
