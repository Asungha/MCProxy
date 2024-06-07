package utils

import "fmt"

type PrometheusFormatter struct {
	data string
}

func (f *PrometheusFormatter) Add(metricName string, value string, filter map[string]string) *PrometheusFormatter {
	if len(filter) == 0 {
		f.data += fmt.Sprintf("\n%s %s", metricName, value)
	} else {
		buf := "{"
		for k, v := range filter {
			if buf != "{" {
				buf += ","
			}
			buf += fmt.Sprintf(`%s="%s"`, k, v)
		}
		buf += "}"
		f.data += fmt.Sprintf("\n%s%s %s", metricName, buf, value)
	}
	return f
}

func (f *PrometheusFormatter) Get() string {
	return f.data
}
