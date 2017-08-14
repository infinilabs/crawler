package tasks

import (
	"bytes"
	"fmt"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
)

// GetDomainRow return html blocks to display a domain info
func GetDomainRow(domain model.Domain) string {
	var buffer bytes.Buffer
	link := fmt.Sprintf("<a href=\"?domain=%v\">%v(%v)</a>", domain.Host, domain.Host, domain.LinksCount)
	writeTag(&buffer, "span", link)
	return buffer.String()
}

// GetTaskRow return html blocks to display a task info
func GetTaskRow(task model.Task) string {
	var buffer bytes.Buffer
	buffer.WriteString("<tr>")

	writeTag(&buffer, "td", util.SubStringWithSuffix(task.Url, 60, "..."))

	if task.SnapshotCreated != nil {
		date1 := util.FormatTimeWithLocalTZ(task.SnapshotCreated)
		buffer.WriteString("<td class='timeago' title='" + date1 + "' >" + date1 + "</td>")
	} else {
		buffer.WriteString("<td >N/A</td>")
	}

	if task.NextCheck != nil {
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
