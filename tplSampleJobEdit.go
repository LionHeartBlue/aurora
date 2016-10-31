package main

import (
	"bytes"
	"html"

	"github.com/kr/beanstalk"
)

// tplSampleJobEdit render a sample job edit form.
func tplSampleJobEdit(key string, alert string) string {
	var err error
	var buf, action, title, name, savedTo, saveTo, data, ST, tubeList bytes.Buffer
	if key == "" {
		action.WriteString(`?action=actionNewSample`)
		title.WriteString(`<h4 class="text-info">New sample job</h4>`)
	} else {
		action.WriteString(`?action=actionEditSample&key=`)
		action.WriteString(key)
		name.WriteString(html.EscapeString(getSampleJobNameByKey(key)))
		data.WriteString(html.EscapeString(getSampleJobDataByKey(key)))
		title.WriteString(`<h4 class="text-info">Edit: `)
		title.WriteString(name.String())
		title.WriteString(`</h4>`)
		for _, j := range sampleJobs.Jobs {
			if key == j.Key {
				for _, t := range j.Tubes {
					saveTo.WriteString(`<div class="control-group"><div class="controls"><label class="checkbox-inline"><input type="checkbox" name="tubes[`)
					saveTo.WriteString(t)
					saveTo.WriteString(`]" value="1" checked="checked">`)
					saveTo.WriteString(t)
					saveTo.WriteString(`</label></div></div>`)
				}
			}
		}
	}

	for _, server := range selfConf.Servers {
		var bstkConn *beanstalk.Conn
		tubeList.Reset()
		if bstkConn, err = beanstalk.Dial("tcp", server); err != nil {
			continue
		}
		tubes, _ := bstkConn.ListTubes()
		bstkConn.Close()
		for _, v := range tubes {
			var checked string
			for _, j := range sampleJobs.Jobs {
				if j.Key == key {
					for _, t := range j.Tubes {
						if t == v {
							checked = `checked="checked"`
						}
					}
				}
			}
			tubeList.WriteString(`<div class="control-group"><div class="controls"><label class="checkbox-inline"><input type="checkbox" name="tubes[`)
			tubeList.WriteString(v)
			tubeList.WriteString(`]" value="1" `)
			tubeList.WriteString(checked)
			tubeList.WriteString(`>`)
			tubeList.WriteString(v)
			tubeList.WriteString(`</label></div></div>`)
		}
		ST.WriteString(`<div class="pull-left" style="padding-right: 35px;">`)
		ST.WriteString(server)
		ST.WriteString(`<blockquote>`)
		ST.WriteString(tubeList.String())
		ST.WriteString(`</blockquote></div>`)
	}
	if name.String() != "" {
		savedTo.WriteString(`<div class="pull-left" style="padding-right: 35px;">Saved to: <blockquote>`)
		savedTo.WriteString(saveTo.String())
		savedTo.WriteString(`</blockquote></div>`)
	}
	buf.WriteString(`<form name="sampleJobsEdit" action="`)
	buf.WriteString(action.String())
	buf.WriteString(`" method="POST"><div class="clearfix form-group"><div class="pull-left">`)
	buf.WriteString(title.String())
	buf.WriteString(`</div><div class="pull-right"><a href="?action=manageSamples" class="btn btn-default btn-small"><i class="glyphicon glyphicon-list"></i> Manage samples</a></div></div><div class=" form-group"><fieldset>`)
	buf.WriteString(alert)
	buf.WriteString(`<div class="control-group"><label class="control-label" for="addsamplename"><b>Name *</b></label><div class="controls form-group"><input class="input-xlarge focused" id="addsamplename" name="name" type="text" value="`)
	buf.WriteString(name.String())
	buf.WriteString(`" autocomplete="off"></div></div></fieldset><div class="clearfix"><label class="control-label"><b>Available on tubes *</b></label><br/>`)
	buf.WriteString(savedTo.String())
	buf.WriteString(ST.String())
	buf.WriteString(`</div><div><label class="control-label" for="jobdata"><b>Job data *</b></label><textarea name="jobdata" id="jobdata" style="width:100%" rows="3">`)
	buf.WriteString(data.String())
	buf.WriteString(`</textarea></div></div><div><input type="submit" class="btn btn-success" value="Save"/></div></form>`)
	return buf.String()
}
