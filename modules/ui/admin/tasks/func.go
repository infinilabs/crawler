package tasks

import (
	"bytes"
	"fmt"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
)

func GetDomainRow(domain model.Domain) string {
	var buffer bytes.Buffer
	link := fmt.Sprintf("<a href=\"?domain=%v\">%v(%v)</a>", domain.Host, domain.Host, domain.LinksCount)
	writeTag(&buffer, "span", link)
	return buffer.String()
}

func GetTaskRow(task model.Task) string {
	var buffer bytes.Buffer
	buffer.WriteString("<tr>")

	writeTag(&buffer, "td", util.SubStringWithSuffix(task.Url, 83, "..."))
	date := util.FormatTime(task.UpdateTime)
	buffer.WriteString("<td class='timeago' title='" + date + "' >" + date + "</td>")

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
