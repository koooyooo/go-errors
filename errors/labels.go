package errors

import "strings"

type Label string

type Labels struct {
	Labels []Label
}

func (k Labels) AddedClone(ls ...Label) *Labels {
	var added []Label
	for _, l := range k.Labels {
		added = append(added, l)
	}
	for _, l := range ls {
		added = append(added, l)
	}
	return &Labels{added}
}

func (k Labels) String() string {
	var s []string
	for _, l := range k.Labels {
		s = append(s, string(l))
	}
	return "[" + strings.Join(s, ",") + "]"
}

var NoLabel = L()

func L(args ...interface{}) *Labels {
	var labels []Label
	for _, i := range args {
		switch v := i.(type) {
		case Label:
			labels = append(labels, v)
		case []Label:
			for _, l := range v {
				labels = append(labels, l)
			}
		case Labels:
			for _, l := range v.Labels {
				labels = append(labels, l)
			}
		case *Labels:
			for _, l := range v.Labels {
				labels = append(labels, l)
			}
		}
	}
	return &Labels{Labels: labels}
}
