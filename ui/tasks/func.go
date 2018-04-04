package tasks

import (
	"bytes"
	"fmt"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
)

// GetDomainRow return html blocks to display a domain info
func GetDomainRow(host string, count interface{}) string {
	var buffer bytes.Buffer
	link := fmt.Sprintf("<a href=\"?host=%v\">%v(%v)</a>", host, host, count)
	writeTag(&buffer, "li", link)
	return buffer.String()
}

// GetTaskRow return html blocks to display a task info
func GetTaskRow(task model.Task) string {
	var buffer bytes.Buffer
	buffer.WriteString("<tr>")

	linkUrl := fmt.Sprintf("/admin/task/view/%s", task.ID)
	title := fmt.Sprintf("<a title='%s' href='%s'>%s</a>", task.Url, linkUrl, util.SubStringWithSuffix(task.Url, 50, "..."))

	writeTag(&buffer, "td", title)

	if !task.SnapshotCreated.IsZero() {
		date1 := util.FormatTimeWithLocalTZ(task.SnapshotCreated)
		buffer.WriteString("<td class='timeago' title='" + date1 + "' >" + date1 + "</td>")
	} else {
		buffer.WriteString("<td >N/A</td>")
	}

	if !task.NextCheck.IsZero() {
		date2 := util.FormatTimeWithLocalTZ(task.NextCheck)
		buffer.WriteString("<td class='timeago' title='" + date2 + "' >" + date2 + "</td>")

	} else {
		buffer.WriteString("<td >N/A</td>")
	}

	buffer.WriteString("<td >" + model.GetTaskStatusText(task.Status) + "</td>")

	buffer.WriteString("</tr>")
	return buffer.String()

}

func writeTag(buff *bytes.Buffer, tag string, innerblock string) {
	buff.WriteString("<")
	buff.WriteString(tag)
	buff.WriteString(">")
	buff.WriteString(innerblock)
	buff.WriteString("</")
	buff.WriteString(tag)
	buff.WriteString(">")
}

func GetStatusCount(key string, kvs map[string]interface{}) interface{} {
	v := kvs[key]
	if v == nil {
		return 0
	}
	return v
}

func GetActive(i, j int) string {
	if i == j {
		return "class=uk-active"
	}
	return ""
}
