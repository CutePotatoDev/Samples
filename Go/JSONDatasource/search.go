package main

import (
	"github.com/labstack/echo"
	"net/http"
	str "strings"
)

type Search struct {
	Target string `json:"target"`
}

func (h *Handler) search(c echo.Context) error {
	search := new(Search)

	if err := c.Bind(search); err != nil {
		return err
	}

	queryparts := str.Split(search.Target, "/")

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
			data = append(data, queryparts[0]+"/"+e.Host)
		}

		return c.JSON(http.StatusOK, data)
	} else if cap(queryparts) < 4 {

		hosts, err := h.GetZabbixHosts(queryparts[1], groups[0].GroupId)
		CheckErr(err, c)

		items, err := h.GetZabbixItems(queryparts[2], hosts[0].HostId)
		CheckErr(err, c)

		for _, e := range items {
			data = append(data, queryparts[0]+"/"+queryparts[1]+"/"+e.Name)
		}

		return c.JSON(http.StatusOK, data)
	}
	return c.JSON(http.StatusOK, data)
}
