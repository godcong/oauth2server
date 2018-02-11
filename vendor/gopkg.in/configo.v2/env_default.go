package configo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func envDefault(c *Config, f *os.File) error {
	var ckey string
	popn := false
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		var p interface{}
		if process := envDefaultLine(&b, &p); process == PROCESS_COMMON {
			ckey = p.(string)
			if ckey != "" {
				popn = true
				(c.Configure).(Default)[ckey] = make(Property)

			} else {
				popn = false
			}
		} else if process == PROCESS_PROPERTY {
			if popn {
				prop := (p).(Property)
				for k, v := range prop {
					(c.Configure).(Default)[ckey][k] = v
				}
			}
		} else {
			//skip
		}

	}
	fmt.Println(c.Configure)
	return nil
}

func envDefaultLine(line *[]byte, key *interface{}) int {
	proc := string(*line)
	//is common
	proc = strings.TrimSpace(proc)
	sta := strings.Index(proc, "[")
	end := strings.LastIndex(proc, "]")
	//check error
	rlt := PROCESS_NONE
	strarr := []string{"#", "//"}

	//get group name
	if sta >= 0 && end >= 0 {
		*key = ""
		if sta+1 < end {
			*key = proc[sta+1 : end]
		}
		rlt = PROCESS_COMMON
		goto LAST
	}

	//delete comment
	for _, v := range strarr {
		if idx := strings.Index(proc, v); idx > 0 {
			proc = strings.TrimSpace(strings.Split(proc, v)[0])
			break
		} else if idx == 0 {
			goto LAST
		} else {

		}

	}

	//string splite to map
	if mp, err := envStringSplit(proc, "="); err == nil {
		*key = mp
		rlt = PROCESS_PROPERTY
	}

LAST:
	return rlt
}

func envStringSplit(s string, sep string) (Property, error) {
	ss := strings.Split(s, sep)
	if len(ss) == 2 {
		key := strings.TrimSpace(ss[0])
		val := strings.TrimSpace(ss[1])
		return Property{key: val}, nil
	}

	return nil, ERROR_SPLIT_NO_DATA

}

func envDefaultGet(def Default, s string) *Property {
	if v, ok := def[s]; ok {
		return &v
	}
	return nil
}
